package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	m "ProductAPIAgreggator/part-1/src/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

var PS []m.Product

func worker(id int, wg *sync.WaitGroup, final bool, f *os.File) {
	factor := 8000
	init := id * factor
	for i := init; i < init+factor; i++ {
		var p m.Product
		p.ID = strconv.Itoa(i)
		p.Name = "Nome-Teste"
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(b)

		v := []byte(",")
		b = append(b, v[0])
		f.Write(b)

	}

	wg.Done()
}

func Fixture() {
	dao.Connect()
	//rbf stands for REALLY BIG FILE. Was thinking about BFF, but nevermind. Any chance you are a DOOM player?
	f, err := os.Create("rbf.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(f, "[")

	var wg sync.WaitGroup

	for i := 0; i <= 7; i++ {
		wg.Add(1)
		go worker(i, &wg, false, f)
	}
	wg.Wait()
	log.Println("1/4 - 25%")
	for i := 8; i <= 15; i++ {
		wg.Add(1)
		go worker(i, &wg, false, f)
	}
	wg.Wait()
	log.Println("2/4- 50%")
	for i := 16; i <= 23; i++ {
		wg.Add(1)
		go worker(i, &wg, false, f)
	}

	wg.Wait()
	log.Println("3/4 - 75%")
	for i := 24; i <= 32; i++ {
		wg.Add(1)
		go worker(i, &wg, false, f)
	}

	wg.Wait()
	log.Println("4/4 - 100%")
	var p m.Product
	p.ID = "99999999999"
	p.Name = "Nome-Teste Final"
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Write(b)
	fmt.Fprintln(f, "]")

}
