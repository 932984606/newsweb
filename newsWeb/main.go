package main

import (
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
	_"newsWeb/models"
)

func main() {
	beego.AddFuncMap("prePage",prePage)
	beego.AddFuncMap("nextPage",nextPage)
	beego.Run()
}

func prePage(index int) int{
	if index <= 1{
		return 1
	}
	return index-1
}

func nextPage(index int,count int) int{
	if count < 1{
		count = 1
	}
	if index >= count{
		return count
	}
	return index+1
}

