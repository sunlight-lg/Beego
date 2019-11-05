package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
)

func UploadFile(this *beego.Controller,filePath string) string{
	beego.Info("进来了没有")
	//处理文件上传
	file,head,err:=this.GetFile(filePath)
	if head.Filename=="" {
		return "NoImg"
	}

	if err!=nil {
		this.Data["errmsg"]="文件上传失败"
		this.TplName="add.html"
		return ""
	}

	defer file.Close()

	//文件大小
	if head.Size>500000 {
		this.Data["errmsg"]="文件太大，上传失败"
		this.TplName="add.html"
		return ""
	}

	//文件格式
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"]="文件格式错误，上传失败"
		this.TplName="add.html"
		return ""
	}

	//防止重名
	fileName:=time.Now().Format("2006-01-02 15:04:05")+ext
	//存储
	this.SaveToFile(filePath,"./static/img/upload/" +fileName)

	return "/static/img/upload/"+fileName
}

