package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	UserName string
	Password string
	Articles[] *Article `orm:"rel(m2m)"`
}

type Article struct {
	Id int    `orm:"pk;auto"`        //文章id
	Title string  `orm:"size(20)"`    //文章标题
	Content string  `orm:"size(500)"`  //内容
	Img string    `orm:"size(50);null"`    //图片路径
	//Type string       //类型
	Time time.Time  `orm:"type(datetime);auto_now_add"`  //发布时间
	Count int   `orm:"default(0)"`      // 阅读量
	ArticleType *ArticleType  `orm:"rel(fk)"`
	Users[] *User `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles[] *Article  `orm:"reverse(many)"`
}

func init()  {
	//连接数据库
	orm.RegisterDataBase("default","mysql",
		"root:root@(localhost:3306)/newsWeb?charset=utf8")
	//注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//生成表
	orm.RunSyncdb("default",false,true)
}
