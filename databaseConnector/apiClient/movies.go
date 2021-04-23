package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func(client *APIClient)GetMovieTrends() []Movie{
	fmt.Println("Retrieve movie-trend-information from TMDb-API...")
	var movies []Movie
	for i:=1 ; i<=1 ; i++ {
		movies = append(movies, client.GetMovieTrendPage(i)...)
	}
	return movies
}

func(client *APIClient)GetMovieTrendPage(n int) []Movie{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%v&page=%v", client.APIKey, n))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var trendingResultSet TrendResultMovie
	err = json.Unmarshal(res,&trendingResultSet)
	if err != nil {
		panic(err)
	}

	return trendingResultSet.Results
}

func(client *APIClient)GetMovies(trends []Movie) []Movie {
	//In den Trends stehen nicht alle Attribute
	var movies []Movie
	for _, movie := range trends {
		movies = append(movies, client.GetMovieByID(movie.ID))
	}
	return movies
}

func(client *APIClient)GetMovieByID(id int) Movie {
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/%v?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var movie Movie
	err = json.Unmarshal(res, &movie)
	if err != nil {
		panic(err)
	}

	credits := client.GetCreditsForMovie(movie.ID)
	movie.Crew = credits.Crew
	movie.Cast = credits.Cast
	providers := client.GetStreamingProvidersForMovie(movie.ID)
	movie.WatchProviders = providers
	return movie
}

func(client *APIClient)GetCreditsForMovie(id int) CreditsForMovieOrTV {
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/%v/credits?api_key=%v", id, client.APIKey))
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

func(client *APIClient)GetStreamingProvidersForMovie(id int) ProviderResultSetDE{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/%v/watch/providers?api_key=%v", id, client.APIKey))
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