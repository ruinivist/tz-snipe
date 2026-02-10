package core

import (
	"cmp"
	"fmt"
	"slices"
	"tz-snipe/internal/geodata"
)

func GetStatsForTZs(tzs []string) ([]Prediction, error) {
	// assumes tzs at this stage are uniq
	// 1. get geodata
	var predictions []Prediction
	tzMap, err := geodata.Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to get geodata")
	}

	// 2. aggregate pop
	// summing over countries pop was an afterthough in impl here
	// the model from geodata should've been better

	// GO QUIRK
	// var init of slice and map are both nil
	// but append handles init in case of slice
	// while maps are more strict
	// so here I need a map
	seenCountry := make(map[string]struct{})
	// another go quirk: structs take up 0 space
	// hence above
	var totalPop int64
	for _, curTz := range tzs {
		for _, country := range tzMap[curTz] {
			if _, exists := seenCountry[country.Name]; !exists {
				totalPop += country.Population
				seenCountry[country.Name] = struct{}{}
			}
		}
	}

	// 3. build predictions
	clear(seenCountry)
	for _, tz := range tzs {
		countries := tzMap[tz]

		// for new countries
		for _, country := range countries {
			if _, existsCountry := seenCountry[country.Name]; existsCountry {
				continue
			}

			// I like this style of avoiding indents, but is it good?
			predictions = append(predictions, Prediction{
				Country:     country.Name,
				Probability: float64(country.Population) / float64(totalPop),
			})
		}
	}

	// 4. sort
	slices.SortFunc(predictions, func(a, b Prediction) int {
		return cmp.Compare(b.Probability, a.Probability) // desc
	})

	return predictions, nil
}
