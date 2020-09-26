package main

import (
	"ProductAPIAgreggator/part-2/src/utils"
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ch chan ChannelObject
var myClient *http.Client

// planetSorter joins a By function and a slice of Planets to be sorted.
type productSorter struct {
	products []Product
	by       func(p1, p2 *Product) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *productSorter) Len() int {
	return len(s.products)
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *Product) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(products []Product) {
	ps := &productSorter{
		products: products,
		by:       by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// Swap is part of sort.Interface.
func (s *productSorter) Swap(i, j int) {
	s.products[i], s.products[j] = s.products[j], s.products[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *productSorter) Less(i, j int) bool {
	return s.by(&s.products[i], &s.products[j])
}

type Product struct {
	ProductID string `json:"productId"`
	Image     string `json:"image"`
}

func generateWorkers(len int) {
	for i := 0; i < len; i++ {
		go worker(ch)
	}
}

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

	/* 	for _, p := range products {
		fmt.Println("By pid:", p.ProductID)

	} */

	var wg sync.WaitGroup
	i := 0
	var workerBatch []Product
	ch = make(chan ChannelObject)
	maxWorkers := 400
	generateWorkers(maxWorkers)
	req = 0
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
			ch <- co

			workerBatch = nil
			//	time.Sleep(200 * time.Millisecond)
		}
		i++

	}
	wg.Wait()
	log.Println("Aggregation Done!")
	//passar pids com suas linhas para canais

	//canais fazem requests e compilam ate 3 imagens por produto (descartando requests acima da 3 OK)

	//monta request para o atualizador

	/*
		log.Println(f)
		data := ProductNodes{}
		_ = json.Unmarshal([]byte(f), &data)

		for i := 0; i < len(data.ProductNodes); i++ {
			fmt.Println("Product Id: ", data.ProductNodes[i].productID)
			fmt.Println("Quantity: ", data.ProductNodes[i].image)
		} */

}

type ChannelObject struct {
	pid string
	ps  []Product
	wg  *sync.WaitGroup
}

//ProductOutput objeto transformados para ser enviado ao updater
type ProductOutput struct {
	ProductID string   `json:"productId"`
	Images    []string `json:"images"`
}

var req int

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
		//We don't wanna kill the api crazy gopher...

		o.wg.Done()
	}
}

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
