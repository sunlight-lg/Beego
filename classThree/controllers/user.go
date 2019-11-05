package controllers

import (
	"BeeGo-project/classThree/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

/**
登录操作
 */
type LoginController struct {
	beego.Controller
}

func (this *LoginController) ShowLogin() {
	name:=this.Ctx.GetCookie("userName")
	if name!="" {
		this.Data["userName"]=name
		this.Data["check"]="checked"
	}
	this.TplName = "login.html"
}

func (this *LoginController) HandleLogin() {
	//拿到浏览器传递的数据
	name:=this.GetString("userName")
	passwd:=this.GetString("password")

	//数据处理
	if name==""||passwd=="" {
		beego.Info("用户名或者密码输入有误")
		this.TplName = "login.html"
		return
	}

	//查询数据库
	o:=orm.NewOrm()
	user:=models.User{}
	user.UserName=name
	err:=o.Read(&user,"userName")
	if err!=nil {
		beego.Info("用户名或者密码输入错误")
		this.TplName = "login.html"
		return
	}

	//判断密码是否一致
	if user.Password!=passwd {
		beego.Info("密码输入错误")
		this.TplName = "login.html"
		return
	}

	check:=this.GetString("userName")
	if check=="on" {
		this.Ctx.SetCookie("userName",name,time.Second*3600)
	}else{
		this.Ctx.SetCookie("userName","ssss",-1)
	}

	this.SetSession("userName",name)


	//this.TplName="index.html"
	this.Redirect("/Article/ShowArticle",302)
}

/**
注册操作
 */
type RegController struct {
	beego.Controller
}

func (this *RegController) ShowReg() {
	this.TplName = "register.html"
}

func (this *RegController) HandleReg() {
	//拿到浏览器传递的数据
	name:=this.GetString("userName")
	passwd:=this.GetString("password")

	//数据处理
	if name==""||passwd=="" {
		beego.Info("用户名或者密码输入有误")
		this.TplName = "login.html"
		return
	}

	//查询数据库
	o:=orm.NewOrm()
	user:=models.User{}
	user.UserName=name
	user.Password=passwd

	_,err:=o.Insert(&user)
	if err!=nil {
		beego.Info("注册失败")
	}

	//返回登录页面
	this.Redirect("/",302)
}

/**
退出登录
 */
func (this *LoginController) Logout() {
	this.DelSession("userName")
	this.Redirect("/",302)
}