package model

import (
	"encoding/json"
	"mixed/beas/file"
	"mixed/beas/pool"
	"sync"
	"time"
)

// cc_order_products_user
type OrderProductsUser struct {
	Id             int
	ProductId      int64
	ProductSkuId   int
	ProductVersion int
	ProductImgUrl  string
}

// cc_product_skus
type ProductSkus struct {
	Id         int
	ProductId  int
	SkuPicture string
}

// cc_product_basic
type ProductBasic struct {
	Id            int
	ProductId     int
	ImageUrlsHead string `xorm:"varchar(2048)"`
}

// get original order data
func GetOrderWorker(begin int, end int, workers pool.RoutinePool, wg sync.WaitGroup) {
	var orders = make([]OrderProductsUser, 0)

	//  select * from cc_order_products_user where id between ? and ?;
	if err := ShopEngine.Where("id between ? and ?", begin, end).Find(&orders); err == nil {
		if len(orders) == 0 {
			workers.Return()
		}
		for _, order := range orders {
			if len(order.ProductImgUrl) == 99 {
				format := "ID=%v---ProductId=%v---SkuId=%v---ProductVersion=%v---ProductUrl=%v\n"
				file.WriteDataFile(file.OriFd, format, order.Id, order.ProductId, order.ProductSkuId, order.ProductVersion, order.ProductImgUrl)
				wg.Add(1)
				go GetProductWorker(order, wg)
			}
		}
	} else {
		fmt := "Get order data err : %v\n"
		file.WriteDataFile(file.ErrFd, fmt, err.Error())
		panic("Init xorm engine again err :" + err.Error())
	}
	time.Sleep(500 * time.Millisecond)
	workers.Return()
}

// get corect product data
func GetProductWorker(orderProduct OrderProductsUser, wg sync.WaitGroup) {
	defer wg.Done()
	// get data from cc_product_skus
	// UPDATE cc_order_products_user set `product_img_url` = "http://xxx.com/test.jpg" where id = 123
	skus := new(ProductSkus)
	format := "UPDATE cc_order_products_user set `product_img_url` = \"%v\" where id = %v\n"
	if orderHas, err := ShopEngine.Where("id = ?", orderProduct.ProductSkuId).Get(skus); err == nil {
		if orderHas && (len(skus.SkuPicture) != 0) {
			file.WriteDataFile(file.NewFd, format, skus.SkuPicture, orderProduct.Id)
		} else {
			// get data from cc_product_basic
			products := make([]ProductBasic, 0)
			if err = ShopEngine.Where("product_id = ?", orderProduct.ProductId).Limit(1, 0).Find(&products); err == nil {
				pLen := len(products)
				if pLen > 0 {
					data := products[pLen-1]
					urls := make([]string, 0)
					err = json.Unmarshal([]byte(data.ImageUrlsHead), &urls)
					if err == nil {
						if imgLen := len(urls); imgLen > 0 {
							img := urls[0]
							file.WriteDataFile(file.NewFd, format, img, orderProduct.Id)
						} else {
							errMsg := "Get Product Img Err:orderId=%v,productId=%v,productUrl=%v\n"
							file.WriteDataFile(file.ErrFd, errMsg, orderProduct.Id, orderProduct.ProductId, data.ImageUrlsHead)
						}

					}
				}
			} else {
				fmt := "Get product data err : %v\n"
				file.WriteDataFile(file.ErrFd, fmt, err.Error())
				panic("Init xorm engine again err :" + err.Error())

			}
		}
	} else {
		fmt := "Get sku data err : %v\n"
		file.WriteDataFile(file.ErrFd, fmt, err.Error())
		panic("Init xorm engine again err :" + err.Error())
	}
}
