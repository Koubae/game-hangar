package testings

import (
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../tests/.env.test")
	if err != nil {
		panic("Error loading .env.test file: " + err.Error())
	}
}
