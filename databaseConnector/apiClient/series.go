package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func(client *APIClient)GetTVTrends() []Series{
	fmt.Println("Retrieve tv-trend-information from TMDb-API...")
	var series []Series
	for i:=1 ; i<=5 ; i++ {
		series = append(series, client.GetTVTrendPage(i)...)
	}
	return series
}

func(client *APIClient)GetTVTrendPage(n int) []Series{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/tv/popular?api_key=%v&page=%v", client.APIKey, n))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var trendingResultSet TrendResultTV
	err = json.Unmarshal(res,&trendingResultSet)
	if err != nil {
		panic(err)
	}

	return trendingResultSet.Results
}

func(client *APIClient)GetSeries(trends []Series) []Series {
	//In den Trends stehen nicht alle Attribute
	var series []Series
	for _, serie := range trends {
		series = append(series, client.GetSeriesByID(serie.ID))
	}
	return series
}

func(client *APIClient)GetSeriesByID(id int) Series {
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/tv/%v?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var series Series
	err = json.Unmarshal(res, &series)
	if err != nil {
		panic(err)
	}

	credits := client.GetCreditsForTV(series.ID)
	series.Cast = credits.Cast
	series.Crew = credits.Crew

	providers := client.GetStreamingProvidersForMovie(series.ID)
	series.WatchProviders = providers
	return series
}

func(client *APIClient)GetCreditsForTV(id int) CreditsForMovieOrTV {
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/tv/%v/credits?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var credits CreditsForMovieOrTV
	err = json.Unmarshal(res, &credits)
	if err != nil {
		panic(err)
	}
	return credits
}

func(client *APIClient)GetStreamingProvidersForSeries(id int) ProviderResultSetDE{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/tv/%v/watch/providers?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var resultSet StreamingProvidersResultSet
	err = json.Unmarshal(res, &resultSet)
	if err != nil {
		panic(err)
	}

	return resultSet.Results.DE
}