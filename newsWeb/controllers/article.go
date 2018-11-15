package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"math"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowArticle(){
	typename := this.GetString("select")
	o := orm.NewOrm()
	articles := make([]models.Article,0)

	qs := o.QueryTable("Article")
	//qs.All(&articles)

	var Count int64
	if typename == ""{
		Count,_ = qs.Count()
	}else {
		Count,_ = qs.Filter("ArticleType__TypeName",typename).RelatedSel("ArticleType").Count()
	}

	pageSize := 5
	pageCount := math.Ceil(float64(Count)/float64(pageSize))
	pageIndex,err := this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}
	if pageCount <= 0{
		pageCount = 1
	}
	start := pageSize*(pageIndex-1)

	if typename == ""{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else {
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typename).All(&articles)
	}

	for i:=0;i<len(articles);i++{
		if articles[i].ArticleType == nil{
			articles[i].ArticleType = &models.ArticleType{TypeName:""}
		}
	}

	articletypes := make([]models.ArticleType,0)

	qs = o.QueryTable("ArticleType")
	qs.All(&articletypes)

	this.Data["articletypes"] = articletypes
	this.Data["TName"] = typename
	this.Data["pageIndex"] = pageIndex
	this.Data["Count"] = Count
	this.Data["pageCount"] = int(pageCount)
	this.Data["articles"] = articles
	Layout(this,"后台管理页面")
	this.TplName = "index.html"
}

func (this *ArticleController) ShowAddarticle(){
	articletypes := make([]models.ArticleType,0)

	o := orm.NewOrm()
	qs := o.QueryTable("ArticleType")
	qs.All(&articletypes)

	this.Data["articletypes"] = articletypes

	Layout(this,"添加文章内容")
	this.TplName = "add.html"
}

func (this *ArticleController) HandlerAddarticle(){
	title := this.GetString("articleName")
	content := this.GetString("content")
	if title == "" || content == ""{
		this.Data["errmsg"] = "文章名称和内容不能为空"

		Layout(this,"添加文章内容")
		this.TplName = "add.html"
		return
	}

	file,head,err := this.GetFile("uploadname")
	if err != nil{
		this.Data["errmsg"] = "文件传输错误"

		Layout(this,"添加文章内容")
		this.TplName = "add.html"
		return
	}
	defer file.Close()
	if head.Size > 500000{
		this.Data["errmsg"] = "文件大小不符合要求"

		Layout(this,"添加文章内容")
		this.TplName = "add.html"
		return
	}

	fileExt := path.Ext(head.Filename)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg"{
		this.Data["errmsg"] = "文件类型不符合要求"

		Layout(this,"添加文章内容")
		this.TplName = "add.html"
		return
	}

	fileName := time.Now().Format("20060102150405")+fileExt
	this.SaveToFile("uploadname","./static/image/"+fileName)

	//添加文章类型
	articletypename := this.GetString("select")
	//beego.Info("articletype:",articletype)
	var articletype models.ArticleType
	articletype.TypeName = articletypename

	o := orm.NewOrm()
	var article models.Article
	article.Content = content
	article.Title = title
	article.Image = "/static/image/"+fileName

	err = o.Read(&articletype,"TypeName")
	if err != nil{
		this.Data["errmsg"] = "添加文章失败"
		this.TplName = "add.html"
		return
	}
	article.ArticleType = &articletype

	_,err = o.Insert(&article)
	if err != nil{
		this.Data["errmsg"] = "添加文章失败"

		Layout(this,"添加文章内容")
		this.TplName = "add.html"
		return
	}

	this.Redirect("/article/articleList",302)
}

func (this *ArticleController) ShowArticleDetial(){
	id,err := this.GetInt("id")
	if err != nil{
		beego.Error("访问路径错误",err)
		this.Redirect("/article/articleList",302)
		return
	}
	var articel models.Article
	articel.Id = id
	o := orm.NewOrm()
	err = o.Read(&articel)
	if err != nil{
		beego.Error("访问路径错误",err)
		this.Redirect("/article/articleList",302)
		return
	}
	if articel.ArticleType == nil{
		articel.ArticleType = &models.ArticleType{TypeName:""}
	}else {
		_,err = o.LoadRelated(&articel,"ArticleType")
	}
	if err != nil{
		beego.Error("访问路径错误",err)
		this.Redirect("/article/articleList",302)
		return
	}

	m2m := o.QueryM2M(&articel,"Users")
	var user models.User
	username := this.GetSession("username")
	user.Username = username.(string)
	err = o.Read(&user,"Username")
	if err != nil{
		beego.Error("读取数据失败：",err)
		this.Redirect("/article/articleList",302)
		return
	}
	m2m.Add(user)

	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)

	this.Data["users"] = users
	this.Data["article"] = articel

	Layout(this,"文章详情")
	this.TplName = "content.html"

	articel.ReadCount++
	if articel.ArticleType.TypeName == ""{
		articel.ArticleType = nil
	}
	_,err = o.Update(&articel)
	if err != nil{
		beego.Error("修改阅读次数失败：",err)
		this.Redirect("/article/articleList",302)
		return
	}
}

func (this *ArticleController) ShowUpdate(){
	id,err := this.GetInt("id")
	errmsg := this.GetString("errmsg")

	if err != nil{
		beego.Error("访问路径错误",err)
		this.Redirect("/article/articleList",302)
		return
	}

	var article models.Article
	article.Id = id

	o := orm.NewOrm()
	err = o.Read(&article)
	if err != nil{
		beego.Error("查询失败：",err)
		this.Redirect("/article/articleList",302)
		return
	}
	this.Data["select"] = this.GetString("select")
	this.Data["errmsg"] = errmsg
	this.Data["article"] = article

	Layout(this,"更新文章内容")
	this.TplName = "update.html"
}

func (this *ArticleController) HandleUpdate(){
	title := this.GetString("articleName")
	content := this.GetString("content")
	image := UploadFile(this,"uploadname")
	id,err := this.GetInt("id")

	if err != nil{
		this.Redirect("/article/articleList",302)
		return
	}

	if title == "" || content == "" || image == ""{
		errmsg := "内容不能为空"
		this.Redirect("/article/articleUpdate?errmsg="+errmsg+"&id="+strconv.Itoa(id),302)
		return
	}

	var article models.Article
	article.Id = id

	o := orm.NewOrm()
	err = o.Read(&article)
	if err != nil{
		errmsg := "不存在此文章"
		this.Redirect("/article/articleUpdate?errmsg="+errmsg+"&id="+strconv.Itoa(id),302)
		return
	}

	article.Title = title
	article.Content = content
	article.Image = image
	_,err = o.Update(&article)
	if err != nil{
		beego.Error("更新错误：",err)
		return
	}

	typename := this.GetString("select")
	this.Redirect("/article/articleList?select="+typename,302)
}

func UploadFile(this *ArticleController,name string) string {
	file,head,err := this.GetFile(name)

	if err != nil{
		return ""
	}
	defer file.Close()

	if head.Size > 500000{
		return ""
	}

	fileExt := path.Ext(head.Filename)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg"{
		return ""
	}

	fileName := time.Now().Format("20060102150405")+fileExt
	this.SaveToFile(name,"./static/image/"+fileName)
	return "/static/image/"+fileName
}

func (this *ArticleController) DelArticle(){
	id,err := this.GetInt("id")
	typename := this.GetString("select")
	if err != nil{
		beego.Error("获取id失败",err)
		this.Redirect("/article/articleList",302)
		return
	}

	var article models.Article
	article.Id = id

	o := orm.NewOrm()
	_,err = o.Delete(&article)
	if err != nil{
		beego.Error("数据库删除失败",err)
		this.Redirect("/article/articleList",302)
		return
	}

	this.Redirect("/article/articleList?select="+typename,302)
}

func (this *ArticleController) ShowArticleType(){
	articletypes := make([]models.ArticleType,0)

	o := orm.NewOrm()
	qs := o.QueryTable("ArticleType")
	qs.All(&articletypes)

	this.Data["articletypes"] = articletypes

	Layout(this,"编辑文章类型")
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Scripts"] = "scripts.html"
	this.TplName = "addType.html"
}

func (this *ArticleController) AddArticleType(){
	typename := this.GetString("typeName")
	if typename == ""{
		this.Redirect("/article/addarticletype",302)
		return
	}

	var articletype models.ArticleType
	articletype.TypeName = typename

	o := orm.NewOrm()
	_,err := o.Insert(&articletype)
	if err != nil{
		beego.Error("插入数据错误：",err)
		return
	}

	this.Redirect("/article/addarticletype",302)
}

func (this *ArticleController) DelArticleType(){
	id,err := this.GetInt("id")
	if err != nil{
		beego.Error(err)
		this.Redirect("/article/addarticletype",302)
		return
	}

	var articletype models.ArticleType
	articletype.Id = id

	o := orm.NewOrm()
	_,err = o.Delete(&articletype)
	if err != nil{
		beego.Error(err)
		this.Redirect("/article/addarticletype",302)
		return
	}

	this.Redirect("/article/addarticletype",302)
}

func Layout(this *ArticleController,title string){
	username := this.GetSession("username")
	this.Data["username"] = username.(string)
	this.Data["title"] = title
	this.Layout = "layout.html"
}