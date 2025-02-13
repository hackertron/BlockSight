package main

import (
	"context"
	"log"

	"github.com/hackertron/blocksight/internal/client"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create Alchemy client
	alchemyClient, err := client.NewAlchemyClient()
	if err != nil {
		log.Fatalf("Failed to create Alchemy client: %v", err)
	}

	// Get latest block
	ctx := context.Background()
	block, err := alchemyClient.GetLatestBlock(ctx)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}

	log.Printf("Latest block number: %v", block.Number)
	log.Printf("Number of transactions: %d", len(block.Transactions))

	// Get receipt for the first transaction if any
	if len(block.Transactions) > 0 {
		receipt, err := alchemyClient.GetTransactionReceipt(ctx, block.Transactions[0].Hash)
		if err != nil {
			log.Fatalf("Failed to get transaction receipt: %v", err)
		}
		log.Printf("Transaction status: %d", receipt.Status)
	}
}
