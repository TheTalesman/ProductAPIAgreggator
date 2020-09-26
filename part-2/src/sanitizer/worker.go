package main

import (
	"ProductAPIAgreggator/part-2/src/utils"
	"log"
	"net/http"
	"time"
)

func generateWorkers(len int) {
	for i := 0; i < len; i++ {
		go worker(ch)
	}
}

func setTransport() {
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 400
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 400
	myClient = &http.Client{Timeout: time.Second * 2}
}

//Recebe todas as linhas do mesmo pid, faz ate 3 requests ok e responde para o canal com linha agregada.
func worker(co <-chan ChannelObject) {
	for o := range co {

		var po ProductOutput
		po.ProductID = o.pid

		for _, p := range o.ps {
			url := p.Image
			//url = "http://localhost:8081"

			resp, err := myClient.Get(url)

			req++

			err = utils.Check(err, 0)
			if err == nil {

				if resp.StatusCode == 200 {
					po.Images = append(po.Images, p.Image)
					if len(po.Images) == 3 {
						break

					}
				}
			} else {
				//Caso ocorra erro de timeout ampliar tempo no metodo setTransport acima.
				log.Println("erro: ", err)
				log.Println("Caso ocorra erro de timeout ampliar tempo no metodo setTransport acima.")

			}
			//Throw the body away  to deal with TIME_WAIT
			//io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()

		}

		log.Println("Req Count: ", req)
		log.Println("Aggregate object: ", po)
		/* cor := ChannelObjectResponse{po}
		chr <- cor */

		o.wg.Done()
	}
}
