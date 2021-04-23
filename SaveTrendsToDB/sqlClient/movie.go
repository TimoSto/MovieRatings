package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
	"strconv"
	"strings"
)

type Movie struct {
	ID          int             `json:id`
	Title       sql.NullString  `json:title`
	Overview    sql.NullString  `json:overview`
	Popularity  sql.NullFloat64 `json:popularity`
	Revenue     sql.NullString  `json:revenue`
	PosterPath  sql.NullString  `json:posterPath`
	ReleaseDate sql.NullString  `json:release_date`
	VoteAVG     sql.NullFloat64 `json:voteAvg`
	VoteCount   sql.NullInt64   `json:VoteCount`
	Runtime     sql.NullInt64   `json:runtime`
	Tagline     sql.NullString  `json:tagline`
}

func(client *SQLClient)ExtendOrUpdateMovies(movies []apiClient.Movie) {

	for _, movie := range movies {
		sqlmovie,n := client.GetMovieByID(movie.ID)

		if n == -1 {
			client.CreateMovieEntry(movie)
		} else {
			client.UpdateMovieEntry(movie, sqlmovie)
		}

		if client.MovieTrendExists(movie.ID, apiClient.WeekNr) == false {
			client.WriteMovieTrendToSQL(movie, apiClient.WeekNr)
		}
	}
}

func(client *SQLClient)CreateMovieEntry(movie apiClient.Movie) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create MovieEntry "+movie.Title)
	sqlstr := fmt.Sprintf("INSERT INTO Movies(id, title, overview, popularity, releaseDate, posterPath, voteCount, voteAvg, revenue, runtime, tagline) VALUES(%v,'%v','%v',%v,'%v','%v',%v, %v,'%v', %v, '%v')",movie.ID, movie.Title, movie.Overview, movie.Popularity, movie.Release_Date, movie.Poster_Path, movie.Vote_Count, movie.Vote_Average, movie.Revenue, movie.Runtime, movie.Tagline)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
	for _,genre := range movie.Genres {

		client.CreateMovieGenreEntry(movie.ID, genre.ID)
	}
	for _,country := range movie.Production_countries {

		client.CreateMovieCountryEntry(movie.ID, country.ISO_3166_1)
	}
	for i,person := range movie.Cast {
		if i > 5 {
			break
		}
		client.CreateMoviePersonEntry(movie.ID, person.ID, "Actor_"+person.Character)
	}
	for i,person := range movie.Crew {
		if i > 5 {
			break
		}

		client.CreateMoviePersonEntry(movie.ID, person.ID, person.Job)
	}
	for _,provider := range movie.WatchProviders.Buy {

		client.CreateMovieProviderEntry(movie.ID, provider.Provider_id, "buy")
	}
	for _,provider := range movie.WatchProviders.Rent {

		client.CreateMovieProviderEntry(movie.ID, provider.Provider_id, "rent")
	}
	for _,provider := range movie.WatchProviders.Flatrate {

		client.CreateMovieProviderEntry(movie.ID, provider.Provider_id, "flat")
	}
}

func(client *SQLClient)CreateMovieGenreEntry(movie int, genre int) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create MovieGenreEntry", movie, genre)
	sqlstr := fmt.Sprintf("INSERT INTO MovieGenre(movieId, genreId) VALUES(%v,'%v')",movie, genre)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)CreateMovieCountryEntry(movie int, country string) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create MovieCountry-Entry", movie, country)
	sqlstr := fmt.Sprintf("INSERT INTO MovieCountry(movieId, countryId) VALUES(%v,'%v')",movie, country)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)CreateMoviePersonEntry(movie int, person int, job string) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create MovieCreditsEntry", movie, person, job)
	job = strings.Replace(job, "'", "\\'", -1)
	sqlstr := fmt.Sprintf("INSERT INTO MovieCredits(movieId, personId, job) VALUES(%v,%v,'%v')",movie, person, job)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)CreateMovieProviderEntry(movie int, provider int, service string) {
	fmt.Println("Create MovieProviderEntry", movie, provider)
	service = strings.Replace(service, "'", "\\'", -1)
	sqlstr := fmt.Sprintf("INSERT INTO MovieProvider(movieId, provider, service) VALUES(%v,%v,'%v')", movie, provider, service)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)GetMovieByID(id int) (Movie, int){
	sqlstr := fmt.Sprintf("SELECT * FROM Movies WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var movie Movie
	err := row.Scan(&movie.ID, &movie.Title, &movie.Overview, &movie.Popularity, &movie.Revenue, &movie.PosterPath, &movie.ReleaseDate, &movie.VoteAVG, &movie.VoteCount, &movie.Runtime, &movie.Tagline)
	if err == sql.ErrNoRows {
		return Movie{}, -1
	}
	if err != nil {
		panic(err)
	}
	return movie, 1
}

func(client *SQLClient)UpdateMovieEntry(movie apiClient.Movie, sqlmovie Movie) {

	revenue, err := strconv.ParseFloat(sqlmovie.Revenue.String, 64)
	if err != nil {
		panic(err)
	}
	sqlmovie.Title.String = strings.Replace(sqlmovie.Title.String,"'","\\'",-1)
	sqlmovie.Tagline.String = strings.Replace(sqlmovie.Tagline.String,"'","\\'",-1)
	sqlmovie.Overview.String = strings.Replace(sqlmovie.Overview.String,"'","\\'",-1)

	different := sqlmovie.Popularity.Float64 != movie.Popularity ||
		sqlmovie.VoteCount.Int64 != int64(movie.Vote_Count) ||
		sqlmovie.VoteAVG.Float64 != movie.Vote_Average ||
		revenue != movie.Revenue ||
		strings.Compare(sqlmovie.Title.String, movie.Title) != 0 ||
		strings.Compare(sqlmovie.Overview.String, movie.Overview) != 0 ||
		strings.Compare(sqlmovie.Tagline.String, movie.Tagline) != 0 ||
		strings.Compare(sqlmovie.PosterPath.String, movie.Poster_Path) != 0

	//Eintrag für Film in SQL-DB aktualisieren
	if different {
		fmt.Println("Update movie with id ",movie.ID)
		sqlstr := fmt.Sprintf("Update Movies SET title='%v', overview='%v', popularity=%v, releaseDate='%v', posterPath='%v', voteCount=%v, voteAvg=%v, revenue=%v, runtime=%v, tagline='%v' WHERE id=%v", movie.Title, movie.Overview, movie.Popularity, movie.Release_Date, movie.Poster_Path, movie.Vote_Count, movie.Vote_Average, movie.Revenue, movie.Runtime, movie.Tagline, movie.ID)
		_, err := client.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
	}

	//for _,genre := range movie.Genres {
	//	if client.GetGenreByID(genre.ID).ID.Valid == false {
	//		client.CreateGenreEntry(genre)
	//	}
	//
	//	if client.MovieGenreEntryExists(movie.ID, genre.ID) == false {
	//		client.CreateMovieGenreEntry(movie.ID, genre.ID)
	//	}
	//}
	//
	//for _, person := range movie.Cast {
	//	if client.GetPersonByID(person.ID).ID.Valid == false {
	//		client.CreatePersonEntry(person)
	//	} else {
	//		client.UpdatePersonEntry(person)
	//	}
	//}
}

func(client *SQLClient) WriteMovieTrendToSQL(movie apiClient.Movie, weekNr int) {
	fmt.Println("Create MovieTrend Entry", movie.Title, weekNr)
	sql := fmt.Sprintf("INSERT INTO MovieWeekPopularity(movieId, weekNr, popularity, voteAVG, voteCount) VALUES ('%v', %v, %v, %v, %v)",movie.ID, weekNr, movie.Popularity, movie.Vote_Average, movie.Vote_Count)
	_, err := client.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	WriteSQLToFile(sql)
}

func(client *SQLClient) MovieTrendExists(movieID int, weekNr int) bool{
	sqlstr := fmt.Sprintf("SELECT movieid FROM MovieWeekPopularity WHERE movieid=%v AND weekNr=%v", movieID, weekNr)
	row := client.DB.QueryRow(sqlstr)
	var found_id int
	err := row.Scan(&found_id)
	if err == sql.ErrNoRows{
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}