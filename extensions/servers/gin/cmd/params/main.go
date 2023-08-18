package main

import (
	"github.com/gin-gonic/gin"
)

type Person struct {
	Id   string `form:"id" binding:"required,uuid"`
	Name string `form:"name" binding:"required" json:"name"`
}

func main() {
	route := gin.Default()
	route.GET("/data", func(c *gin.Context) {

		//println(fmt.Sprintf("%+v", c.Request.URL.Query()))
		var person Person

		//values := c.Request.URL.Query()
		//if err := mapForm(&person, values); err != nil {
		//	println(">>>>>>", err)
		//} else {
		//	println(">>>> >!> !>>! >! ")
		//}

		//e := c.ShouldBindWith(&person, binding.Query)
		//println("!!!! <<<<<< ", e.Error())
		if err := c.BindQuery(&person); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(200, gin.H{"name": person.Name})
	})
	route.Run(":8088")
}
