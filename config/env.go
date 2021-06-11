package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Getenvfile() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("file .env tidak ada")
	}
}
