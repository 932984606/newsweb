package controllers

import (
	"github.com/astaxie/beego"
	"newsWeb/models"
	"github.com/astaxie/beego/orm"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister(){
	this.TplName = "register.html"
}

func (this *UserController) HandleReg(){
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	if userName == "" || pwd == ""{
		this.Data["errmsg"] = "用户名和密码不能为空"
		this.TplName = "register.html"
		return
	}

	var user models.User
	user.Username = userName
	user.Pwd = pwd
	o := orm.NewOrm()
	_,err := o.Insert(&user)
	if err != nil{
		this.Data["errmsg"] = "注册失败"
		this.TplName = "register.html"
		return
	}
	//this.Ctx.WriteString("注册成功")
	//this.TplName = "login.html"
	this.Redirect("/login",302)
}

func (this *UserController) ShowLogin(){
	username := this.Ctx.GetCookie("userName")
	etc,_ := base64.StdEncoding.DecodeString(username)
	if username != ""{
		this.Data["username"] = string(etc)
		this.Data["checked"] = "checked"
	}

	this.TplName = "login.html"
}

func (this *UserController) HandleLogin(){
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	if userName == "" || pwd == ""{
		this.Data["errmsg"] = "用户名和密码不能为空"
		this.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Username = userName
	err := o.Read(&user,"Username")
	if err != nil{
		this.Data["errmsg"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}

	if user.Pwd != pwd{
		this.Data["errmsg"] = "密码错误，请重新输入"
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")
	if remember == "on"{
		etc := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",etc,3600)
	}else {
		this.Ctx.SetCookie("userName",userName,-1)
	}

	this.SetSession("username",userName)
	//this.Ctx.WriteString("登录成功")
	this.Redirect("/article/articleList",302)
}

func (this *UserController) Logout(){
	this.DelSession("username")

	this.Redirect("/login",302)
}
