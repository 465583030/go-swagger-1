# swaggen 

本项目旨在将 swagger.json 转化 为 service类的代码

使用方法

``` bash
go get github.com/inu1255/go-swagger/swaggen
swaggen -u="swagger.json的url或者文件路径" -t="template路径" -s="service导出路径" -ext="导出文件后缀(如:js)" 
# -e="可选:导出entity路径"
# 
# 示例
swaggen -u="http://localhost:8080/api/swagger.json" -t=$GOPATH/src/github.com/inu1255/go-swagger/swaggen/tmpl
cat service/test.js
```
将得到如下内容

``` js
import Model from './model.js'

function Test() {
}
Test.prototype = new Model("test")

//{//MainTestBody
//    name string //名字
//}
// 测试post
/*  */
Test.prototype.post = function(id,title,data) {
    return this.request("post/"+id+"?title="+title+"",data)
}

// 测试get
/*  */
Test.prototype.get = function() {
    return this.request("get")
}

export default new Test()
```
