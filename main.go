package main

import (
	"github.com/gin-gonic/gin"
	swaggin "github.com/inu1255/go-swagger/gin"
)

type TestBody struct {
	Name string `json:"name,omitempty" gev:"名字"`
}

type TestData struct {
	Name  string `json:"name,omitempty" gev:"名字"`
	Id    string `json:"id,omitempty" gev:"id"`
	Title string `json:"title,omitempty" gev:"标题"`
}

func main() {
	app := gin.Default()
	test := app.Group("test")
	g := swaggin.NewRouter(test, "测试")
	g.Info("测试post").Body(
		new(TestBody),
	).Data(
		new(TestData),
	).Params(
		g.PathParam("id", "path param"),
		g.QueryParam("title", "query param"),
	).POST("/post/:id", func(c *gin.Context) {
		body := new(TestBody)
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		data := new(TestData)
		data.Id = c.Param("id")
		data.Title = c.Query("title")
		data.Name = body.Name
		c.JSON(200, data)
	})
	g.Info("测试get").GET("/get", func(c *gin.Context) {
		c.JSON(200, "hello world!")
	})
	swaggin.Swag.WriteJson("api/swagger.json")
	app.Static("api", "api")
	app.Run()
}
