package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"./cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file containing the story")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	story, err := cyoa.JSONstory(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	h := cyoa.NewStoryHandler(story)

	fmt.Println("Server is starting at port 8080.")
	log.Fatal(http.ListenAndServe(":8080", h))
}
