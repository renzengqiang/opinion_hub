package models

import "github.com/astaxie/beego/orm"

type Opinion struct {
	Id            int64 `orm:"pk"`
	NickName     string
	Feeling 	string
	Score 		  int64
	Img 		string
	CreateTime    int64
	LastUpdate    int64
}


func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Opinion))
}

func (this *Opinion) TableName() string {
	return "op_opinion"
}

func FetchOpinion(o orm.Ormer) (Opinion []*Opinion, err error) {
	_, err = o.QueryTable("op_opinion").OrderBy("-id").Exclude("nick_name", "").Limit(200).All(&Opinion)
	return
}
