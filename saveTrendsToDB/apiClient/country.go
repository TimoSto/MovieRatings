package apiClient

import (
	"fmt"
	"strings"
)

type Country struct {
	ISO_3166_1 string `json:iso_3166_1`
	Name       string `json:name`
}

func(client *APIClient)GetCountries(trends []Movie) []Country {
	var countries []Country
	fmt.Print("Analysing countries")
	for _,movie := range trends {
		fmt.Print(".")
		for _,country := range movie.Production_countries {
			if findCountryInSlice(countries, country.ISO_3166_1) == -1 {
				countries = append(countries, country)
			}
		}
	}
	return countries
}

func(client *APIClient)GetCountriesTV(trends []Series) []Country {
	var countries []Country
	fmt.Print("Analysing countries")
	for _,movie := range trends {
		fmt.Print(".")
		for _,country := range movie.Production_countries {
			if findCountryInSlice(countries, country.ISO_3166_1) == -1 {
				countries = append(countries, country)
			}
		}
	}
	return countries
}

func findCountryInSlice(arr []Country, id string) int {
	for i,c := range arr {
		if strings.Compare(c.ISO_3166_1, id) == 0 {
			return i
		}
	}
	return -1
}