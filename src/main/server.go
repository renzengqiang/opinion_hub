package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/achun/tom-toml"
	"os"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"main/handler"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"net/http"
	"os/signal"
	"syscall"
	"github.com/martini-contrib/binding"
	"main/models"
	"strconv"
	"mime/multipart"
	"io/ioutil"
)

var running bool = true

func registerSignalHandler() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		for {
			sig := <-c
			logs.Info("Signal %d received", sig)

			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				stop()
				time.Sleep(time.Second)
				os.Exit(0)
			}
		}
	}()
}

func stop() {
	running = false
}

type UploadForm struct {
	NickName    string  `form:"nickName"`
	Feeling 	string	`form:"feeling"`
	Score		string	`form:"score"`
	//TextUpload *multipart.FileHeader `form:"txtUpload"`
}

type UploadFormV2 struct {
	NickName    string  `form:"nickName"`
	Feeling 	string	`form:"feeling"`
	Score		string	`form:"score"`
	FileUpload *multipart.FileHeader `form:"fileUpload"`
}

func main() {
	// set runtime env
	pwd, _ := os.Getwd()
	execDir := flag.String("d", pwd, "execute directory")
	flag.Parse()
	fmt.Println("Current execute directory:", *execDir)
	os.Setenv("EXECDIR", *execDir)

	registerSignalHandler()


	// init conf
	conf, err := toml.LoadFile(fmt.Sprintf("%s/%s", os.Getenv("EXECDIR"), "conf/conf.toml"))
	if err != nil {
		fmt.Println(err)
		return
	}

	m := martini.Classic()
	m.Map(conf)
	fmt.Println("Run on: ", conf["basic.addr"].String())
	// init log
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s/logs/opinion_hub.log","separate":["error", "warning"]}`, os.Getenv("EXECDIR")))
	logger := logs.GetBeeLogger()
	m.Map(logger)
	// init db
	err = orm.RegisterDataBase("default", "mysql", fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8&autocommit=true`,
		conf["database_master.user"].String(),
		conf["database_master.pswd"].String(),
		conf["database_master.host"].String(),
		conf["database_master.port"].String(),
		conf["database_master.name"].String(),
		))
	if err != nil {
		fmt.Println(err)
		return
	}

	orm.SetMaxIdleConns("default", 30)
	orm.SetMaxOpenConns("default", 30)
	readOrm := orm.NewOrm()
	readOrm.Using("default")
	m.Map(readOrm)

	// midleware
	m.Use(render.Renderer(render.Options{
		Directory: "templates", // Specify what path to load the templates from.
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		//IndentJSON: true, // Output human readable JSON
		//IndentXML: true, // Output human readable XML
		//HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
	}))
	m.Use(handler.Logger())
	//m.Use(handler.Recovery())

	m.Group("/opinion", func(r martini.Router) {
		r.Get("/index", IndexV2)
		r.Get("/IndexV2", IndexV2)
		r.Post("/add", binding.Bind(UploadForm{}), Add)
		r.Post("/addV2", binding.MultipartForm(UploadFormV2{}), AddV2)
	})


	martini.SetENV(conf["basic.env"].String())
	m.RunOnAddr(conf["basic.addr"].String())
	m.Run()
}



func Index(r render.Render, param martini.Params, logger *logs.BeeLogger, o orm.Ormer, conf toml.Toml, req *http.Request) {
	ops, _ := models.FetchOpinion(o)
	var ret []string
	for _, o := range ops {
		createTime := time.Unix(o.CreateTime, 0).Format("2006-01-02 15:04:05")
		row := fmt.Sprintf("昵称：%s,    感受：%s,    评分：%d,    创建时间：%s", o.NickName, o.Feeling, o.Score, createTime)
		ret = append(ret, row)
	}
	r.HTML(200, "opinion/index", ret)
}

func IndexV2(r render.Render, param martini.Params, logger *logs.BeeLogger, o orm.Ormer, conf toml.Toml, req *http.Request) {
	ops, _ := models.FetchOpinion(o)
	type Row struct {
		NickName string
		Feeling string
		Score int64
		Img 	string
		CreateTime	string
	}
	var ret []*Row
	for _, o := range ops {
		createTime := time.Unix(o.CreateTime, 0).Format("2006-01-02 15:04:05")
		row := &Row{
			NickName: o.NickName,
			Feeling: o.Feeling,
			Score: o.Score,
			Img: o.Img,
			CreateTime: createTime,
		}
		ret = append(ret, row)
	}
	r.HTML(200, "opinion/index_v2", ret)
}

func Add(uf UploadForm, r render.Render, param martini.Params, logger *logs.BeeLogger, o orm.Ormer, conf toml.Toml, req *http.Request) {
	score, _ := strconv.ParseInt(uf.Score, 10, 64)
	if score > int64(10) {
		score = int64(10)
	}
	if score < 0 {
		score = 0
	}
	m := &models.Opinion{
		NickName: uf.NickName,
		Feeling: uf.Feeling,
		Score: score,
		CreateTime: time.Now().Unix(),
		LastUpdate: time.Now().Unix(),
	}
	_, err := o.InsertOrUpdate(m)
	if err != nil {
		r.HTML(200, "opinion/add/fail",  err.Error())
	} else {
		r.HTML(200, "opinion/add/success", "提交成功")
	}
}


func AddV2(uf UploadFormV2, r render.Render, param martini.Params, logger *logs.BeeLogger, o orm.Ormer, conf toml.Toml, req *http.Request) {
	score, _ := strconv.ParseInt(uf.Score, 10, 64)
	if score > int64(10) {
		score = int64(10)
	}
	if score < 0 {
		score = 0
	}
	file, err := uf.FileUpload.Open()
	defer file.Close()
	fileBytes, _ := ioutil.ReadAll(file)

	var fileName string
	if len(fileBytes)== 0 || uf.FileUpload.Filename == "" {
		fileName = ""
	} else {
		fileName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), uf.FileUpload.Filename)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", "./public/imgs", fileName), fileBytes, 0755)
	}

	m := &models.Opinion{
		NickName: uf.NickName,
		Feeling: uf.Feeling,
		Score: score,
		Img: fileName,
		CreateTime: time.Now().Unix(),
		LastUpdate: time.Now().Unix(),
	}
	_, err = o.InsertOrUpdate(m)
	if err != nil {
		r.HTML(200, "opinion/add/fail",  err.Error())
	} else {
		r.HTML(200, "opinion/add/success", "提交成功")
	}
}
