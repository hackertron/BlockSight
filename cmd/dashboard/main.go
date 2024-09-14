package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/hackertron/blocksight/internal/data/retrieval"
	"github.com/joho/godotenv"
)

func main() {
	envFile, _ := godotenv.Read(".env")

	//alchemyAPIKey := os.Getenv("ALCHEMY_API_KEY")
	alchemyAPIKey := envFile["ALCHEMY_API_KEY"]
	if alchemyAPIKey == "" {
		log.Fatalf("ALCHEMY_API_KEY not set in environment") // Log the error and exit
	}
	r, err := retrieval.NewRetrieval(alchemyAPIKey)
	if err != nil {
		log.Fatalf("Failed to create retrieval instance: %v", err)
	}
	defer r.Close()

	blockNumber := big.NewInt(12345678)
	block, err := r.RetrieveBlockByNumber(blockNumber)
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}

	fmt.Printf("Block Number: %d\n", block.Number().Uint64())
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())

	// ... rest of the code ...
}
