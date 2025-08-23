package main

import (
	"fmt"
	xhtml "golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"time"
)

var url = "https://crawler-test.com/mobile/separate_desktop_with_different_h1"

func main() {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	z := xhtml.NewTokenizer(response.Body)

	for {
		tt := z.Next()

		switch tt {
		case xhtml.ErrorToken:
			if z.Err() == io.EOF {
				return
			}
		case xhtml.DoctypeToken:
			fmt.Println("HTML Found")
			break
		case xhtml.StartTagToken, xhtml.SelfClosingTagToken:
			fmt.Printf("%s", z.Token().Data)
			break
		default:
			continue
		}
	}
}
