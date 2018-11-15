package models

import ("github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int `orm:"pk;auto"`
	Username string `orm:"size(20);unique"`
	Pwd string `orm:"size(20)"`
	Articles []*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(100)"`
	Content string `orm:"size(500)"`
	Time time.Time `orm:"type(datetime);auto_now"`
	ReadCount int  `orm:"default(0)"`
	Image string `orm:"null"`
	ArticleType *ArticleType `orm:"rel(fk);null;on_delete(set_null)"`
	Users []*User `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(100);unique"`
	Articles []*Article `orm:"reverse(many)"`
}

func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsweb?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}