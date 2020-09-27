package main

import (
	"sort"
	"sync"
)

//ChannelObject is sent from sanitizer to worker with data to test and aggregate
type ChannelObject struct {
	pl []ProductLine
	wg *sync.WaitGroup
}

//ProductLine represents a line from dump file
type ProductLine struct {
	pid string
	ps  []Product
}

//ChannelObjectResponse goes from worker to sanitizer to be sent to updater
type ChannelObjectResponse struct {
	po []ProductOutput
}

//ProductOutput objeto transformados para ser enviado ao updater
type ProductOutput struct {
	ProductID string   `json:"productId"`
	Images    []string `json:"images"`
}

//Product express a line in the dump file
type Product struct {
	ProductID string `json:"productId"`
	Image     string `json:"image"`
}

//SORT METHODS INTERFACE

// productSorter joins a By function and a slice of Planets to be sorted.
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
