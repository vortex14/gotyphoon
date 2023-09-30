package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type Data struct {
	QrTare *string `json:"qr_tare" binding:"required" form:"qr_tare"`
}

type RequestModel struct {
	Payload []Data `json:"payload" binding:"required,min=1,dive" form:"payload"`
}

func getData(c *gin.Context, model interface{}) {
	if err := c.ShouldBindBodyWith(model, binding.JSON); err != nil {
		println("!!!")
	}

}

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *gin.Context) {
		var json RequestModel

		var json2 RequestModel

		//if err := c.ShouldBindJSON(&json); err != nil {
		//
		//}

		getData(c, &json)

		getData(c, &json2)

		//c.Set("body", &json)

		//var json2 RequestModel
		//if err2 := c.Bind(&json2); err2 != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		//	return
		//}
		//

		//v, s := c.Get("body")

		println(fmt.Sprintf("%+v", json2))
		println(fmt.Sprintf("%+v", json))
		//if json.User != "manu" || json.Password != "123" {
		//	c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		//	return
		//}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding XML (
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<root>
	//		<user>manu</user>
	//		<password>123</password>
	//	</root>)
	router.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding a HTML form (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8086")
}
