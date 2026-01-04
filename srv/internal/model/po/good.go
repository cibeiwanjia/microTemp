package po

import (
	"time"
)

type Goods struct {
	DeletedAt       time.Time `gorm:"column:deleted_at"`
	ID              int       `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	AddTime         time.Time `gorm:"column:add_time;NOT NULL"`
	IsDeleted       int       `gorm:"column:is_deleted"`
	UpdateTime      time.Time `gorm:"column:update_time;NOT NULL"`
	CategoryID      int       `gorm:"column:category_id;NOT NULL"`
	BrandsID        int       `gorm:"column:brands_id;NOT NULL"`
	OnSale          int       `gorm:"column:on_sale;NOT NULL"`
	GoodsSn         string    `gorm:"column:goods_sn;NOT NULL"`
	Name            string    `gorm:"column:name;NOT NULL"`
	ClickNum        int       `gorm:"column:click_num;NOT NULL"`
	SoldNum         int       `gorm:"column:sold_num;NOT NULL"`
	FavNum          int       `gorm:"column:fav_num;NOT NULL"`
	MarketPrice     float64   `gorm:"column:market_price;NOT NULL"`
	ShopPrice       float64   `gorm:"column:shop_price;NOT NULL"`
	GoodsBrief      string    `gorm:"column:goods_brief;NOT NULL"`
	ShipFree        int       `gorm:"column:ship_free;NOT NULL"`
	Images          string    `gorm:"column:images;NOT NULL"`
	DescImages      string    `gorm:"column:desc_images;NOT NULL"`
	GoodsFrontImage string    `gorm:"column:goods_front_image;NOT NULL"`
	IsNew           int       `gorm:"column:is_new;NOT NULL"`
	IsHot           int       `gorm:"column:is_hot;NOT NULL"`
}
