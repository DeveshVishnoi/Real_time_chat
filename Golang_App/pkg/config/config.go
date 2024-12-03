package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvData() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error getting loading the env file : ", err)
		return
	}

	fmt.Println("Successfully Load the Env Data")

}
