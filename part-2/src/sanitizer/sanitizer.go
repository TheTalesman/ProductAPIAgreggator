package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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
	var workerBatch []Product
	ch = make(chan ChannelObject)
	maxWorkers := 400
	generateWorkers(maxWorkers)
	req = 0

	//passar pids com suas linhas para canais
	for i < len(products)-1 {

		currentPid := products[i].ProductID
		nextPid := products[i+1].ProductID
		workerBatch = append(workerBatch, products[i])
		if currentPid != nextPid {
			wg.Add(1)

			log.Println("sent to batch: ", currentPid)
			co := ChannelObject{
				currentPid,
				workerBatch, &wg,
			}
			//workers fazem requests e compilam ate 3 imagens por produto (descartando requests acima da 3 OK)
			ch <- co

			workerBatch = nil
			//	time.Sleep(200 * time.Millisecond)
		}
		i++

	}
	wg.Wait()
	log.Println("Aggregation Done!")

	//monta request para o atualizador
	for i := range chr {
		fmt.Println(i)
	}
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
