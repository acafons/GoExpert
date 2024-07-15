package main

import (
	"log"
	"os"
	"time"
)

type Message struct {
	API     string
	Address string
}

const HttpResponseTimeout = 1 * time.Second

func main() {
	for _, cep := range os.Args[1:] {
		log.Println("[info] Searching for address of CEP: " + cep)

		channelViaCep := make(chan Message)
		channelBrasilApi := make(chan Message)

		go GetAddressViaCEP(cep, channelViaCep)
		go GetAddressBrasilApi(cep, channelBrasilApi)

		select {
		case msg := <-channelViaCep:
			log.Println("[info] " + msg.API + ": " + msg.Address)
		case msg := <-channelBrasilApi:
			log.Println("[info] " + msg.API + ": " + msg.Address)
		case <-time.After(HttpResponseTimeout):
			log.Println("[error] Timeout exceeded for CEP: " + cep)
		}
	}
}
