package controllers

import (
	"BeeGo-project/classThree/models"
	"bytes"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"math"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

/**
显示文章列表
*/
func (this *ArticleController) ShowArticleList() {
	userName:=this.GetSession("userName")
	if userName==nil {
		this.Redirect("/",302)
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)

	pageIndex := this.GetString("pageIndex")
	pageIndex1, err := strconv.Atoi(pageIndex)
	if err != nil {
		pageIndex1 = 1
	}

	//获取总页面数
	pageSize := 3
	start := pageSize * (pageIndex1 - 1)

	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles) //pageSize每页显示多少条数据，2.start起始位置


	count, err := qs.RelatedSel("ArticleType").Count() //返回数据条目数

	pageCount := float64(count) / float64(pageSize)
	pageCount1 := math.Ceil(pageCount)

	if err != nil {
		beego.Info("查询错误")
		return
	}

	//首页末页数据处理
	FirstPage:=false
	if pageIndex1 == 1 {
		FirstPage=true
	}

	EndPage:=false
	if pageIndex1==int(pageCount1) {
		EndPage=true
	}

	//获取类型数据
	var types []models.ArticleType
	conn,err:=redis.Dial("tcp",":6379")
	buffeer,err:=redis.Bytes(conn.Do("get","types"))
	if err!=nil {
		beego.Info("获取redis数据错误")
	}

	dec:=gob.NewDecoder(bytes.NewReader(buffeer))
	dec.Decode(&types)

	beego.Info("dec=",dec)

	if len(types)==0 {
		//从MySQL取数据
		o.QueryTable("ArticleType").All(&types)
		var buffer bytes.Buffer
		enc:=gob.NewEncoder(&buffer)
		err:=enc.Encode(&types)

		defer conn.Close()
		_,err=conn.Do("set","types",buffer.Bytes())
		if err!=nil {
			beego.Info("redis数据库操作错误")
			return
		}

		beego.Info("从MySQL中取数据")
	}


	//根据类型获取数据
	typeName:=this.GetString("select")
	var articleswithtype []models.Article
	if typeName==""{
		//beego.Info("下拉框获取数据失败")
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articleswithtype)
	}else{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articleswithtype)
	}

	this.Data["userName"]=userName
	this.Data["types"]=types
	this.Data["typeName"]=typeName
	this.Data["FirstPage"] = FirstPage
	this.Data["EndPage"] = EndPage
	this.Data["count"] = count
	this.Data["pageCount1"] = pageCount1
	this.Data["articles"] = articleswithtype
	this.Data["pageIndex"] = pageIndex1

	this.Layout="layout.html"
	this.TplName = "index.html"
}

//处理下拉框改变发的请求
func (this *ArticleController) HandleSelect() {
	typeName:=this.GetString("select")
	if typeName==""{
		beego.Info("下拉框获取数据失败")
		return
	}

	o:=orm.NewOrm()
	var articles []models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	beego.Info("articles=",articles)
}

/**
添加
*/
func (this *ArticleController) ShowAddArticle() {
	//查询类型数据
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	this.Data["types"]=types

	this.Layout="layout.html"
	this.TplName = "add.html"
}

func (this *ArticleController) HandleAddArticle() {
	//拿数据
	artiName := this.GetString("articleName")
	artiContent := this.GetString("content")

	f, h, err := this.GetFile("uploadname")
	defer f.Close()

	//文件格式处理
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Info("上传文件格式不正确")
		return
	}

	//文件大小判断
	if h.Size > 50000000 {
		beego.Info("上传文件过大，请重新上传")
		return
	}

	//不能重名
	fileName := time.Now().Format("2006-01-02 15:04:05")

	this.SaveToFile("uploadname", "./static/img/upload/"+fileName+ext)
	if err != nil {
		beego.Info("文件上传失败")
		return
	}

	o := orm.NewOrm()
	article := models.Article{}
	article.Title = artiName
	article.Content = artiContent
	article.Img = "./static/img/upload/" + fileName + ext

	//给Article对象赋值
	typeName:=this.GetString("select")
	if typeName=="" {
		beego.Info("类型数据错误")
		return
	}
	var articleType models.ArticleType
	articleType.TypeName=typeName
	err=o.Read(&articleType,"TypeName")
	if err!=nil {
		beego.Info("获取类型错误")
		return
	}
	article.ArticleType=&articleType

	//插入数据
	_, err = o.Insert(&article)
	if err != nil {
		beego.Info("增加文章失败")
		return
	}

	this.Redirect("/Article/ShowArticle", 302)
}

func (this *ArticleController) ShowArticleDetail() {
	id, err := this.GetInt("articleId")
	if err != nil {
		beego.Info("传递的链接有误")
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	o.Read(&article)

	article.Count += 1
	//article=models.Article{Id:id}
	m2m:=o.QueryM2M(&article,"Users")
	userName:=this.GetSession("userName")
	user:=models.User{}
	user.UserName=userName.(string)
	o.Read(&user,"UserName")
	_,err=m2m.Add(&user)
	if err!=nil {
		beego.Info("插入错误")
		return
	}
	o.Update(&article)

	var users[] models.User
	//o.LoadRelated(&article,"Users")
	o.QueryTable("User").Filter("Articles__Article__id",id).Distinct().All(&users)

	beego.Info("article user=",article)

	this.Data["article"] = article
	this.Data["users"] = users

	this.Layout="layout.html"
	this.TplName = "content.html"
}

func (this *ArticleController) HandleDelete() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("传递的链接有误")
	}

	o := orm.NewOrm()
	article := models.Article{
		Id: id,
	}

	o.Delete(&article)

	this.Redirect("/Article/ShowArticle", 302)
}

func (this *ArticleController) ShowUpdateArticle() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("请求文章错误")
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)

	//返回给视图
	this.Data["article"] = article

	this.Layout="layout.html"
	this.TplName = "update.html"
}

func (this *ArticleController) HandleUpdateArticle() {
	id, err := this.GetInt("id")
	aname := this.GetString("articleName")
	content := this.GetString("content")
	filePath := UploadFile(&this.Controller, "uploadname")
	beego.Info("aname=", aname, "content=", content, "filePath=", filePath)
	if err != nil || aname == "" || content == "" || filePath == "" {
		beego.Info("请求错误")
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		beego.Info("更新的文章不存在")
		return
	}

	article.Title = aname
	article.Content = content
	if filePath != "NoImg" {
		article.Img = filePath
	}
	o.Update(&article)

	//返回视图
	this.Redirect("/Article/ShowArticle", 302)
}

func (this *ArticleController) ShowAddType() {
	//显示数据
	o:=orm.NewOrm()
	var artiType[] models.ArticleType

	//查询
	_,err:=o.QueryTable("ArticleType").All(&artiType)
	if err!=nil {
		beego.Info("查询类型错误")
	}

	this.Data["types"]=artiType

	this.Layout="layout.html"
	this.TplName="addType.html"
}

func (this *ArticleController) HandleAddType() {
	tname:=this.GetString("typeName")

	if tname=="" {
		beego.Info("数据不能为空")
		return
	}

	o := orm.NewOrm()
	var aType models.ArticleType
	aType.TypeName=tname
	_,err:=o.Insert(&aType)
	if err!=nil {
		beego.Info("增加失败")
		return
	}

	this.Redirect("/Article/AddArticleType",302)
}
