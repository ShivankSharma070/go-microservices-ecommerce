package main

import (
	"log"
	"time"

	"github.com/ShivankSharma070/go-microservices-ecommerce/orders"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_SERVICE_URL"`
	CatalogURL  string `envconfig:"DATABASE_SERVICE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r orders.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = orders.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("Unable to connect to db: ",err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port 8080...")
	s := orders.NewService(r)
	log.Fatal(orders.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
