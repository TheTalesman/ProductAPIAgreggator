package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func Upsert(p Product) bool {
	pMap := structs.Map(p)
	ok, _ := dao.Upsert(pMap)
	return ok
}

func Find(p Product) (bool, Product) {
	return false, Product{"", 0}
}
