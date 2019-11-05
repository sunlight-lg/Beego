package controllers

import (
	"BeeGo-project/shFresh/models"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
	"strconv"
)

type UserController struct {
	beego.Controller
}

//展示登录页面
func (this *UserController) ShowLogin() {
	name:=this.GetString("userName")
	//解码
	temp,_:=base64.StdEncoding.DecodeString(name)
	if string(temp)=="" {
		this.Data["userName"]=""
		this.Data["checked"]=""
	}else{
		this.Data["userName"]=string(temp)
		this.Data["checked"]="checked"
	}
	beego.Info(temp)
	this.TplName="login.html"
}

//登录操作的处理
func (this *UserController) HandleLogin() {
	name:=this.GetString("username")
	pwd:=this.GetString("pwd")

	if name==""||pwd=="" {
		this.Data["errmsg"]="用户名密码不能为空"
		this.TplName="login.html"
		return
	}

	o:=orm.NewOrm()
	var user models.User
	user.Name=name
	user.PassWord=pwd

	err:=o.Read(&user,"Name")
	if err!=nil {
		this.Data["errmsg"]="用户名密码错误，请重新输入"
		this.TplName="login.html"
		return
	}

	if user.PassWord!=pwd {
		this.Data["errmsg"]="用户名密码错误，请重新输入"
		this.TplName="login.html"
		return
	}

	if user.Active!=true {
		this.Data["errmsg"]="用户未激活，请先往邮箱激活在进行登录操作"
		this.TplName="login.html"
		return
	}

	//记住用户名
	remember:=this.GetString("remember")

	//base64加密
	if remember=="on" {
		//加密
		temp:=base64.StdEncoding.EncodeToString([]byte(name))
		this.Ctx.SetCookie("userName",temp,24*3600*30)
	}else{
		this.Ctx.SetCookie("userName",name,-1)
	}

	this.TplName="index.html"
}

//注册页面
func (this *UserController) ShowReg() {
	this.TplName="register.html"
}

//注册处理操作
func (this *UserController) HandleReg() {
	userName:=this.GetString("user_name")
	pwd:=this.GetString("pwd")
	cpwd:=this.GetString("cpwd")
	email:=this.GetString("email")

	if userName==""||pwd==""||cpwd==""||email=="" {
		this.Data["errmsg"]="数据不完整，请重新注册"
		this.TplName="register.html"
		return
	}

	if pwd!=cpwd {
		this.Data["errmsg"]="两次输入密码不一致，请重新注册"
		this.TplName="register.html"
		return
	}

	reg,_:=regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res:=reg.FindString(email)
	if res=="" {
		this.Data["errmsg"]="邮箱格式不正确，请重新输入"
		this.TplName="register.html"
		return
	}

	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	user.PassWord=pwd
	user.Email=email

	_,err:=o.Insert(&user)
	if err!=nil {
		this.Data["errmsg"]="注册失败，请更换数据重新注册"
		this.TplName="register.html"
		return
	}

	//发送邮件`{"username":"563364657@qq.com","password":"cgapyzgkkczubdea","host":"smtp.qq.com","port":587}`
	config := `{"username":"1503741454@qq.com","password":"ycewnentkegqhdei","host":"smtp.qq.com","port":587}`
	emailConn := utils.NewEMail(config)
	emailConn.To = []string{email}
	emailConn.From = "1503741454@qq.com"   //这个一定是发件人邮箱地址
	emailConn.Subject = "天天生鲜用户注册"
	emailConn.HTML = "http://192.168.229.130/active?id="+strconv.Itoa(user.Id)
	// email.AttachFile("1.jpg") // 附件
	// email.AttachFile("1.jpg", "1") // 内嵌资源
	err = emailConn.Send()
	if err != nil {
		this.Data["errmsg"]="发送邮件失败，请重新注册"
		fmt.Println("邮件发送失败原因是=",err)
		this.TplName="register.html"
		return
	}

	//返回视图
	this.Ctx.WriteString("注册成功，请激活用户在进行登录操作")
}

//激活处理
func (this *UserController) ActiveUser() {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["errmsg"]="要激活的用户不存在"
		this.TplName="register.html"
		return
	}

	o:=orm.NewOrm()
	var user models.User
	user.Id=id
	err=o.Read(&user)
	if err!=nil {
		this.Data["errmsg"]="要激活的用户不存在"
		this.TplName="register.html"
		return
	}

	user.Active=true
	o.Update(&user)

	//返回页面
	this.Redirect("/login",302)
}
