package geodata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func fetchLive() (TimezoneMap, error) {
	client := http.Client{Timeout: 10 * time.Second}

	// from: https://restcountries.com
	url := "https://restcountries.com/v3.1/all?fields=name,population,timezones"

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch country data")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch country data")
	}

	var apiData []countryApiReponse
	if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
		return nil, fmt.Errorf("failed to fetch ")
	}

	return buildTzMap(apiData), nil
}

func Fetch() (TimezoneMap, error) {
	// try to get from cache
	tzMap, err := loadFromCache()
	if err == nil {
		return tzMap, nil
	}

	// fetch live and save
	tzMap, err = fetchLive()
	if err != nil {
		return nil, err
	}

	if err = saveToCache(tzMap); err != nil {
		return nil, err
	}

	return tzMap, nil
}

func buildTzMap(data []countryApiReponse) TimezoneMap {
	tzMap := make(TimezoneMap)

	for _, countryRaw := range data {
		country := Country{
			Name:       countryRaw.Name.Common,
			Population: countryRaw.Population,
		}

		for _, rawTZ := range countryRaw.Timezones {
			rawTZ = strings.TrimPrefix(rawTZ, "UTC")
			rawTZ = strings.ReplaceAll(rawTZ, ":", "")

			if rawTZ == "" {
				rawTZ = "+0000" // api quirk
			}

			tzMap[rawTZ] = append(tzMap[rawTZ], country)
		}
	}

	return tzMap
}
