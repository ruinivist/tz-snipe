package geodata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheDuration = 30 * 24 * time.Hour

func getCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "git-chronos", "geodata.json")
}

func loadFromCache() (TimezoneMap, error) {
	path := getCachePath()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cacheContainer cacheCountryData
	if err := json.NewDecoder(file).Decode(&cacheContainer); err != nil {
		return nil, err
	}

	if time.Since(cacheContainer.FetchedAt) > cacheDuration {
		return nil, fmt.Errorf("cache expired")
	}

	return cacheContainer.Data, nil
}

func saveToCache(tzMap TimezoneMap) error {
	path := getCachePath()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil { // creates if not exists
		return fmt.Errorf("failed to create cache dir")
	}

	file, err := os.Create(path) // create if not exists
	if err != nil {
		return err
	}

	defer file.Close()

	container := cacheCountryData{
		FetchedAt: time.Now(),
		Data:      tzMap,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ") // so it looks better
	return encoder.Encode(container)
}
