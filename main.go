package main

import (
	"fmt"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	h := host{
		apiUrl:     "https://kesdev.com",
		contentKey: os.Getenv("CONTENT_KEY"),
	}
	c := newCache()
	resp, err := getPosts(c, h, 1)
	check(err)
	fmt.Println(resp)
	resp, err = getPosts(c, h, 1)
	check(err)
	fmt.Println(resp)
}
