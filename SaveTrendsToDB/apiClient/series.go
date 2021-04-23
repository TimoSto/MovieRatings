package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Series struct {
	ID                 int     `json:id`
	Name               string  `json:name`
	Tagline            string  `json:tagline`
	Overview           string  `json:overview`
	Popularity           float64   `json:popularity`
	Number_of_seasons    int       `json:number_of_seasons`
	Number_of_episodes   int       `json:number_of_episodes`
	Poster_Path          string    `json:poster_path`
	First_air_date       string    `json:first_air_date`
	Last_air_date        string    `json:last_air_date`
	Vote_Average         float64   `json:vote_average`
	Vote_Count           float64   `json:vote_count`
	Genres               []Genre   `json:genres`
	In_Production        bool      `json:in_production`
	Networks             []Network `json:networks`
	Production_countries []Country `json:production_countires`
	Cast                 []PersonJob
	Crew                 []PersonJob
	WatchProviders		 ProviderResultSetDE
}

type TrendResultTV struct {
	Page          int `json:page`
	Results       []Series `json:results`
}

func(client *APIClient)GetTVTrends() []Series{
	fmt.Println("Retrieve tv-trend-information from TMDb-API...")
	var series []Series
	for i:=1 ; i<=5 ; i++ {
		fmt.Println(i)
		trends := client.GetTVTrendPage(i)
		for _, serie := range trends {
			series = append(series, client.GetSeriesByID(serie.ID))
		}
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
	for _, movie := range trends {
		series = append(series, client.GetSeriesByID(movie.ID))
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
	if len(credits.Crew) > 8 {
		series.Crew = credits.Crew[:8]
	} else {
		series.Crew = credits.Crew
	}
	if len(credits.Cast) > 8 {
		series.Cast = credits.Cast[:8]
	} else {
		series.Cast = credits.Cast
	}

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