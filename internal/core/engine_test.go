package core

import "testing"

func TestGetStatsForTZs(t *testing.T) {
	preds, err := getStatsForTZs([]string{"+0100", "+0000"})
	if err != nil {
		t.Fatal(err)
	}

	if len(preds) == 0 {
		t.Fatalf("failed to get predictions")
	}
}
