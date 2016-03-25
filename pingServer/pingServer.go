package main

import (
	"flag"
	"fmt"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/packet"
	"os"
	"time"

	"github.com/go-ndn/log"
)

var path = flag.String("p", "/ndn/ping", "Listen path")
var server = flag.String("s", ":6363", "Register nfd Server")
var keyPath = flag.String("k", "key/default.pri", "default key path")

func main() {
	flag.Parse()
	conn, err := packet.Dial("tcp", *server)
	if err != nil {
		log.Fatalln(err)
	}

	recv := make(chan *ndn.Interest)
	face := ndn.NewFace(conn, recv)
	defer face.Close()

	pem, err := os.Open(*keyPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer pem.Close()
	key, _ := ndn.DecodePrivateKey(pem)

	m := mux.New()

	m.HandleFunc(*path, func(w ndn.Sender, i *ndn.Interest) {
		w.SendData(&ndn.Data{
			Name:    i.Name,
			Content: []byte(time.Now().UTC().String()),
		})
		log.Println("Recv Interest ", i.Name.String())
	})

	fmt.Println("Successful listen on", *path)

	m.Run(face, recv, key)
}
