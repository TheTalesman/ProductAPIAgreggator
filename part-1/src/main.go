package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	fixture := os.Args[1:]
	if fixture[0] == "true" {
		Fixture()
		return
	}

	dao.Connect()
	dao.ConnectRedis()
	r := gin.Default()
	r.Use(Guard())
	r.GET("/products", FindProducts)
	r.POST("/products", UpsertProducts)
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
