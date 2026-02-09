package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// now the nice thing is that I don't NEED to define a model always,
// go can be "loose" if you want it to be
// THis is unlike the geodata package where I made a model and all...

// Uses the user's public github events api, filters for push events and returns
// an array of .patch commits
// Eg https://api.github.com/users/xyz/events/public to
// https://api.github.com/repos/USER/REPO/commits/HASH
func getPushPatches(user string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", user)

	// another neat way is that you don't NEED a Client
	// it's good to have it for production for doing timeouts etc though
	// this uses the global shared default client
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github api failed")
	}

	var apiResult []map[string]any
	// do YOU know diff b/w Unmarshal and NewDecoder?
	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return nil, fmt.Errorf("failed to decode github response")
	}

	// in place filter
	end := 0
	for _, entry := range apiResult {
		if entry["type"].(string) == "PushEvent" {
			apiResult[end] = entry // THESE are NOT copies, do YOU know why?
			end++
		}
	}

	// build patch urls
	// do YOU know the diff b/w make vs declaring it using var
	var patches []string
	for i := 0; i < end; i++ {
		repoPrefix := apiResult[i]["repo"].(map[string]any)["url"].(string)
		commitHash := apiResult[i]["payload"].(map[string]any)["head"].(string)

		patchUrl := fmt.Sprintf("%s/commits/%s", repoPrefix, commitHash)
		patches = append(patches, patchUrl)
	}

	return patches, nil
}

func getUniqueTzFromPatches(patchUrls []string) ([]string, error) {
	var tzs []string

	for _, patchUrl := range patchUrls {

		req, err := http.NewRequest("GET", patchUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to build req")
		}

		// Refer: https://docs.github.com/en/rest/commits/commits?apiVersion=2022-11-28#get-a-commit
		// this is how YOU do custom stuff
		req.Header.Set("Accept", "application/vnd.github.patch")

		// can potentially run for INIFINITY as no timeout
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return nil, fmt.Errorf("failed to get patch")
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			// the DOCS say that some larges ones may time out
			// and they DID
			continue
		}
		// ^ THis is ugly

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read dpatch body")
		}

		// the date is the 3rd line
		lines := strings.Split(string(body), "\n")
		dateLine := lines[2]

		if !strings.HasPrefix(dateLine, "Date: ") {
			// my invariant failed and the world is ending
			panic("invariant of date on 3rd line failed")
		}

		// split by space and pick last
		spaceSplits := strings.Split(dateLine, " ")
		if len(spaceSplits) == 0 {
			panic("invariant of dateline having spaces failed")
		}
		tz := spaceSplits[len(spaceSplits)-1]

		tzs = append(tzs, tz)
	}

	// uniq it now
	// go seems to be getting a little TOO verbose for my taste now
	seen := make(map[string]int)
	end := 0
	for _, tz := range tzs {
		// the first is the value which is 0 if missing, second is the exists
		if _, exists := seen[tz]; !exists {
			seen[tz] = 1
			tzs[end] = tz
			end++
		}
	}

	return tzs[:end], nil
}
