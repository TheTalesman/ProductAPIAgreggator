package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ch chan ChannelObject
var myClient *http.Client
var chr chan ChannelObjectResponse

func main() {

	setTransport()
	f, err := os.Open("../dump/input-dump")
	defer f.Close()
	if err != nil {
		log.Fatal("Falha lendo arquivo")

	}

	scanner := bufio.NewScanner(f)

	var products []Product

	for scanner.Scan() {
		var p Product
		//reader := strings.NewReader(scanner.Text())
		err := json.Unmarshal([]byte(scanner.Text()), &p)
		//err = json.NewDecoder(reader).Decode(&p)
		if err != nil {
			log.Println("Erro")
		}

		products = append(products, p)
	}

	//ordenar pids
	pid := func(p1, p2 *Product) bool {
		pid1, pid2, _ := parseAsInt(p1, p2)

		return pid1 < pid2
	}
	By(pid).Sort(products)

	var wg sync.WaitGroup
	i := 0
	var productBatch []Product
	var workerBatch ChannelObject

	//FINE TUNE SANITIZER

	nWorkers := 5
	channelBuffer := 100
	batchSize := 10
	sleepTime := 200 * time.Millisecond

	ch = make(chan ChannelObject, channelBuffer)
	generateWorkersAggregate(nWorkers)

	req = 0
	//passar pids com suas linhas para canais
	//go func() {
	for i < len(products) {

		currentPid := products[i].ProductID

		productBatch = append(productBatch, products[i])

		nextPid := ""
		if i < len(products)-1 {
			nextPid = products[i+1].ProductID
		}

		if nextPid != currentPid {
			//finish one product batch
			co := ProductLine{
				currentPid,
				productBatch,
			}

			workerBatch.pl = append(workerBatch.pl, co)
			productBatch = nil
		}

		//Batches de x produtos enviados para agregar
		//workers farão requests e compilam ate 3 imagens por produto (descartando requests acima da 3 OK)
		if len(workerBatch.pl) >= batchSize {
			wg.Add(1)
			workerBatch.wg = &wg

			ch <- workerBatch
			for _, pl := range workerBatch.pl {
				log.Println("sent to batch: ", pl.pid)
			}
			time.Sleep(sleepTime)
			workerBatch.pl = nil
		} else {
			i++
		}
	}

	//Não perder possíveis que não entraram em um batch full
	wg.Add(1)
	workerBatch.wg = &wg

	ch <- workerBatch
	for _, pl := range workerBatch.pl {
		log.Println("sent to batch: ", pl.pid)
	}
	//}()
	close(ch)

	wg.Wait()
	log.Println("Aggregation Done!")
	log.Println("Posting done!")

	/*
		log.Println(f)
		data := ProductNodes{}
		_ = json.Unmarshal([]byte(f), &data)

		for i := 0; i < len(data.ProductNodes); i++ {
			fmt.Println("Product Id: ", data.ProductNodes[i].productID)
			fmt.Println("Quantity: ", data.ProductNodes[i].image)
		} */

}

var req int

func parseAsInt(p1, p2 *Product) (pid1 int, pid2 int, err error) {
	pidString := strings.Split(p1.ProductID, "pid")
	pidString2 := strings.Split(p2.ProductID, "pid")
	pid1, err = strconv.Atoi(pidString[1])
	pid2, err = strconv.Atoi(pidString2[1])

	if err != nil {
		log.Println("Error: ", err)
	}

	return
}
