package main

import (
	"log"
	"time"
	"github.com/kelseyhightower/envconfig"
	"github.com/ShivankSharma070/go-microservices-ecommerce/account"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if  err != nil {
		log.Fatal(err)
	}


	var r account.Repository
	retry.ForeverSleep(2 * time.Second,func(_ int)(err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL) 
		if err != nil {
			log.Println("Database connection error:",err)
		}

		return
	})

	defer r.Close()

	log.Println("Listening of port 8080...")
	s := account.NewService(r)

	log.Fatal(account.ListenGRPC(s, 8080))
}
