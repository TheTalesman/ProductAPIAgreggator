package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	ok := dao.Connect()
	if !ok {
		log.Fatal("DB not connected")
	}
	log.Println("OK")
	r.GET("/products", FindProducts)
	r.POST("/products", UpsertProducts)
	log.Println("Listening...")
	r.Run(":8082")

}
