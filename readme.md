### 意见收集小工具
#### 背景

GO-Martini 实现的小 web 程序，支持表单提交，图片上传的基本模式

只需要一个二进制代码包，更轻，更好维护


#### 构建

建议 go 1.9.2 以上版本编译

#### 设计

1. 基本结构

	结合多个轻量好用工具
	- 使用了 [Golang Martini](https://github.com/go-martini/martini/blob/master/translations/README_zh_cn.md) 做 web
	- 中间件使用了 [Beego](https://beego.me/docs/intro/) 的 [log](https://beego.me/docs/module/logs.md) 模块和 [orm](https://beego.me/docs/mvc/model/overview.md) 模块
	- 配置使用了 [TOML](https://github.com/achun/tom-toml/blob/master/README_CN.md)
	- 其他
		- web 对外渲染模块使用了[martini-contrib/render](https://github.com/martini-contrib/render)


2. 建议

	自己编写构建脚本，自己增加控制脚本，静态目录根据自己的实际服务环境创建即可
