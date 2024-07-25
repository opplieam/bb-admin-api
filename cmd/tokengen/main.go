package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type Params struct {
	Duration time.Duration
	UserID   int
}

func main() {
	var param Params
	flag.DurationVar(&param.Duration, "duration", time.Hour*1, "expire duration of token generation")
	flag.IntVar(&param.UserID, "userid", 1, "user_id to store in token")
	flag.Parse()

	token, err := utils.GenerateToken(param.Duration, int32(param.UserID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("token:", token)
	fmt.Println("duration:", param.Duration)
	fmt.Println("userid:", param.UserID)
}
