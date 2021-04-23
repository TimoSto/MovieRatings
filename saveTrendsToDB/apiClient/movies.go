package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Movie struct {
	ID                   int       `json:id`
	Title                string    `json:title`
	Overview             string    `json:overview`
	Popularity           float64   `json:popularity`
	Revenue              float64   `json:revenue`
	Poster_Path          string    `json:poster_path`
	Release_Date         string    `json:release_date`
	Vote_Average         float64   `json:vote_average`
	Vote_Count           float64   `json:vote_count`
	Genres               []Genre   `json:genres`
	Runtime              int       `json:runtime`
	Tagline              string    `json:tagline`
	Production_countries []Country `json:production_countires`
	Cast                 []PersonJob
	Crew                 []PersonJob
	WatchProviders		 ProviderResultSetDE
}

type TrendResultMovie struct {
	Page          int `json:page`
	Results       []Movie `json:results`
}

func(client *APIClient)GetMovieTrends() []Movie{
	fmt.Println("Retrieve movie-trend-information from TMDb-API...")
	var movies []Movie
	for i:=1 ; i<=5 ; i++ {
		fmt.Println(i)
		trends := client.GetMovieTrendPage(i)
		for _,movie := range trends {
			movies = append(movies, client.GetMovieByID(movie.ID))
		}
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
	movie.Title = strings.Replace(movie.Title, "'","\\'", -1)
	movie.Overview = strings.Replace(movie.Overview, "'","\\'", -1)
	movie.Tagline = strings.Replace(movie.Tagline, "'","\\'", -1)
	credits := client.GetCreditsForMovie(movie.ID)
	if len(credits.Crew) > 8 {
		movie.Crew = credits.Crew[:8]
	} else {
		movie.Crew = credits.Crew
	}
	if len(credits.Cast) > 8 {
		movie.Cast = credits.Cast[:8]
	} else {
		movie.Cast = credits.Cast
	}
	providers := client.GetStreamingProvidersForMovie(movie.ID)
	movie.WatchProviders = providers
	return movie
}