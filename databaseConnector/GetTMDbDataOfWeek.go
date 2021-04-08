package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Movie struct {
	ID           string          `json:id`
	Title        sql.NullString  `json:title`
	Overview     sql.NullString  `json:overview`
	Popularity   sql.NullFloat64 `json:popularity`
	Revenue      sql.NullFloat64 `json:revenue`
	Poster       sql.NullString  `json:posterPath`
	ReleaseDate  sql.NullString  `json:release_date`
	VoteAVG      sql.NullFloat64 `json:voteAvg`
	VoteCount    sql.NullInt64   `json:VoteCount`
	InProduction sql.NullBool    `json:inProduction`
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

	trends := GetTop100TrendingMovies()

	UpdateFilmTable(trends)

	WriteTrendsToSQL(trends)

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
	//err = GetTop100TrendingMovies()
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
		err = results.Scan(&film.ID, &film.Title, &film.Overview, &film.Popularity, &film.Revenue, &film.Poster, &film.ReleaseDate, &film.VoteAVG, &film.VoteCount, &film.InProduction)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Namme attribute
		movies = append(movies, film)
	}
	return movies
}

func GetTop100TrendingMovies() []TMDbMovie {

	var trends []TMDbMovie
	for i:=1 ; i<=5 ; i++ {
		trendsN := GetTrendingPage(i)
		trends = append(trends, trendsN...)
	}
	return trends
}

func GetTrendingPage(n int) []TMDbMovie{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/trending/movie/week?api_key=b97e33a6b0c4283466ad23df952ebd6a&page=%v", n))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
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
	db, err := sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	for _, movie := range trends {
		if movie.ID == 0 {
			continue
		}
		sqlstr := fmt.Sprintf("SELECT id FROM Movies WHERE id=%v", movie.ID)
		row := db.QueryRow(sqlstr)
		var id int
		switch err := row.Scan(&id); err {
		case sql.ErrNoRows:
			CreateMovieEntry(movie.ID, db)
		case nil:
			fmt.Println("Movie already exists")
		default:
			panic(err)
		}
	}
}

func CreateMovieEntry(id int, db *sql.DB) {
	//db, err := sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
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
	sql := fmt.Sprintf("INSERT INTO Movies(id, title, overview, popularity, releaseDate, posterPath, voteCount, voteAvg, revenue) VALUES(%v,'%v','%v',%v,'%v','%v',%v, %v,'%v')",movie.ID, movie.Title, movie.Overview, movie.Popularity, movie.Release_Date, movie.Poster_Path, movie.Vote_Average, movie.Vote_Count, movie.Revenue)
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}
	//3. SQL-Befehl in SQL-Datei schreiben, um die Daten auch woanders füllen zu können
	WriteSQLToFile(sql)
	//4. Genres ggf ergänzen
	for _,genre := range movie.Genres {
		CheckIfGenreExists(genre, db)
	}
}

func CheckIfGenreExists(genre Genre, db *sql.DB) {
	sqlstr := fmt.Sprintf("SELECT id FROM Genres WHERE id=%v", genre.ID)
	row := db.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		CreateGenreEntry(genre, db)
	case nil:
		fmt.Println("Genre "+genre.Name+" already exists")
	default:
		panic(err)
	}
}

func CreateGenreEntry(genre Genre, db *sql.DB) {
	fmt.Println("Create Genre-Entry for Genre "+genre.Name+"...")
	sql := fmt.Sprintf("INSERT INTO Genres(id, genre) VALUES (%v,'%v')", genre.ID, genre.Name)
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	WriteSQLToFile(sql)
}

func WriteTrendsToSQL(trends []TMDbMovie) {
	db, err := sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	for _, movie := range trends {

		CheckIfTrendEntryExist(movie, db)
	}
}

func CheckIfTrendEntryExist(trend TMDbMovie, db *sql.DB) {
	sqlstr := fmt.Sprintf("SELECT movieid FROM MovieWeekPopularity WHERE movieid=%v AND weekNr=%v", trend.ID, weekNr)
	row := db.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		fmt.Println("Create MovieWeekPopularity-Entry")
		WriteTrendToSQL(trend, db)
	case nil:
		fmt.Println("Trend-Entry already exists")
	default:
		panic(err)
	}
}

func WriteTrendToSQL(movie TMDbMovie, db *sql.DB) {
	sql := fmt.Sprintf("INSERT INTO MovieWeekPopularity(movieId, weekNr, popularity, voteAVG, voteCount) VALUES ('%v', %v, %v, %v, %v)",movie.ID, weekNr, movie.Popularity, movie.Vote_Average, movie.Vote_Count)
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