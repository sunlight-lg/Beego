package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	_"github.com/go-sql-driver/mysql"
)

//用户表
type User struct {
	Id int
	Name string `orm:"size(20);unique"`   //用户名
	PassWord string `orm:"size(20)"`      //登录密码
	Email string `orm:"size(20)"`  //邮箱
	Active bool `orm:"default(false)"`    //是否激活
	Power int `orm:"default(0)"`          //权限设置，0表示未激活，1表示激活
	Address []*Address `orm:"reverse(many)"`
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

//地址表
type Address struct {
	Id int
	Receiver string `orm:"size(20)"`   //收件人
	Addr string `orm:"size(50)"`       //收件地址
	Zipcode string `orm:"size(20)"`    //邮编
	Phone string `orm:"size(20)"`      //联系方式
	Isdefault bool `orm:"default(false)"`  //是否默认，0为非默认，1为默认
	User *User `orm:"rel(fk)"`  //用户id
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

//商品SPU表
type Goods struct {
	Id int
	Name string `orm:"size(50)"`  //商品名称
	Datail string `orm:"size(250)"` //详细描述
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
}

//商品类型表
type GoodsType struct {
	Id int
	Name string   //种类名称
	Logo string   //logo
	Image string  //图片
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"`
}

//商品SKU表
type GoodsSKU struct {
	Id int
	Goods *Goods `orm:"rel(fk)"`   //商品SPU
	GoodsType *GoodsType `orm:"rel(fk)"`  //商品所属种类
	Name string   //商品名称
	Desc string   //商品简介
	Price int   //商品价格
	Unite string   //商品单位
	Image string   //商品图片
	Stock int `orm:"default(1)"`  //商品库存
	Sales int `orm:"default(0)"`  //商品销量
	Status int `orm:"default(1)"`  //商品状态
	Time time.Time `orm:"auto_now_add"`  //添加时间
	GoodsImage []*GoodsImage `orm:"reverse(many)"`
	IndexGoodsBanner []*IndexGoodsBanner `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"`
	OrderGoods []*OrderGoods  `orm:"reverse(many)"`
}

//商品图片表
type GoodsImage struct {
	Id int
	Image string //商品图片
	GoodsSKU *GoodsSKU  `orm:"rel(fk)"`  //商品SKU
}

//首页轮播商品展示表
type IndexGoodsBanner struct {
	Id int
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`  //商品SKU
	Image string                        //商品图片
	Index int    `orm:"default(0)"`     //展示顺序
}

//首页分类商品展示表
type IndexTypeGoodsBanner struct {
	Id int
	GoodsType *GoodsType `orm:"rel(fk)"` //商品类型
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`  //商品SKU
	DisplayType int  `orm:"default(1)"` //展示类型，0表示文字，1表示图片
	Index int    `orm:"default(0)"`     //展示顺序
}

//首页促销商品展示表
type IndexPromotionBanner struct {
	Id int
	Name string `orm:"size(20)"`    //活动名称
	Url string `orm:"size(50)"`     //活动链接
	Image string `orm:"size(20)"`   //活动图片
	Index int `orm:"default(0)"`     //展示顺序
}

//订单表
type OrderInfo struct {
	Id int
	OrderId string `orm:"unique"`
	User *User `orm:"rel(fk)"`   //用户
	Address *Address `orm:"rel(fk)"`   //地址
	PayMethod int    //付款方式
	TotalCount int `orm:"default(1)"`   //商品数量
	TotalPrice int  //商品价格
	TransitPrice int  //运费
	Orderstatus int `orm:"default(1)"` //订单状态
	TradeNo string `orm:"default('')"` //支付编号
	Time time.Time `orm:"auto_now_add"`  //评论时间

	OrderGoods []*OrderGoods `orm:"reverse(many)"`
}

//订单商品表
type OrderGoods struct {
	Id int
	OrderInfo *OrderInfo `orm:"rel(fk)"`   //订单
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`   //商品
	Count int `orm:"default(1)"`   //商品数量
	Price int //商品价格
	Comment string `orm:"default('')"`  //评论
}

func init()  {
	orm.RegisterDataBase("default","mysql",
		"root:root@tcp(127.0.0.1:3306)/dailyfresh?charset=utf8")

	orm.RegisterModel(new(User),new(Address),new(Goods),new(GoodsType),new(GoodsSKU),new(GoodsImage),
		new(IndexGoodsBanner),new(IndexTypeGoodsBanner),
		new(IndexPromotionBanner),new(OrderInfo),new(OrderGoods),)

	orm.RunSyncdb("default",false,true)
}