package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type Params struct {
	Seed string
}

func main() {
	var param Params
	flag.StringVar(&param.Seed, "seed", "all", "data to seed ex. seed=user")
	flag.Parse()

	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	switch param.Seed {
	case "all":
		if err := utils.SeedUsers(db); err != nil {
			log.Fatal(err)
		}
		if err := utils.SeedCategory(db); err != nil {
			log.Fatal(err)
		}
		if err := utils.SeedMatchCategory(db); err != nil {
			log.Fatal(err)
		}
	case "users":
		if err := utils.SeedUsers(db); err != nil {
			log.Fatal(err)
		}
	case "category":
		if err := utils.SeedCategory(db); err != nil {
			log.Fatal(err)
		}
	case "match_category":
		if err := utils.SeedMatchCategory(db); err != nil {
			log.Fatal(err)
		}
	}

}
