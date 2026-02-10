package core

import (
	"math"
	"testing"
)

func TestGetStatsForTZs(t *testing.T) {
	preds, err := GetStatsForTZs([]string{"+0100", "+0000"})
	if err != nil {
		t.Fatal(err)
	}

	if len(preds) == 0 {
		t.Fatalf("failed to get predictions")
	}
}

func TestPredSum(t *testing.T) {
	preds, err := GetStatsForTZs([]string{"+0100"})
	if err != nil {
		t.Fatal(err)
	}

	// sum must be 1
	sum := 0.0
	for _, pred := range preds {
		sum += pred.Probability
	}

	if math.Abs(sum-1) > 1e-6 {
		t.Fatalf("Probabilities did not sum to 1. Got %.4f instead", sum)
	}
}
