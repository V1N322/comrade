package main

import (
	lm "comrade/LM"
	"fmt"
	"os"

	"log"

	"github.com/joho/godotenv"

	"bufio"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}



	LLM := lm.NewComradeLM(os.Getenv("COMRADE_URL"), os.Getenv("COMRADE_TOKEN"), "MistralAI_Mixtral7b_Instruct")
	var result string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		scanner.Scan()
		input := scanner.Text()
		result, err = LLM.SendMessage(input)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}



	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}