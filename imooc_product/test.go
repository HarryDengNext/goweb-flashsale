package main

import (
	"imooc_product/common"
	"imooc_product/datamodels"
)

func main() {
	data := map[string]string{"ID":"1", "productName": "imooc测试结构体", "productNum": "3", "productImage": "123", "productUrl": "http://url"}
	product := &datamodels.Product{}
	common.DataToStructByTagSql(data, product)
}
