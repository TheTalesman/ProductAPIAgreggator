package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	. "ProductAPIAgreggator/part-1/src/model"
	"log"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	_ "gopkg.in/mgo.v2/bson"
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

//Upsert a p Product as parameter.
//Returns OK flags success.
func Upsert(p Product) bool {
	pMap := structs.Map(p)
	ok, _ := dao.Upsert(pMap)
	return ok
}

//Find Receive a p Product as parameter
//Returns the relative product in database by ID. OK flags success.
func Find(p Product) (ok bool, pResponse Product) {

	ok, pInterface := dao.FindByID(p.ID)
	err := mapstructure.Decode(pInterface, &pResponse)
	if err != nil {
		log.Println(" Erro no Decode()")
		log.Fatal(err)
		ok = false
	}
	return
}
