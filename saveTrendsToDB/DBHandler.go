package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"savetrends.com/apiClient"
	MySQL "savetrends.com/sqlClient"
	"strconv"
	"strings"
)

/*
Dieses Programm kann genutzt werden, um händisch die Covid-Daten einzutragen und um die bestehende Datenbank in eine kompakte SQL-Datei zu komprimieren
Mögliche Parameter:
--convert : Erstellen der SQL-Datei
--covid week=<int> cases=<int> deaths=<int> recovered=<int>
 */

var sqlClient MySQL.SQLClient

func ReadConfig2() MySQL.Config{
	file, err := ioutil.ReadFile("../config.json")

	if err != nil {
		panic(err)
	}
	var conf MySQL.Config
	err = json.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func main() {
	conf := ReadConfig2()
	sqlClient.EstablishConnectionToDB(conf)
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		printHelp()
	}
	switch argsWithoutProg[0] {
	case "--convert":
		convert()
	case "--covid":
		addCovid(argsWithoutProg)
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("USAGE:\\ DBHandler --<COMMAND> <PARAMS>")
	fmt.Println("COMMAND:\n convert          Creates SQL-File to recreate the current database\n covid            Add covid values for week\n	Parameters:\n		week: current weekNr\n		cases: current covid cases\n		deaths: current covid deaths\n		recovered: current recovered covid patients")
}

func convert() {
	var sqlStr string
	err, sqlStr1 := convertMoviesTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr1

	err, sqlStr2 := convertGenresTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr2

	err, sqlStr3 := convertMovieGenreTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr3

	err, sqlStr4 := convertProvidersTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr4

	err, sqlStr5 := convertMovieProviderTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr5

	err, sqlStr6 := convertCountriesTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr6

	err, sqlStr7 := convertMovieCountriesTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr7

	err, sqlStr8 := convertMovieWeekTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr8

	err, sqlStr9 := convertSeriesTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr9

	err, sqlStr10 := convertSeriesGenreTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr10

	err, sqlStr11 := convertSeriesCountriesTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr11

	err, sqlStr12 := convertSeriesWeekTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr12

	err, sqlStr13 := convertSeriesProviderTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr13

	err, sqlStr14 := convertNetworksTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr14

	err, sqlStr15 := convertSeriesNetworkTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr15

	err, sqlStr16 := convertPersons()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr16

	err, sqlStr17 := convertPersonWeekTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr17

	err, sqlStr18 := convertMovieCredits()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr18

	err, sqlStr19 := convertSeriesCredits()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr19

	err, sqlStr20 := convertCovidTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlStr = sqlStr+sqlStr20

	sqlStr = sqlStr + "INSERT into Wochentage Values(0,\"Montag\");\nINSERT into Wochentage Values(1,\"Dienstag\");\nINSERT into Wochentage Values(2,\"Mittwoch\");\nINSERT into Wochentage Values(3,\"Donnerstag\");\nINSERT into Wochentage Values(4,\"Freitag\");\nINSERT into Wochentage Values(5,\"Samstag\");\nINSERT into Wochentage Values(6,\"Sonntag\");"

	ioutil.WriteFile("../database/rebuild.sql", []byte(sqlStr), 0755)
}

func convertCovidTable() (error, string){
	rows, err := sqlClient.DB.Query("select * from covid")
	if err != nil {
		return err, ""
	}
	defer rows.Close()

	var cs []MySQL.Covid
	for rows.Next() {
		var c MySQL.Covid
		err := rows.Scan(&c.WeekNr, &c.Cases, &c.Deaths, &c.Recovered)
		if err != nil {
			return err, ""
		}
		cs = append(cs, c)
	}

	var sqlCommands string
	for _,c := range cs {
		sqlstr := fmt.Sprintf("INSERT INTO covid VALUES(%v,%v, %v, %v)",c.WeekNr.Int64, c.Cases.Int64, c.Deaths.Int64, c.Recovered.Int64) + ";\n"
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func convertMoviesTable() (error, string){
	rows, err := sqlClient.DB.Query("select * from movies")
	if err != nil {
		return err, ""
	}
	defer rows.Close()

	var movies []MySQL.Movie
	for rows.Next() {
		var movie MySQL.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Overview, &movie.Popularity.Float64, &movie.Revenue, &movie.PosterPath, &movie.ReleaseDate, &movie.VoteAVG, &movie.VoteCount, &movie.Runtime, &movie.Tagline)
		if err != nil {
			return err, ""
		}
		movies = append(movies, movie)
	}

	var sqlCommands string
	for _,movie := range movies {
		sqlstr := fmt.Sprintf("INSERT INTO Movies(id, title, overview, popularity, releaseDate, posterPath, voteCount, voteAvg, revenue, runtime, tagline) VALUES(%v,\"%v\",\"%v\",%v,\"%v\",\"%v\",%v, %v,\"%v\", %v, \"%v\")",movie.ID, movie.Title.String, movie.Overview.String, movie.Popularity.Float64, movie.ReleaseDate.String, movie.PosterPath.String, movie.VoteCount.Int64, movie.VoteAVG.Float64, movie.Revenue.String, movie.Runtime.Int64, movie.Tagline.String) + ";\n"
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func convertGenresTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from genres")
	if err != nil {
		return err, ""
	}
	defer rows.Close()

	var genres []MySQL.Genre
	for rows.Next() {
		var genre MySQL.Genre
		err := rows.Scan(&genre.ID, &genre.Genre)
		if err != nil {
			return err, ""
		}
		genres = append(genres, genre)
	}

	var sqlCommands string
	for _,genre := range genres {
		sqlstr := fmt.Sprintf("INSERT INTO Genres(id, genre) VALUES(%v,\"%v\");\n",genre.ID.Int64, genre.Genre.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type MovieGenre struct {
	MovieId sql.NullInt64 `json:movieId`
	GenreId sql.NullInt64 `json:genreId`
}

func convertMovieGenreTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from moviegenre")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []MovieGenre
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(&mg.MovieId, &mg.GenreId)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO MovieGenre(movieId, genreId) VALUES(%v,%v);\n",mg.MovieId.Int64, mg.GenreId.Int64)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func convertProvidersTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from providers")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var ps []MySQL.Provider
	for rows.Next() {
		var p MySQL.Provider
		err := rows.Scan(&p.ID, &p.PName)
		if err != nil {
			return err, ""
		}
		ps = append(ps, p)
	}
	var sqlCommands string
	for _,p := range ps {
		sqlstr := fmt.Sprintf("INSERT INTO Providers(id, pname) VALUES(%v,\"%v\");\n",p.ID.Int64, p.PName.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type MovieProvider struct {
	MovieId sql.NullInt64 `json:movieId`
	Provider sql.NullInt64 `json:provider`
	Service sql.NullString `json:service`
}

func convertMovieProviderTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from movieprovider")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mps []MovieProvider
	for rows.Next() {
		var mp MovieProvider
		err := rows.Scan(&mp.MovieId, &mp.Provider, &mp.Service)
		if err != nil {
			return err, ""
		}
		mps = append(mps, mp)
	}
	var sqlCommands string
	for _,mp := range mps {
		sqlstr := fmt.Sprintf("INSERT INTO MovieProvider(movieId, provider, service) VALUES(%v,%v,\"%v\");\n", mp.MovieId.Int64, mp.Provider.Int64, mp.Service.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func convertCountriesTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from countries")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []MySQL.Country
	for rows.Next() {
		var mg MySQL.Country
		err := rows.Scan(&mg.ID, &mg.CName)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO Countries(id, cname) VALUES(\"%v\",\"%v\");\n",mg.ID.String, mg.CName.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type MovieCountry struct {
	MovieID sql.NullInt64 `json:movieId`
	CountryId sql.NullString `json:countryId`
}

func convertMovieCountriesTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from moviecountry")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []MovieCountry
	for rows.Next() {
		var mg MovieCountry
		err := rows.Scan(&mg.MovieID, &mg.CountryId)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO MovieCountry(movieId, countryId) VALUES(%v,\"%v\");\n",mg.MovieID.Int64, mg.CountryId.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type MWP struct {
	MovieId sql.NullInt64 `json:movieId`
	WeekNr sql.NullInt64 `json:weekNr`
	Popularity sql.NullFloat64 `json:popularity`
	VoteAvg sql.NullFloat64 `json:voteAvg`
	VoteCount sql.NullInt64 `json:voteCount`
}

func convertMovieWeekTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from movieweekPopularity")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []MWP
	for rows.Next() {
		var mg MWP
		err := rows.Scan(&mg.MovieId, &mg.WeekNr, &mg.Popularity, &mg.VoteAvg, &mg.VoteCount)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO MovieWeekPopularity(movieId, weekNr, popularity, voteAVG, voteCount) VALUES (%v, \"%v\", %v, %v, %v);\n",mg.MovieId.Int64, mg.WeekNr.Int64, mg.Popularity.Float64, mg.VoteAvg.Float64, mg.VoteCount.Int64)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func convertSeriesTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from series")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []MySQL.Series
	for rows.Next() {
		var series MySQL.Series
		err := rows.Scan(&series.ID, &series.Title, &series.Overview, &series.Popularity, &series.Seasons, &series.Episodes, &series.Poster, &series.VoteAVG, &series.VoteCount, &series.FirstAit, &series.LastAir, &series.Tagline)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, series)
	}
	var sqlCommands string
	for _,series := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO Series(id, title, overview, popularity, seasons, episodes, posterPath, voteCount, voteAvg, firstAir, lastAir, tagline) VALUES(%v,\"%v\",\"%v\",%v,%v,%v,\"%v\",%v, %v,\"%v\", \"%v\", \"%v\");\n", series.ID, series.Title.String, series.Overview.String, series.Popularity.Float64, series.Seasons.Int64, series.Episodes.Int64, series.Poster.String, series.VoteCount.Int64, series.VoteAVG.Float64, series.FirstAit.String, series.LastAir.String, series.Tagline.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type SeriesGenre struct {
	SeriesId sql.NullInt64 `json:seriesId`
	GenreId sql.NullInt64 `json:genreId`
}

func convertSeriesGenreTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from seriesgenre")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []SeriesGenre
	for rows.Next() {
		var mg SeriesGenre
		err := rows.Scan(&mg.SeriesId, &mg.GenreId)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO SeriesGenre(movieId, genreId) VALUES(%v,%v);\n",mg.SeriesId.Int64, mg.GenreId.Int64)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type SeriesCountry struct {
	SeriesID sql.NullInt64 `json:seriesId`
	CountryId sql.NullString `json:countryId`
}

func convertSeriesCountriesTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from seriescountry")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []SeriesCountry
	for rows.Next() {
		var mg SeriesCountry
		err := rows.Scan(&mg.SeriesID, &mg.CountryId)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO SeriesCountry(movieId, countryId) VALUES(%v,\"%v\");\n",mg.SeriesID.Int64, mg.CountryId.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type SWP struct {
	SeriesId sql.NullInt64 `json:seriesId`
	WeekNr sql.NullInt64 `json:weekNr`
	Popularity sql.NullFloat64 `json:popularity`
	VoteAvg sql.NullFloat64 `json:voteAvg`
	VoteCount sql.NullInt64 `json:voteCount`
}

func convertSeriesWeekTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from seriesweekPopularity")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []SWP
	for rows.Next() {
		var mg SWP
		err := rows.Scan(&mg.SeriesId, &mg.WeekNr, &mg.Popularity, &mg.VoteAvg, &mg.VoteCount)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO SeriesWeekPopularity(movieId, weekNr, popularity, voteAVG, voteCount) VALUES (%v, \"%v\", %v, %v, %v);\n",mg.SeriesId.Int64, mg.WeekNr.Int64, mg.Popularity.Float64, mg.VoteAvg.Float64, mg.VoteCount.Int64)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type SeriesProvider struct {
	SeriesId sql.NullInt64 `json:seriesId`
	Provider sql.NullInt64 `json:provider`
	Service sql.NullString `json:service`
}

func convertSeriesProviderTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from seriesprovider")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mps []SeriesProvider
	for rows.Next() {
		var mp SeriesProvider
		err := rows.Scan(&mp.SeriesId, &mp.Provider, &mp.Service)
		if err != nil {
			return err, ""
		}
		mps = append(mps, mp)
	}
	var sqlCommands string
	for _,mp := range mps {
		sqlstr := fmt.Sprintf("INSERT INTO SeriesProvider(movieId, provider, service) VALUES(%v,%v,\"%v\");\n", mp.SeriesId.Int64, mp.Provider.Int64, mp.Service.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func convertNetworksTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from networks")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var ps []MySQL.Network
	for rows.Next() {
		var p MySQL.Network
		err := rows.Scan(&p.ID, &p.NName, &p.Logo, &p.OriginCountry)
		if err != nil {
			return err, ""
		}
		ps = append(ps, p)
	}
	var sqlCommands string
	for _,network := range ps {
		sqlstr := fmt.Sprintf("INSERT INTO Networks(id, nname, logo, originCountry) VALUES(\"%v\",\"%v\", \"%v\", \"%v\");\n", network.ID.Int64, network.NName.String, network.Logo.String, network.OriginCountry.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type SeriesNetwork struct {
	SeriesId sql.NullInt64 `json:seriesId`
	NetworkId sql.NullInt64 `json:networkId`
}

func convertSeriesNetworkTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from seriesnetwork")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var ps []SeriesNetwork
	for rows.Next() {
		var p SeriesNetwork
		err := rows.Scan(&p.SeriesId, &p.NetworkId)
		if err != nil {
			return err, ""
		}
		ps = append(ps, p)
	}
	var sqlCommands string
	for _,sn := range ps {
		sqlstr := fmt.Sprintf("INSERT INTO SeriesNetwork(seriesId, networkId) VALUES(%v,\"%v\");\n", sn.SeriesId.Int64, sn.NetworkId.Int64)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func convertPersons() (error, string) {
	rows, err := sqlClient.DB.Query("select * from personen")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var ps []MySQL.Person
	for rows.Next() {
		var person MySQL.Person
		err := rows.Scan(&person.ID, &person.Name, &person.Birthday, &person.Deathday, &person.Popularity, &person.ProfilePath, &person.Gender, &person.Profession)
		if err != nil {
			return err, ""
		}
		ps = append(ps, person)
	}
	var sqlCommands string
	for _,person := range ps {
		sqlstr := fmt.Sprintf("INSERT INTO Personen(id, name, birthday, deathday, popularity, profilePath, gender, profession) VALUES(%v,\"%v\",\"%v\",\"%v\",%v,\"%v\",%v,\"%v\");\n", person.ID.Int64, person.Name.String, person.Birthday.String, person.Deathday.String, person.Popularity.Float64, person.ProfilePath.String, person.Gender.Int64, person.Profession.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type PWP struct {
	PersonId sql.NullInt64 `json:personId`
	WeekNr sql.NullInt64 `json:weekNr`
	Popularity sql.NullFloat64 `json:popularity`
}

func convertPersonWeekTable() (error, string) {
	rows, err := sqlClient.DB.Query("select * from personweek")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs []PWP
	for rows.Next() {
		var mg PWP
		err := rows.Scan(&mg.PersonId, &mg.WeekNr, &mg.Popularity)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO PersonWeek(personId, weekNr, popularity) VALUES (%v, %v, %v);\n", mg.PersonId.Int64, mg.WeekNr.Int64, mg.Popularity.Float64)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type MovieCredits struct {
	MovieId sql.NullInt64 `json:movieId`
	PersonId sql.NullInt64 `json:personId`
	Job sql.NullString `json:job`
}

func convertMovieCredits() (error, string) {
	rows, err := sqlClient.DB.Query("select * from moviecredits")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs [] MovieCredits
	for rows.Next() {
		var mg MovieCredits
		err := rows.Scan(&mg.MovieId, &mg.PersonId, &mg.Job)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO MovieCredits(movieId, personId, job) VALUES(%v,%v,\"%v\");\n",mg.MovieId.Int64, mg.PersonId.Int64, mg.Job.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

type SeriesCredits struct {
	SeriesId sql.NullInt64 `json:seriesId`
	PersonId sql.NullInt64 `json:personId`
	Job sql.NullString `json:job`
}

func convertSeriesCredits() (error, string) {
	rows, err := sqlClient.DB.Query("select * from moviecredits")
	if err != nil {
		return err, ""
	}
	defer rows.Close()
	var mgs [] SeriesCredits
	for rows.Next() {
		var mg SeriesCredits
		err := rows.Scan(&mg.SeriesId, &mg.PersonId, &mg.Job)
		if err != nil {
			return err, ""
		}
		mgs = append(mgs, mg)
	}
	var sqlCommands string
	for _,mg := range mgs {
		sqlstr := fmt.Sprintf("INSERT INTO SeriesCredits(seriesId, personId, job) VALUES(%v,%v,\"%v\");\n", mg.SeriesId.Int64, mg.PersonId.Int64, mg.Job.String)
		sqlCommands = sqlCommands + sqlstr
	}
	return nil, sqlCommands
}

func addCovid(args []string) {
	data := apiClient.CovidResult{
		Cases:     0,
		Deaths:    0,
		Recovered: 0,
	}
	args = args[1:]
	for _,d := range args {
		parts := strings.Split(d, "=")
		value,_ := strconv.ParseInt(parts[1],10,64)
		switch parts[0] {
		case "week":
			apiClient.WeekNr = int(value)
		case "cases":
			data.Cases = int(value)
		case "deaths":
			data.Deaths = int(value)
		case "recovered":
			data.Recovered = int(value)
		default:
			printHelp()
			return
		}
	}
	sqlClient.ExtendOfUpdateCovidTable(data)
}