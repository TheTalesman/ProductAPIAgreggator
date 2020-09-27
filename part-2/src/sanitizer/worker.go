package main

import (
	"ProductAPIAgreggator/part-2/src/utils"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func generateWorkersAggregate(len int) {
	for i := 0; i < len; i++ {
		go workerAggregate()
	}

}

func setTransport() {
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 400
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 400
	myClient = &http.Client{Timeout: time.Second * 5}
}

//Recebe todas as linhas do mesmo pid, faz ate 3 requests ok e responde para o canal com linha agregada.
func workerAggregate() {
	var batch []ProductOutput

	for co := range ch {
		for _, o := range co.pl {
			var po ProductOutput
			po.ProductID = o.pid

			for _, p := range o.ps {
				url := p.Image
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
				resp.Body.Close()

				//Throw the body away  to deal with TIME_WAIT
				io.Copy(ioutil.Discard, resp.Body)

			}

			batch = append(batch, po)
			po.Images = nil
		}
		go workerRequester(batch)
		batch = nil

		co.wg.Done()

	}
}

func workerRequester(pos []ProductOutput) {
	log.Println("Requesting:", pos[0].ProductID, " to ", pos[len(pos)-1].ProductID)

	response, err := json.Marshal(pos)
	_, err = myClient.Post("http://localhost/products", "application/json", bytes.NewBuffer(response))
	utils.Check(err, 0)
	//	log.Println(resp)
}
