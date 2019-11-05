package routers

import (
	"BeeGo-project/classThree/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	//路由过滤器
	beego.InsertFilter("/Article/*",beego.BeforeRouter,filterFunc)

	//登录
    beego.Router("/", &controllers.LoginController{},
    		"get:ShowLogin;post:HandleLogin")

	//退出登录
    beego.Router("/Logout", &controllers.LoginController{},
    		"get:Logout")

    //注册
	beego.Router("/register", &controllers.RegController{},
		"get:ShowReg;post:HandleReg")

	//首页
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{},
		"get:ShowArticleList;post:HandleSelect")

	//添加文章
	beego.Router("/Article/AddArticle", &controllers.ArticleController{},
		"get:ShowAddArticle;post:HandleAddArticle")

	//显示文章详情
	beego.Router("/Article/showArticleDetail",&controllers.ArticleController{},
					"get:ShowArticleDetail")

	//显示文章详情
	beego.Router("/Article/DeleteArticle",&controllers.ArticleController{},
					"get:HandleDelete")

    //编辑
    beego.Router("/Article/updateArticle",&controllers.ArticleController{},
    	"get:ShowUpdateArticle;post:HandleUpdateArticle")

    //添加类型
	beego.Router("/Article/AddArticleType",&controllers.ArticleController{},
		"get:ShowAddType;post:HandleAddType")
}

var filterFunc=func (ctx *context.Context)  {
	userName:=ctx.Input.Session("userName")
	if userName==nil {
		ctx.Redirect(302,"/")
	}
}
