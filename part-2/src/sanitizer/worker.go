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
		go workerAggregate(ch)
	}
}

func setTransport() {
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 400
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 400
	myClient = &http.Client{Timeout: time.Second * 5}
}

//Recebe todas as linhas do mesmo pid, faz ate 3 requests ok e responde para o canal com linha agregada.
func workerAggregate(co <-chan ChannelObject) {

	for o := range co {

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

						go sendRequest(po)

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

		o.wg.Done()

	}

}

func sendRequest(po ProductOutput) {
	cor := ChannelObjectResponse{po}
	chr <- cor

}

func workerReq(chr <-chan ChannelObjectResponse) {
	for range chr {

		log.Println(len(chr))
	}
	if len(chr) >= 100 {

		var pos []ProductOutput
		for cor := range chr {
			pos = append(pos, cor.po)
		}
		workerRequester(pos)

	}
}
func workerRequester(pos []ProductOutput) {

	response, err := json.Marshal(pos)
	resp, err := myClient.Post("http://localhost:8082/products", "application/json", bytes.NewBuffer(response))
	utils.Check(err, 0)
	log.Println(resp)

	/*
		DEBUG
		var bd interface{}
		json.NewDecoder(resp.Body).Decode(&bd)
		log.Println("response Body:", bd) */
}
