package main

import (
	"fmt"
	TMDb "savetrends.com/apiClient"
	MySQL "savetrends.com/sqlClient"
	"time"
)

func main() {
	apiClient := TMDb.APIClient{
		APIKey: "b97e33a6b0c4283466ad23df952ebd6a",
	}

	tn := time.Now().UTC()
	fmt.Println(tn)
	_,TMDb.WeekNr = tn.ISOWeek()

	sqlClient := MySQL.SQLClient{}
	sqlClient.EstablishConnectionToDB()



	//Zunächst werden die Film-Trends ermittelt (dazu werden erst die Trends abgerufen und dann nochmal die Infos zu jedem Film in den Trends einzeln)
	movieTrends := apiClient.GetMovieTrends()

	//Im folgenden werden zunächst die Referenz-Tabellen ergänzt
	persons := apiClient.GetPersonObjects(movieTrends)

	sqlClient.ExtendOrUpdatePersonTable(persons)

	genres := apiClient.GetGenres(movieTrends)

	sqlClient.ExtendGenresTable(genres)

	countries := apiClient.GetCountries(movieTrends)

	sqlClient.ExtendCountriesTable(countries)

	fmt.Println(movieTrends[1].WatchProviders)

	providers := apiClient.GetStreamingProvidersForMovieTrends(movieTrends)

	sqlClient.ExtendProviderTable(providers)

	//Nun werden die Tabellem Movies, Movie-Genre, ... ergänzt
	sqlClient.ExtendOrUpdateMovies(movieTrends)

	//TV-Trends
	seriesTrends := apiClient.GetTVTrends()

	//Referenz-Tabellen
	persons = apiClient.GetPersonObjectsTV(seriesTrends)

	sqlClient.ExtendOrUpdatePersonTable(persons)

	genres = apiClient.GetGenresTV(seriesTrends)

	sqlClient.ExtendGenresTable(genres)

	countries = apiClient.GetCountriesTV(seriesTrends)

	sqlClient.ExtendCountriesTable(countries)

	providers = apiClient.GetStreamingProvidersForTVTrends(seriesTrends)

	sqlClient.ExtendProviderTable(providers)

	networks := apiClient.GetNetworksForTVTrends(seriesTrends)

	sqlClient.ExtendNetworkTable(networks)

	defer sqlClient.DB.Close()
}