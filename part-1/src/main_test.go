package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	"testing"
)

func TestConnect(t *testing.T) {
	ok := dao.Connect()
	if !ok {
		t.Error("Falha na conexao")
	}
}
func TestUpsertFind(t *testing.T) {
	product := Product{
		ID:   1,
		Name: "Roupa do goku",
	}
	ok := Upsert(product)
	if !ok {
		t.Error("Falha na insercao")
	}

	//pp = Persisted Product
	ok, pp := Find(product)
	if !ok {
		t.Error("Falha na persistencia")
	}

	if product != pp {
		t.Error("Falha nos dados ")
	}

	return
}
