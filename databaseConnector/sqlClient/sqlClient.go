package sqlClient

import (
	"database/sql"
	"dbconn.com/apiClient"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
)

type SQLClient struct {
	DB *sql.DB
}

func(client *SQLClient)Exec(sqlstr string) (sql.Result,error) {
	res, err := client.DB.Exec(sqlstr)
	if err == nil {
		WriteSQLToFile(sqlstr)
	}
	return res, err
}

func WriteSQLToFile(sqlstr string){
	f, err := os.OpenFile("../database/FILLDB.sql",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//sqlstr = strings.Replace(sqlstr,"'","\\'",-1)
	if _, err := f.WriteString(sqlstr +";\n"); err != nil {
		panic(err)
	}
}

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

type Series struct {
	ID           int         `json:id`
	Title        sql.NullString  `json:title`
	Overview    sql.NullString  `json:overview`
	Popularity  sql.NullFloat64 `json:popularity`
	Seasons     sql.NullInt64   `json:seasons`
	Episodes    sql.NullInt64   `json:episodes`
	Poster      sql.NullString  `json:posterPath`
	ReleaseDate sql.NullString  `json:release_date`
	VoteAVG     sql.NullFloat64 `json:voteAvg`
	VoteCount   sql.NullInt64   `json:VoteCount`
	FirstAit    sql.NullString  `json:firstAir`
	LastAir     sql.NullString  `json:lastAir`
	Tagline     sql.NullString  `json:tagline`
}

type Country struct {
	ID sql.NullString `json:id`
	CName sql.NullString `json:cname`
}

type Network struct {
	ID            sql.NullInt64  `json:id`
	NName         sql.NullString `json:nname`
	Logo          sql.NullString `json:logo`
	OriginCountry sql.NullString `json:originCountry`
}

type Genre struct {
	ID sql.NullInt64 `json:id`
	Genre sql.NullString `json:genre`

}

func(client *SQLClient)EstablishConnectionToDB() {
	fmt.Println("Trying to connect to to mySQL-DB...")
	var err error
	client.DB, err = sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	if err != nil {
		panic(err)
	}
}

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
			fmt.Println("Create")
			client.CreateMovieEntry(movie)
		}
	}
}

func(client *SQLClient)CreateMovieEntry(movie apiClient.Movie) {
	//Eintrag für Film in SQL-DB hinzufügen
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
}

func(client *SQLClient)CreateGenreEntry(genre apiClient.Genre) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create GenreEntry "+genre.Name)
	sqlstr := fmt.Sprintf("INSERT INTO Genres(id, genre) VALUES(%v,'%v')",genre.ID, genre.Name)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)CreateMovieGenreEntry(movie int, genre int) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create MovieGenreEntry ")
	sqlstr := fmt.Sprintf("INSERT INTO MovieGenre(movieId, genreId) VALUES(%v,'%v')",movie, genre)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)GetGenreByID(id int) Genre{
	sqlstr := fmt.Sprintf("SELECT * FROM Genres WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var genre Genre
	err := row.Scan(&genre.ID, &genre.Genre)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return genre
}

func(client *SQLClient)GetCountryByID(id string) Country{
	sqlstr := fmt.Sprintf("SELECT * FROM Countries WHERE id='%v'", id)
	row := client.DB.QueryRow(sqlstr)
	var country Country
	err := row.Scan(&country.ID, &country.CName)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return country
}

func(client *SQLClient)CreateCountryEntry(country apiClient.Country) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create CountryEntry "+country.Name)
	sqlstr := fmt.Sprintf("INSERT INTO Countries(id, cname) VALUES('%v','%v')",country.ISO_3166_1, country.Name)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)CreateMovieCountryEntry(movie int, country string) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create MovieCountryEntry ")
	sqlstr := fmt.Sprintf("INSERT INTO MovieCountry(movieId, countryId) VALUES(%v,'%v')",movie, country)
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

func(client *SQLClient)SeriesWithIdExists(id int) bool {
	sqlstr := fmt.Sprintf("SELECT Count(*) FROM Series WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var count int
	err := row.Scan(&count)
	return err == nil && count != 0
}

func(client *SQLClient)ExtendOrUpdateTVTable(series []apiClient.Series) {
	fmt.Println("Check for changes in series and extend or update database accordingly...")
	for _, serie := range series {
		serie.Name = strings.Replace(serie.Name,"'","\\'",-1)
		serie.Tagline = strings.Replace(serie.Tagline,"'","\\'",-1)
		serie.Overview = strings.Replace(serie.Overview,"'","\\'",-1)
		if client.SeriesWithIdExists(serie.ID) {
			client.UpdateSeriesEntry(serie)
		} else {
			fmt.Println("Create")
			client.CreateSeriesEntry(serie)
		}
	}
}

func(client *SQLClient)CreateSeriesEntry(series apiClient.Series) {
	//Eintrag für Film in SQL-DB hinzufügen
	sqlstr := fmt.Sprintf("INSERT INTO Series(id, title, overview, popularity, seasons, episodes, posterPath, voteCount, voteAvg, firstAir, lastAir, tagline) VALUES(%v,'%v','%v',%v,%v,%v,'%v',%v, %v,'%v', '%v', '%v')", series.ID, series.Name, series.Overview, series.Popularity, series.Number_of_seasons, series.Number_of_episodes, series.Poster_Path, series.Vote_Count, series.Vote_Average, series.First_air_date, series.Last_air_date, series.Tagline)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
	for _,genre := range series.Genres {
		if client.GetGenreByID(genre.ID).ID.Valid == false {
			client.CreateGenreEntry(genre)
		}
		client.CreateSeriesGenreEntry(series.ID, genre.ID)
	}
	for _,country := range series.Production_countries {
		if client.GetCountryByID(country.ISO_3166_1).ID.Valid == false {
			client.CreateCountryEntry(country)
		}
		client.CreateSeriesCountryEntry(series.ID, country.ISO_3166_1)
	}
	for _,network := range series.Networks {
		if client.GetNetworkByID(network.ID).ID.Valid == false {
			client.CreateNetworkEntry(network)
		}
		client.CreateSeriesNetworkEntry(series.ID, network)
	}
}

func(client *SQLClient)CreateSeriesGenreEntry(series int, genre int) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create SeriesGenreEntry ")
	sqlstr := fmt.Sprintf("INSERT INTO SeriesGenre(seriesId, genreId) VALUES(%v,'%v')", series, genre)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)GetNetworkByID(id int) Network{
	sqlstr := fmt.Sprintf("SELECT * FROM Networks WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var network Network
	err := row.Scan(&network.ID, &network.NName, &network.Logo, &network.OriginCountry)
	if err != nil && err != sql.ErrNoRows{
		panic(err)
	}
	return network
}

func(client *SQLClient)CreateNetworkEntry(network apiClient.Network) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create Networks-Entry ", network.ID)
	if client.GetCountryByID(network.Origin_country).ID.Valid == false {
		client.CreateCountryEntry(apiClient.Country{network.Origin_country, ""})
	}
	sqlstr := fmt.Sprintf("INSERT INTO Networks(id, nname, logo, originCountry) VALUES('%v','%v', '%v', '%v')", network.ID, network.Name, network.Logo_Path, network.Origin_country)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)GetSeriesByID(id int) Series{
	sqlstr := fmt.Sprintf("SELECT * FROM Series WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var series Series
	err := row.Scan(&series.ID, &series.Title, &series.Overview, &series.Popularity, &series.Seasons, &series.Episodes, &series.Poster, &series.VoteAVG, &series.VoteCount, &series.FirstAit, &series.LastAir, &series.Tagline)
	if err != nil {
		panic(err)
	}
	return series
}

func(client *SQLClient)UpdateSeriesEntry(series apiClient.Series) {
	sqlmovie := client.GetSeriesByID(series.ID)
	sqlmovie.Title.String = strings.Replace(sqlmovie.Title.String,"'","\\'",-1)
	sqlmovie.Tagline.String = strings.Replace(sqlmovie.Tagline.String,"'","\\'",-1)
	sqlmovie.Overview.String = strings.Replace(sqlmovie.Overview.String,"'","\\'",-1)

	different := sqlmovie.Popularity.Float64 != series.Popularity ||
		sqlmovie.VoteCount.Int64 != int64(series.Vote_Count) ||
		sqlmovie.VoteAVG.Float64 != series.Vote_Average ||
		sqlmovie.Seasons.Int64 != int64(series.Number_of_seasons) ||
		sqlmovie.Episodes.Int64 != int64(series.Number_of_episodes) ||
		strings.Compare(sqlmovie.Title.String, series.Name) != 0 ||
		strings.Compare(sqlmovie.Overview.String, series.Overview) != 0 ||
		strings.Compare(sqlmovie.Tagline.String, series.Tagline) != 0 ||
		strings.Compare(sqlmovie.LastAir.String, series.Last_air_date) != 0 ||
		strings.Compare(sqlmovie.Poster.String, series.Poster_Path) != 0

	//Eintrag für Film in SQL-DB hinzufügen
	if different {
		fmt.Println("Update series with id ", series.ID)
		sqlstr := fmt.Sprintf("Update Series SET title='%v', overview='%v', popularity=%v, seasons=%v, episodes=%v, posterPath='%v', voteCount=%v, voteAvg=%v, firstAir='%v', lastAir='%v', tagline='%v' WHERE id=%v", series.Name, series.Overview, series.Popularity, series.Number_of_seasons, series.Number_of_episodes, series.Poster_Path, series.Vote_Count, series.Vote_Average, series.First_air_date, series.Last_air_date, series.Tagline, series.ID)
		_, err := client.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
	}
}

func(client *SQLClient)CreateSeriesCountryEntry(series int, country string) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create SeriesCountryEntry ")
	sqlstr := fmt.Sprintf("INSERT INTO SeriesCountry(seriesId, countryId) VALUES(%v,'%v')", series, country)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)CreateSeriesNetworkEntry(series int, network apiClient.Network) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create SeriesNetworkEntry ")
	if client.GetNetworkByID(network.ID).ID.Valid == false {
		client.CreateNetworkEntry(network)
	}
	sqlstr := fmt.Sprintf("INSERT INTO SeriesNetwork(seriesId, networkId) VALUES(%v,'%v')", series, network.ID)
	_, err := client.Exec(sqlstr)
	if err != nil {
		fmt.Println(network)
		panic(err)
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
		fmt.Println("Create MovieWeekPopularity-Entry")
		client.WriteMovieTrendToSQL(trend, weekNr)
	case nil:

	default:
		panic(err)
	}
}

func(client *SQLClient) WriteMovieTrendToSQL(movie apiClient.Movie, weekNr int) {
	sql := fmt.Sprintf("INSERT INTO MovieWeekPopularity(movieId, weekNr, popularity, voteAVG, voteCount) VALUES ('%v', %v, %v, %v, %v)",movie.ID, weekNr, movie.Popularity, movie.Vote_Average, movie.Vote_Count)
	_, err := client.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	WriteSQLToFile(sql)
}

func(client *SQLClient) WriteTVTrendsToSQL(trends []apiClient.Series, week int) {
	fmt.Println("Write Movie-Trends to database...")
	for _, movie := range trends {

		client.CheckIfTVTrendEntryExist(movie, week)
	}
}

func(client *SQLClient) CheckIfTVTrendEntryExist(trend apiClient.Series, weekNr int) {
	sqlstr := fmt.Sprintf("SELECT seriesid FROM SeriesWeekPopularity WHERE seriesid=%v AND weekNr=%v", trend.ID, weekNr)
	row := client.DB.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		fmt.Println("Create TVWeekPopularity-Entry")
		client.WriteTVTrendToSQL(trend, weekNr)
	case nil:

	default:
		panic(err)
	}
}

func(client *SQLClient) WriteTVTrendToSQL(series apiClient.Series, weekNr int) {
	sql := fmt.Sprintf("INSERT INTO SeriesWeekPopularity(seriesId, weekNr, popularity, voteAVG, voteCount) VALUES ('%v', %v, %v, %v, %v)", series.ID, weekNr, series.Popularity, series.Vote_Average, series.Vote_Count)
	_, err := client.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	WriteSQLToFile(sql)
}