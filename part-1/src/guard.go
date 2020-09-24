package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//Guard blocks duplicate POSTS in a 10 min range. Use hashes of body requests in REDIS as distributed persistence.
func Guard() gin.HandlerFunc {
	return func(c *gin.Context) {
		m := c.Request.Method
		if m == "POST" {

			var bodyBytes []byte
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)

			h := sha1.New()
			h.Write(bodyBytes)
			hash := hex.EncodeToString(h.Sum(nil))
			log.Println("hash: ", hash)

			tt, err := dao.RClient.Get(context.Background(), hash).Result()
			if err != nil {
				log.Println("Falha no GET REDIS")
				setExpire(hash)

			}

			tNow := time.Now().Unix()
			expire, err := strconv.ParseInt(tt, 10, 64)
			if err != nil {
				log.Println("Falha convertendo tempo da hash")

			}
			log.Println("tt:", tt)
			log.Println("tNow", tNow)
			log.Println("expire", expire)

			if expire > tNow {
				c.JSON(http.StatusForbidden, "Forbidden")
				c.Abort()
			} else {
				setExpire(hash)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			}

		}

		return
	}
}

func setExpire(hash string) {
	//600 = 10 min in unix
	err := dao.RClient.Set(context.Background(), hash, time.Now().Unix()+600, 0).Err()
	if err != nil {
		log.Fatal("Falha no SET REDIS")
		log.Fatal(err)
	}
}
