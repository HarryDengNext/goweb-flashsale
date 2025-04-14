package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"imooc_product/datamodels"
	"imooc_product/encrypt"
	"imooc_product/services"
	"imooc_product/tool"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (u *UserController)GetRegister()mvc.View  {
	return mvc.View{
		Name:"user/register.html",
	}
}

func (u *UserController) PostRegister()  {
	var (
		nickName = u.Ctx.FormValue("nickName")
		userName = u.Ctx.FormValue("userName")
		passWord = u.Ctx.FormValue("passWord")
	)
	// ozzo-validattion

	user := &datamodels.User{
		NickName:     nickName,
		UserName:     userName,
		HashPassword: passWord,
	}
	_, err := u.Service.AddUser(user)
	if err != nil {
		u.Ctx.Application().Logger().Debug(err)
		u.Ctx.Redirect("/user/error")
		return
	}
	u.Ctx.Redirect("/user/login")
}

func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name:"user/login.html",
	}
}

func (u *UserController) PostLogin() mvc.Response {
	// 1. 获取用户提交的表单信息
	var (
		userName  = u.Ctx.FormValue("userName")
		passWord  = u.Ctx.FormValue("passWord")
	)
	// 2. 验证账号密码正确
	user, isOk := u.Service.IsPwdSuccess(userName, passWord)
	if !isOk {
		return mvc.Response{Path:"/user/login",}
	}
	// 3. 写入用户ID到cookie中
	tool.GlobalCookie(u.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err!=nil {
		u.Ctx.Application().Logger().Debug(err)
	}
	// 写入用户浏览器
	tool.GlobalCookie(u.Ctx, "sign", uidString)
	//u.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	return mvc.Response{
		Path:"/product/",
	}
}