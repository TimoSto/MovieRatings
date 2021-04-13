package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	APIKey string
}

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
	Cast                 []Person
	Crew                 []Person
}

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
	Cast                 []Person
	Crew                 []Person
}

type Person struct {
	ID                   int     `json:id`
	Name                 string  `json:name`
	Birthday             string  `json:birthday`
	Deathday             string  `json:deathday`
	Known_for_department string  `json:known_for_department`
	Gender               int     `json:gender`
	Biography            string  `json:biography`
	Popularity           float64 `json:popularity`
	Profile_path         string  `json:profile_path`
	Job string `json:job`
	Character string `json:character`
}

type CreditsForMovieOrTV struct {
	ID int `json:id`
	Cast []Person `json:cast`
	Crew []Person `json:crew`
}

type Genre struct {
	ID   int    `json:id`
	Name string `json:name`
}

type Country struct {
	ISO_3166_1 string `json:iso_3166_1`
	Name       string `json:name`
}

type Network struct {
	Name           string `json:name`
	ID             int    `json:id`
	Logo_Path      string `json:logo_path`
	Origin_country string `json:origin_country`
}

type TrendResultMovie struct {
	Page          int `json:page`
	Results       []Movie `json:results`
}

type TrendResultTV struct {
	Page          int `json:page`
	Results       []Series `json:results`
}

type TrendResultPerson struct {
	Page          int `json:page`
	Results       []Person `json:results`
}

func(client *APIClient)GetMovieTrends() []Movie{
	fmt.Println("Retrieve movie-trend-information from TMDb-API...")
	var movies []Movie
	for i:=1 ; i<=5 ; i++ {
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
	series.Cast = credits.Cast
	series.Crew = credits.Crew
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

func(client *APIClient)GetPersonTrends() []Person{
	var persons []Person
	for i:=1 ; i <=5 ; i++ {
		persons = append(persons, client.GetPersonTrendPage(i)...)
	}
	return persons
}

func(client *APIClient)GetPersonTrendPage(n int) []Person{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/person/popular?api_key=%v&page=%v", client.APIKey, n))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var trendingResultSet TrendResultPerson
	err = json.Unmarshal(res,&trendingResultSet)
	if err != nil {
		panic(err)
	}

	return trendingResultSet.Results
}

func(client *APIClient)GetPersons(trends []Person) []Person {
	var persons []Person
	for _, trend := range trends {
		persons = append(persons, client.GetPersonByID(trend.ID))
	}
	return persons
}

func(client *APIClient)GetPersonByID(id int) Person {
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/person/%v?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var person Person
	err = json.Unmarshal(res, &person)
	if err != nil {
		panic(err)
	}

	return person
}