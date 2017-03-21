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