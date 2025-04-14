package controllers

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"imooc_product/common"
	"imooc_product/datamodels"
	"imooc_product/services"
	"strconv"
)

type OrderController struct {
	Ctx iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) GetAll() mvc.View {
	orderMap, err := o.OrderService.GetAllOrderInfo()
	if err!= nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
		o.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name:"/order/view.html",
		Data:iris.Map{
			"order": orderMap,
		},
	}
}

func (o *OrderController) GetManager() mvc.View  {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	order, err := o.OrderService.GetOrderByID(id)
	if err!= nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name:"order/manager.html",
		Data:iris.Map{
			"order": order,
		},
	}
}

func (o *OrderController) PostUpdate()  {
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{"imooc", false, false})
	if err := dec.Decode(o.Ctx.Request().Form, order);err!= nil{
		o.Ctx.Application().Logger().Debug(err)
	}
	err := o.OrderService.UpdateOrder(order)
	if err!= nil{
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}

func (o *OrderController) PostDelete() []byte{
	//idString := o.Ctx.URLParam("id")
	idString := o.Ctx.PostValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	isOk := o.OrderService.DeleteOrderByID(id)
	if isOk {
		o.Ctx.Application().Logger().Debug("删除订单成功，ID为"+idString)
		response, _ := json.Marshal(iris.Map{"code": 200,"msg": "删除订单成功"})
		return response
	}else {
		o.Ctx.Application().Logger().Debug("删除订单失败，ID为"+idString)
		response, _ := json.Marshal(iris.Map{"code": 201,"msg": "删除订单失败"})
		return response
	}
	//o.Ctx.Redirect("/order/all")
}