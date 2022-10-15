package main

import (
	"bytes"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/leonardom/go-url-shortener/shortener"
)

func httpPort() string {
	port := "8000"
	if os.Getenv("POST") != "" {
		port = os.Getenv("POST")
	}
	return fmt.Sprintf(":%s", port)
}

func main()  {
	address := fmt.Sprintf("http://localhost%s", httpPort())
	redirect := shortener.Redirect{}
	redirect.URL = "https://www.youtube.com/watch?v=QyBXz9SpPqE"

	body, err := msgpack.Marshal(&redirect)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(address, "application/x-msgpack", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	msgpack.Unmarshal(body, &redirect)
	log.Printf("%v\n", redirect)
}