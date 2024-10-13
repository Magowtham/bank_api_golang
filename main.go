package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	//loading environment variables
	error := godotenv.Load(".env")

	if error != nil {
		log.Fatalln("Failed to load .env variable file , error occurred -> ", error.Error())
	}

	storage := NewDataBase()

	listenAddress := os.Getenv("LISTEN_ADDR")

	server := NewAPIServer(listenAddress, storage)

	server.Run()

}
