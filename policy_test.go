package main

import (
	"log"
	"testing"
)

func TestSavePoilcy(t *testing.T) {
	pnumber, _ := save(&request{QuoteNumber: "1234567", ReceiptNumber: "1234567"})
	log.Println("pnumber ", pnumber)
}

func TestMarshallJson(t *testing.T) {
	r, err := marshalPolicy(`{"QuoteNumber": "12343456","ReceiptNumber": "1234345678"}`)
	log.Println(r.QuoteNumber, err)
}
