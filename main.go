package main

import (
	lm "comrade/LM"
	emb "comrade/embedding"

	"fmt"
	"os"

	"log"

	"github.com/joho/godotenv"

	"bufio"
)

func TestLLM() {
	LLM := lm.NewComradeLM(os.Getenv("COMRADE_URL"), os.Getenv("COMRADE_TOKEN"), "MistralAI_Mixtral7b_Instruct", true)
	var result string
	scanner := bufio.NewScanner(os.Stdin)

	var err error

	for {
		fmt.Print("You: ")
		scanner.Scan()
		input := scanner.Text()
		result, err := LLM.SendMessage(input)

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

func TestEmbedding() {
	embedding := emb.NewComradeEmbedding(os.Getenv("COMRADE_URL"), os.Getenv("COMRADE_TOKEN"), "Embeddings")
	result, err := embedding.EmbedText("apple")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	// TestLLM()

	TestEmbedding()


}