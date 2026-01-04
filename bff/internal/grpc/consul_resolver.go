package grpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

type consulResolver struct {
	cc         resolver.ClientConn
	consulAddr string
	service    string
	closeCh    chan struct{}
	wg         sync.WaitGroup
}

func (r *consulResolver) watch() {
	defer r.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.closeCh:
			return
		case <-ticker.C:
			addrs, err := r.resolve()
			if err == nil {
				r.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}
	}
}

func (r *consulResolver) resolve() ([]resolver.Address, error) {
	config := api.DefaultConfig()
	config.Address = r.consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	services, _, err := client.Health().Service(r.service, "", true, nil)
	if err != nil {
		return nil, err
	}

	addrs := make([]resolver.Address, 0, len(services))
	for _, s := range services {
		addrs = append(addrs, resolver.Address{
			Addr: fmt.Sprintf("%s:%d", s.Service.Address, s.Service.Port),
		})
	}

	return addrs, nil
}

func (r *consulResolver) ResolveNow(resolver.ResolveNowOptions) {
	addrs, err := r.resolve()
	if err == nil {
		r.cc.UpdateState(resolver.State{Addresses: addrs})
	}
}

func (r *consulResolver) Close() {
	close(r.closeCh)
	r.wg.Wait()
}
