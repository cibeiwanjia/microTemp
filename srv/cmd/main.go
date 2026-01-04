package main

import (
	"fmt"
	"os"

	//引入app.NewServerCommand包
	app "github.com/cibeiwanjia/microTemp/srv/cmd/apps"
)

func main() {
	command := app.NewServerCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
