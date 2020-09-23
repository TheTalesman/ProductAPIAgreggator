package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	. "ProductAPIAgreggator/part-1/src/model"
	"log"
	"testing"
)

var ProductTest Product
var ObjID string

func TestUpsertFind(t *testing.T) {
	dao.Connect()
	initTestObj(t)
	testInsertFind(t)
	testFindUpsert(t)

	return
}

func testFindUpsert(t *testing.T) {
	log.Println("TESTING testFindUpsert...")
	ok, pp := Find(ProductTest)
	if !ok {
		t.Error("Falha no find")
	}
	newName := "Roupa do Vegeta"
	pp.Name = newName
	ok = Upsert(pp)
	if !ok {
		t.Error("Falha no upsert")
	}

	ok, rpp := Find(ProductTest)
	if !ok {
		t.Error("Falha no find")
	}

	if rpp.Name != pp.Name && rpp.ID != pp.ID {
		if !ok {
			t.Error("Falha nos dados")
		}
	}
	if rpp.Name != newName {
		if !ok {
			t.Error("Falha nos dados NOME")
		}
	}
	log.Println("UPSERT - sucesso")
	log.Println("TESTING testFindUpsert... SUCCESS!")
}

func testInsertFind(t *testing.T) {
	log.Println("TESTING testInsertFind...")
	ok := Upsert(ProductTest)
	if !ok {
		t.Error("Falha na insercao")
	}
	log.Println("Insercao - sucesso")
	//pp = Persisted Product
	ok, pp := Find(ProductTest)
	if !ok {
		t.Error("Falha na persistencia")
	}
	log.Println("Persistencia - sucesso")
	if ProductTest.ID != pp.ID {
		t.Error("Falha nos dados ")
	}
	log.Println("Dados - sucesso")
	log.Println("TESTING testInsertFind... SUCESS!")
}

func initTestObj(t *testing.T) {
	ObjID = "123"

	ProductTest = Product{
		ID:   ObjID,
		Name: "Roupa do goku",
	}

}
