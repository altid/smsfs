package main

import (
	"log"

	"github.com/alexdavid/sigma"
)

func main() {
	// Check if we're on a macOS
	// Run iMessage routine to fetch them
	cl, err := sigma.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()

	_, err = cl.Chats()
	if err != nil {
		log.Fatal(err)
	}
}
