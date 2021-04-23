package main

import (
	"fmt"
	TMDb "savetrends.com/apiClient"
	MySQL "savetrends.com/sqlClient"
)

func main() {
	apiClient := TMDb.APIClient{
		APIKey: "b97e33a6b0c4283466ad23df952ebd6a",
	}

	sqlClient := MySQL.SQLClient{}
	sqlClient.EstablishConnectionToDB()

	//Zunächst werden die Film-Trends ermittelt (dazu werden erst die Trends abgerufen und dann nochmal die Infos zu jedem Film in den Trends einzeln)
	movieTrends := apiClient.GetMovieTrends()

	persons := apiClient.GetPersonObjects(movieTrends)

	sqlClient.ExtendOrUpdatePersonTable(persons)

	genres := apiClient.GetGenres(movieTrends)

	sqlClient.ExtendGenresTable(genres)

	countries := apiClient.GetCountries(movieTrends)

	sqlClient.ExtendCountriesTable(countries)

	fmt.Println(movieTrends[1].WatchProviders)

	providers := apiClient.GetStreamingProvidersForMovieTrends(movieTrends)

	sqlClient.ExtendProviderTable(providers)

	fmt.Println(providers)

	fmt.Println(movieTrends[1].WatchProviders)

	defer sqlClient.DB.Close()
}