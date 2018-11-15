package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {
    beego.InsertFilter("/article/*",beego.BeforeExec,Filterfunc)
    beego.Router("/", &controllers.UserController{},"*:ShowLogin")
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleReg")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/article/articleList",&controllers.ArticleController{},"get:ShowArticle")
    beego.Router("/article/addarticle",&controllers.ArticleController{},"get:ShowAddarticle;post:HandlerAddarticle")
    beego.Router("/article/articleDetial",&controllers.ArticleController{},"get:ShowArticleDetial")
    beego.Router("/article/articleUpdate",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    beego.Router("/article/articleDelete",&controllers.ArticleController{},"get:DelArticle")
    beego.Router("/article/addarticletype",&controllers.ArticleController{},"get:ShowArticleType;post:AddArticleType")
    beego.Router("/article/delarticletype",&controllers.ArticleController{},"get:DelArticleType")
    beego.Router("/logout",&controllers.UserController{},"get:Logout")
}

func Filterfunc(ctx *context.Context){
    username := ctx.Input.Session("username")
    if username == nil{
        ctx.Redirect(302,"/login")
        return
    }
}