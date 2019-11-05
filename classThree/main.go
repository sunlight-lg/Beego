package main

import (
	_ "BeeGo-project/classThree/routers"
	"github.com/astaxie/beego"
	_"BeeGo-project/classThree/models"
	"strconv"
)

func main() {
	beego.AddFuncMap("ShowPrePage",HandlePerPage)
	beego.AddFuncMap("ShowNextPage",HandleNextPage)
	beego.Run()
}

func HandlePerPage(data int) (string) {
	//dataTime,_:=strconv.Atoi(data)
	pageIndex:=data-1
	pageIndex1:=strconv.Itoa(pageIndex)
	return pageIndex1
}

func HandleNextPage(data int) (int) {
	pageIndex:=data+1
	return pageIndex
}