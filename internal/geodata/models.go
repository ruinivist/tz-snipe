package geodata

// Country info exposed to other packages
type Country struct {
	Name       string
	Population int64
}

// Timezone to Country lookups
type TimezoneMap map[string][]Country

// internal model for https://restcountries.com/v3.1/all?fields=name,population,timezones
type countryApiReponse struct {
	Name struct {
		Common string `json:"common"` // it has other fields but I don't want them
	} `json:"name"`
	Population int64    `json:"population"`
	Timezones  []string `json:"timezones"` // example: [ "UTC+00:00" ]
}
