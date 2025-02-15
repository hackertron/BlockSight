package utils

import (
	"math"

	"github.com/hackertron/blocksight/internal/models"
)

// calculateMean calculates the mean of a slice of float64 numbers.
func CalculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculateStdDev calculates the standard deviation of a slice of float64 numbers.
func CalculateStdDev(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range values {
		sumSquares += (v - mean) * (v - mean)
	}
	return math.Sqrt(sumSquares / float64(len(values)))
}

// calculateAverageGasPrice computes the average gas price of transactions in a block.
func CalculateAverageGasPrice(block *models.Block) float64 {
	if len(block.Transactions) == 0 {
		return 0
	}

	sum := 0.0
	for _, tx := range block.Transactions {
		sum += float64(tx.GasPrice.Uint64())
	}
	return sum / float64(len(block.Transactions))
}
