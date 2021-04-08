package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Movie struct {
	ID           int         `json:id`
	Title        sql.NullString  `json:title`
	Overview     sql.NullString  `json:overview`
	Popularity   sql.NullFloat64 `json:popularity`
	Revenue      sql.NullString `json:revenue`
	Poster       sql.NullString  `json:posterPath`
	ReleaseDate  sql.NullString  `json:release_date`
	VoteAVG      sql.NullFloat64 `json:voteAvg`
	VoteCount    sql.NullInt64   `json:VoteCount`
}

type TMDbMovie struct {
	ID           int     `json:id`
	Title        string  `json:title`
	Overview     string  `json:overview`
	Popularity   float64 `json:popularity`
	Revenue      float64 `json:revenue`
	Poster_Path       string  `json:poster_path`
	Release_Date string  `json:release_date`
	Vote_Average float64 `json:vote_average`
	Vote_Count   float64 `json:vote_count`
	Genres     []Genre   `json:genres`
}

type TMDbSeries struct {
	ID           int     `json:id`
	Title        string  `json:title`
	Overview     string  `json:overview`
	Popularity   float64 `json:popularity`
	Seasons      []struct{} `json:seasons`
	Poster_Path       string  `json:poster_path`
	First_air_date string  `json:first_air_date`
	Vote_Average float64 `json:vote_average`
	Vote_Count   float64 `json:vote_count`
	Genres     []Genre   `json:genres`
}

type Genre struct {
	ID    int    `json:id`
	Name string `json:name`
}

type MovieWeek struct {
	ID         int
	Week       int
	Popularity float64
	Revenue    float64
	VoteAVG    float64
	VoteCount  float64
}

type TrendRes struct {
	Page          int `json:"page"`
	Results       []TMDbMovie `json:results`
	Total_Pages   int
	Total_Results int
}

var db *sql.DB

var weekNr int

func main() {
	tn := time.Now().UTC()
	fmt.Println(tn)
	_, weekNr = tn.ISOWeek()

	var err error
	db, err = sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	defer db.Close()
	if err != nil {
		panic(err)
	}

	//movie := GetSQLMovie(11)
	//
	//fmt.Println(movie.Title)

	trends := GetTop100Trending("movie")

	UpdateFilmTable(trends)

	WriteMovieTrendsToSQL(trends)

	trends = GetTop100Trending("tv")

	UpdateSeriesTable(trends)

	WriteTVTrendsToSQL(trends)

	//fmt.Println("Reading Config...")
	//
	//file,err := ioutil.ReadFile("./config.json")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(file, &config)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Requesting relevant data from TMDb via TMDb-API...")
	//
	//db, err = sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = GetTop100Trending()
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer db.Close()
}

func GetMovieByID(id string) TMDbMovie {
	resp, err := http.Get("https://api.themoviedb.org/3/movie/527774?api_key=b97e33a6b0c4283466ad23df952ebd6a")
	res, err := ioutil.ReadAll(resp.Body)
	var movie TMDbMovie
	err = json.Unmarshal(res,&movie)
	if err != nil {
		panic(err)
	}
	fmt.Println(movie.Title)
	return movie
}

func GetAllMovies() []Movie {
	results, err := db.Query("SELECT * FROM film")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Println(results.Columns())

	var movies []Movie

	for results.Next() {
		var film Movie
		// for each row, scan the result into our tag composite object
		err = results.Scan(&film.ID, &film.Title, &film.Overview, &film.Popularity, &film.Revenue, &film.Poster, &film.ReleaseDate, &film.VoteAVG, &film.VoteCount)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Namme attribute
		movies = append(movies, film)
	}
	return movies
}

func GetSQLMovie(mid int) Movie{

	sqlstr := fmt.Sprintf("SELECT * FROM Movies WHERE id=%v", mid)
	row := db.QueryRow(sqlstr)
	var id int
	var title string
	var overview string
	var popularity float64
	var revenue string
	var posterPath string
	var releaseDate string
	var voteAvg float64
	var voteCount int
	err := row.Scan(&id, &title, &overview, &popularity,&revenue,&posterPath,&releaseDate,&voteAvg,&voteCount)
	if err != nil {
		panic(err)
	}
	return Movie{
		id,
		sql.NullString{title, true},
		sql.NullString{overview, true},
		sql.NullFloat64{popularity, true},
		sql.NullString{revenue, true},
		sql.NullString{posterPath, true},
		sql.NullString{releaseDate, true},
		sql.NullFloat64{voteAvg, true},
		sql.NullInt64{int64(voteCount), true},
	}
}

func GetTop100Trending(typ string) []TMDbMovie {

	var trends []TMDbMovie
	for i:=1 ; i<=5 ; i++ {
		trendsN := GetTrendingPage(typ, i)
		trends = append(trends, trendsN...)
	}
	return trends
}

func GetTrendingPage(typ string, n int) []TMDbMovie{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/trending/%v/week?api_key=b97e33a6b0c4283466ad23df952ebd6a&page=%v", typ, n))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var trendingResultSet TrendRes
	err = json.Unmarshal(res,&trendingResultSet)
	if err != nil {
		panic(err)
	}

	return trendingResultSet.Results
}

func UpdateFilmTable(trends []TMDbMovie) {
	//Neue Filme hinzufügen, damit der F-Key in der Week-Trend-Tabelle existiert
	for _, movie := range trends {
		if movie.ID == 0 {
			continue
		}
		sqlstr := fmt.Sprintf("SELECT id FROM Movies WHERE id=%v", movie.ID)
		row := db.QueryRow(sqlstr)
		var id int
		switch err := row.Scan(&id); err {
		case sql.ErrNoRows:
			CreateMovieEntry(movie.ID)
		case nil:
			fmt.Println("Movie already exists")
			UpdateMovieEntry(movie.ID)
		default:
			panic(err)
		}
	}
}

func CreateMovieEntry(id int) {
	//db, err := sqlstr.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	//defer db.Close()
	//1. Daten zum Film über HTTP aus TMDb-API ermitteln
	fmt.Println("Create Movie Entry for ID ",id)
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/%v?api_key=b97e33a6b0c4283466ad23df952ebd6a", id))
	if err != nil {
		panic(err)
	}
	var movie TMDbMovie
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body,&movie)
	if err != nil {
		panic(err)
	}
	movie.Title = strings.Replace(movie.Title,"'","\\'",-1)
	movie.Overview = strings.Replace(movie.Overview,"'","\\'",-1)
	//2. Eintrag für Film in SQL-DB hinzufügen
	sqlstr := fmt.Sprintf("INSERT INTO Movies(id, title, overview, popularity, releaseDate, posterPath, voteCount, voteAvg, revenue) VALUES(%v,'%v','%v',%v,'%v','%v',%v, %v,'%v')",movie.ID, movie.Title, movie.Overview, movie.Popularity, movie.Release_Date, movie.Poster_Path, movie.Vote_Count, movie.Vote_Average, movie.Revenue)
	_, err = db.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
	//3. SQL-Befehl in SQL-Datei schreiben, um die Daten auch woanders füllen zu können
	WriteSQLToFile(sqlstr)
	//4. Genres ggf ergänzen
	for _,genre := range movie.Genres {
		CheckIfGenreExists(genre)
		CreateMovieGenreEntry(movie.ID, genre.ID)
	}
}

func UpdateMovieEntry(id int) {
	//db, err := sqlstr.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	//defer db.Close()
	//1. Daten zum Film über HTTP aus TMDb-API ermitteln
	fmt.Println("Ceck for updates for Movie Entry with ID ",id)
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/%v?api_key=b97e33a6b0c4283466ad23df952ebd6a", id))
	if err != nil {
		panic(err)
	}
	var movie TMDbMovie
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body,&movie)
	if err != nil {
		panic(err)
	}
	movie.Title = strings.Replace(movie.Title,"'","\\'",-1)
	movie.Overview = strings.Replace(movie.Overview,"'","\\'",-1)
	//2. Daten zum Movie aus der DB ziehen
	sqlmovie := GetSQLMovie(id)
	//2. Eintrag für Film in SQL-DB aktualisieren, wenn sich etwas geändert hat
	different := sqlmovie.VoteAVG.Float64 != movie.Vote_Average ||
		sqlmovie.VoteCount.Int64 != int64(movie.Vote_Count) ||
		sqlmovie.Popularity.Float64 != movie.Popularity ||
		strings.Compare(fmt.Sprintf("%f",convertExponentialStringToFloat(sqlmovie.Revenue.String)),fmt.Sprintf("%f",movie.Revenue)) != 0
	if different {
		fmt.Println("Update Movie Entry for ID ", movie.ID)
		sqlstr := fmt.Sprintf("UPDATE Movies SET voteAvg=%v, voteCount=%v, popularity=%v, revenue='%v' WHERE id=%v",movie.Vote_Average, movie.Vote_Count, movie.Popularity, movie.Revenue, movie.ID)
		_, err = db.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
		WriteSQLToFile(sqlstr)
	}
}

func UpdateSeriesTable(trends []TMDbMovie) {
	for _, series := range trends {
		if series.ID == 0 {
			continue
		}
		sqlstr := fmt.Sprintf("SELECT id FROM Series WHERE id=%v", series.ID)
		row := db.QueryRow(sqlstr)
		var id int
		switch err := row.Scan(&id); err {
		case sql.ErrNoRows:
			CreateTVEntry(series.ID)
		case nil:
			fmt.Println("Series already exists")
		default:
			panic(err)
		}
	}
}

func CreateTVEntry(id int) {
	//db, err := sqlstr.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	//defer db.Close()
	//1. Daten zum Film über HTTP aus TMDb-API ermitteln
	fmt.Println("Create Series Entry for ID ",id)
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/tv/%v?api_key=b97e33a6b0c4283466ad23df952ebd6a", id))
	if err != nil {
		panic(err)
	}
	var series TMDbSeries
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body,&series)
	if err != nil {
		panic(err)
	}
	series.Title = strings.Replace(series.Title,"'","\\'",-1)
	series.Overview = strings.Replace(series.Overview,"'","\\'",-1)
	//2. Eintrag für Film in SQL-DB hinzufügen
	sqlstr := fmt.Sprintf("INSERT INTO Series(id, title, overview, popularity, seasons, releaseDate, posterPath, voteCount, voteAvg) VALUES(%v,'%v','%v',%v,%v,'%v','%v', %v, %v)",series.ID, series.Title, series.Overview, series.Popularity, len(series.Seasons), series.First_air_date, series.Poster_Path, series.Vote_Average, series.Vote_Count)
	_, err = db.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
	//3. SQL-Befehl in SQL-Datei schreiben, um die Daten auch woanders füllen zu können
	WriteSQLToFile(sqlstr)
	//4. Genres ggf ergänzen
	for _,genre := range series.Genres {
		CheckIfGenreExists(genre)
		CreateSeriesGenreEntry(series.ID, genre.ID)
	}
}

func CheckIfGenreExists(genre Genre) {
	sqlstr := fmt.Sprintf("SELECT id FROM Genres WHERE id=%v", genre.ID)
	row := db.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		CreateGenreEntry(genre)
	case nil:
		fmt.Println("Genre "+genre.Name+" already exists")
	default:
		panic(err)
	}
}

func CreateGenreEntry(genre Genre) {
	fmt.Println("Create Genre-Entry for Genre "+genre.Name+"...")
	sqlstr := fmt.Sprintf("INSERT INTO Genres(id, genre) VALUES (%v,'%v')", genre.ID, genre.Name)
	_, err := db.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	WriteSQLToFile(sqlstr)
}

func CreateMovieGenreEntry(movieID int, genreID int) {
	sqlstr := fmt.Sprintf("SELECT movieid FROM MovieGenre WHERE movieid=%v AND genreID=%v", movieID, genreID)
	row := db.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		sqlstr := fmt.Sprintf("INSERT INTO MovieGenre(movieId, genreId) VALUES (%v,%v)",movieID, genreID)
		_, err := db.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
		//fmt.Println(res)
		WriteSQLToFile(sqlstr)
	case nil:

	default:
		panic(err)
	}
}

func CreateSeriesGenreEntry(seriesID int, genreID int) {
	sqlstr := fmt.Sprintf("SELECT seriesid FROM SeriesGenre WHERE seriesid=%v AND genreID=%v", seriesID, genreID)
	row := db.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		sqlstr := fmt.Sprintf("INSERT INTO SeriesGenre(seriesId, genreId) VALUES (%v,%v)", seriesID, genreID)
		_, err := db.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
		//fmt.Println(res)
		WriteSQLToFile(sqlstr)
	case nil:

	default:
		panic(err)
	}
}

func WriteMovieTrendsToSQL(trends []TMDbMovie) {
	db, err := sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	for _, movie := range trends {

		CheckIfMovieTrendEntryExist(movie)
	}
}

func CheckIfMovieTrendEntryExist(trend TMDbMovie) {
	sqlstr := fmt.Sprintf("SELECT movieid FROM MovieWeekPopularity WHERE movieid=%v AND weekNr=%v", trend.ID, weekNr)
	row := db.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		fmt.Println("Create MovieWeekPopularity-Entry")
		WriteMovieTrendToSQL(trend)
	case nil:
		fmt.Println("Trend-Entry already exists")
	default:
		panic(err)
	}
}

func WriteMovieTrendToSQL(movie TMDbMovie) {
	sql := fmt.Sprintf("INSERT INTO MovieWeekPopularity(movieId, weekNr, popularity, voteAVG, voteCount) VALUES ('%v', %v, %v, %v, %v)",movie.ID, weekNr, movie.Popularity, movie.Vote_Average, movie.Vote_Count)
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	WriteSQLToFile(sql)
}

func WriteTVTrendsToSQL(trends []TMDbMovie) {
	db, err := sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	for _, series := range trends {

		CheckIfTVTrendEntryExist(series)
	}
}

func CheckIfTVTrendEntryExist(trend TMDbMovie) {
	sqlstr := fmt.Sprintf("SELECT seriesid FROM SeriesWeekPopularity WHERE seriesid=%v AND weekNr=%v", trend.ID, weekNr)
	row := db.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		fmt.Println("Create TVWeekPopularity-Entry")
		WriteTVTrendToSQL(trend)
	case nil:
		fmt.Println("Trend-Entry already exists", trend.ID)
	default:
		panic(err)
	}
}

func WriteTVTrendToSQL(series TMDbMovie) {
	sql := fmt.Sprintf("INSERT INTO SeriesWeekPopularity(seriesId, weekNr, popularity, voteAVG, voteCount) VALUES ('%v', %v, %v, %v, %v)", series.ID, weekNr, series.Popularity, series.Vote_Average, series.Vote_Count)
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	WriteSQLToFile(sql)
}

func WriteSQLToFile(sql string){
	f, err := os.OpenFile("./database/FILLDB.sql",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	//sql = strings.Replace(sql,"'","\\'",-1)
	if _, err := f.WriteString(sql+";\n"); err != nil {
		log.Println(err)
	}
}

func convertExponentialStringToFloat(str string) float64{
	parts := strings.Split(str,"e")
	base, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		panic(err)
	}
	if len(parts) == 1 {
		return base
	}
	exp, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		panic(err)
	}
	return base * math.Pow(10, exp)
}