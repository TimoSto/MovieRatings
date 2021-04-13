package main

import (
	TMDb "dbconn.com/apiClient"
	MySQL "dbconn.com/sqlClient"
)

var apiClient TMDb.APIClient

var sqlClient MySQL.SQLClient

func main() {
	apiClient = TMDb.APIClient{
		APIKey: "b97e33a6b0c4283466ad23df952ebd6a",
	}

	sqlClient = MySQL.SQLClient{}
	sqlClient.EstablishConnectionToDB()

	movieTrends := apiClient.GetMovieTrends()

	//Die Movie-Objekte müssen anhand der ID aus den Trends nochmal einzeln bestimmt werden, da nicht alle Infos in den Trends drinstehen
	movies := apiClient.GetMovies(movieTrends)

	sqlClient.ExtendOrUpdateMovieTable(movies)

	tvTrends := apiClient.GetTVTrends()

	series := apiClient.GetSeries(tvTrends)

	sqlClient.ExtendOrUpdateTVTable(series)

	defer sqlClient.DB.Close()
}