package apiClient

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
	WatchProviders		 ProviderResultSetDE
}