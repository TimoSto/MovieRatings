package apiClient

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