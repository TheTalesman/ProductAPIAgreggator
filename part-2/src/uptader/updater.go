package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	dao.Connect()

	r.POST("/products", UpsertProducts)
	r.Run(":8082") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
