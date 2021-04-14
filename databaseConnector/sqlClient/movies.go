package sqlClient

import (
	"database/sql"
	"dbconn.com/apiClient"
	"fmt"
	"strconv"
	"strings"
)

func(client *SQLClient)MovieWithIdExists(id int) bool {
	sqlstr := fmt.Sprintf("SELECT Count(*) FROM Movies WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var count int
	err := row.Scan(&count)
	return err == nil && count != 0
}

func(client *SQLClient)ExtendOrUpdateMovieTable(movies []apiClient.Movie) {
	fmt.Println("Check for changes in movies and extend or update database accordingly...")
	for _,movie := range movies {
		movie.Title = strings.Replace(movie.Title,"'","\\'",-1)
		movie.Tagline = strings.Replace(movie.Tagline,"'","\\'",-1)
		movie.Overview = strings.Replace(movie.Overview,"'","\\'",-1)
		if client.MovieWithIdExists(movie.ID) {
			client.UpdateMovieEntry(movie)
		} else {
			client.CreateMovieEntry(movie)
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
		if client.GetGenreByID(genre.ID).ID.Valid == false {
			client.CreateGenreEntry(genre)
		}
		client.CreateMovieGenreEntry(movie.ID, genre.ID)
	}
	for _,country := range movie.Production_countries {
		if client.GetCountryByID(country.ISO_3166_1).ID.Valid == false {
			client.CreateCountryEntry(country)
		}
		client.CreateMovieCountryEntry(movie.ID, country.ISO_3166_1)
	}
	for i,person := range movie.Cast {
		if i > 5 {
			break
		}
		if client.GetPersonByID(person.ID).ID.Valid == false {
			client.CreatePersonEntry(person)
		}
		client.CreateMoviePersonEntry(movie.ID, person.ID, "Actor_"+person.Character)
	}
	for i,person := range movie.Crew {
		if i > 5 {
			break
		}
		if client.GetPersonByID(person.ID).ID.Valid == false {
			client.CreatePersonEntry(person)
		}
		client.CreateMoviePersonEntry(movie.ID, person.ID, person.Job)
	}
	for _,provider := range movie.WatchProviders.Buy {
		if client.GetProviderByID(provider.Provider_id).ID.Valid == false {
			client.CreateProviderEntry(apiClient.StreamingProvider(provider))
		}
		client.CreateMovieProviderEntry(movie.ID, provider.Provider_id, "buy")
	}
	for _,provider := range movie.WatchProviders.Rent {
		if client.GetProviderByID(provider.Provider_id).ID.Valid == false {
			client.CreateProviderEntry(apiClient.StreamingProvider(provider))
		}
		client.CreateMovieProviderEntry(movie.ID, provider.Provider_id, "rent")
	}
	for _,provider := range movie.WatchProviders.Flatrate {
		if client.GetProviderByID(provider.Provider_id).ID.Valid == false {
			client.CreateProviderEntry(apiClient.StreamingProvider(provider))
		}
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

func(client *SQLClient)GetMovieByID(id int) Movie{
	sqlstr := fmt.Sprintf("SELECT * FROM Movies WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var movie Movie
	err := row.Scan(&movie.ID, &movie.Title, &movie.Overview, &movie.Popularity, &movie.Revenue, &movie.PosterPath, &movie.ReleaseDate, &movie.VoteAVG, &movie.VoteCount, &movie.Runtime, &movie.Tagline)
	if err != nil {
		panic(err)
	}
	return movie
}

func(client *SQLClient)UpdateMovieEntry(movie apiClient.Movie) {
	sqlmovie := client.GetMovieByID(movie.ID)
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

	//Eintrag für Film in SQL-DB hinzufügen
	if different {
		fmt.Println("Update movie with id ",movie.ID)
		sqlstr := fmt.Sprintf("Update Movies SET title='%v', overview='%v', popularity=%v, releaseDate='%v', posterPath='%v', voteCount=%v, voteAvg=%v, revenue=%v, runtime=%v, tagline='%v' WHERE id=%v", movie.Title, movie.Overview, movie.Popularity, movie.Release_Date, movie.Poster_Path, movie.Vote_Count, movie.Vote_Average, movie.Revenue, movie.Runtime, movie.Tagline, movie.ID)
		_, err := client.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
	}
}

func(client *SQLClient) WriteMovieTrendsToSQL(trends []apiClient.Movie, week int) {
	fmt.Println("Write Movie-Trends to database...")
	for _, movie := range trends {

		client.CheckIfMovieTrendEntryExist(movie, week)
	}
}

func(client *SQLClient) CheckIfMovieTrendEntryExist(trend apiClient.Movie, weekNr int) {
	sqlstr := fmt.Sprintf("SELECT movieid FROM MovieWeekPopularity WHERE movieid=%v AND weekNr=%v", trend.ID, weekNr)
	row := client.DB.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		client.WriteMovieTrendToSQL(trend, weekNr)
	case nil:

	default:
		panic(err)
	}
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