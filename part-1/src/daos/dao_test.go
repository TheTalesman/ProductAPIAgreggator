package dao

import (
	"log"
	"testing"
)

func TestConnect(t *testing.T) {
	ok := Connect()
	if !ok {
		t.Error("Falha na conexao")
	}
	log.Println("Connect - sucesso")

}
