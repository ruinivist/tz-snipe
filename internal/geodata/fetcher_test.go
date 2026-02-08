package geodata

import (
	"os"
	"testing"
)

func TestFetchLive(t *testing.T) {
	tzMap, err := fetchLive()

	if err != nil {
		t.Fatal(err)
	}

	testTzMap(tzMap, t)
}

func TestCache(t *testing.T) {
	// delete cache
	cachePath := getCachePath()
	if err := os.Remove(cachePath); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	_, err := Fetch() // this populates cache
	if err != nil {
		t.Fatal(err)
	}

	tzMap, err := loadFromCache()
	if err != nil {
		t.Fatal(err)
	}

	testTzMap(tzMap, t)
}

func testTzMap(tzMap TimezoneMap, t *testing.T) {
	entry0100 := tzMap["+0100"]
	if len(entry0100) == 0 {
		t.Fatalf("no countries found to timezone")
	}
}
