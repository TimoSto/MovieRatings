package main

import (
	TMDb "dbconn.com/apiClient"
	MySQL "dbconn.com/sqlClient"
	"fmt"
	"time"
)

var apiClient TMDb.APIClient

var sqlClient MySQL.SQLClient

func main() {
	apiClient = TMDb.APIClient{
		APIKey: "b97e33a6b0c4283466ad23df952ebd6a",
	}

	sqlClient = MySQL.SQLClient{}
	sqlClient.EstablishConnectionToDB()

	tn := time.Now().UTC()
	fmt.Println(tn)
	_, weekNr := tn.ISOWeek()

	movieTrends := apiClient.GetMovieTrends()

	//Die Movie-Objekte müssen anhand der ID aus den Trends nochmal einzeln bestimmt werden, da nicht alle Infos in den Trends drinstehen
	movies := apiClient.GetMovies(movieTrends)

	sqlClient.ExtendOrUpdateMovieTable(movies)

	sqlClient.WriteMovieTrendsToSQL(movieTrends, weekNr)

	tvTrends := apiClient.GetTVTrends()

	series := apiClient.GetSeries(tvTrends)

	sqlClient.ExtendOrUpdateTVTable(series)

	sqlClient.WriteTVTrendsToSQL(tvTrends, weekNr)

	personTrends := apiClient.GetPersonTrends()

	persons := apiClient.GetPersons(personTrends)

	sqlClient.ExtendOrUpdatePersonTable(persons)

	sqlClient.WritePersonTrendsToSQL(personTrends, weekNr)


	defer sqlClient.DB.Close()
}