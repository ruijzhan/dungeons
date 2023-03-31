package main

import (
	"log"

	"github.com/ruijzhan/dungeons/cmd/server/app"
)

func main() {
	server := app.NewApp()

	if err := server.Execute(); err != nil {
		log.Fatalln(err)
	}
}
