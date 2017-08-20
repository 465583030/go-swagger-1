package main

import (
	"github.com/gin-gonic/gin"
	"github.com/inu1255/go-swagger/swaggin"
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
	app := swaggin.New()
	app.Info("测试post")
	app.Body(new(TestBody))
	app.Data(new(TestData))
	app.PathParam("id", "path param")
	app.QueryParam("title", "query param")
	app.POST("/post/:id", func(c *gin.Context) {
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
	app.Info("测试get").GET("/get", func(c *gin.Context) {
		c.JSON(200, "hello world!")
	})
	app.Swagger("/api")
	app.Run()
}
