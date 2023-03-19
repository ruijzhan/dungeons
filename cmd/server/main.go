package main

//import cmd/server/app
import (
	"log"

	"github.com/ruijzhan/dungeons/cmd/server/app"
)

func main() {
	server := app.NewServer()

	if err := server.Execute(); err != nil {
		log.Fatalln(err)
	}
}
