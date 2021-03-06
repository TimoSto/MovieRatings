package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
	"strings"
)

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

func(client *SQLClient)ExtendOrUpdateTVTable(series []apiClient.Series) {
	fmt.Println("Check for changes in series and extend or update database accordingly...")
	for _, serie := range series {
		serie.Name = strings.Replace(serie.Name,"'","\\'",-1)
		serie.Tagline = strings.Replace(serie.Tagline,"'","\\'",-1)
		serie.Overview = strings.Replace(serie.Overview,"'","\\'",-1)
		if _,n := client.GetSeriesByID(serie.ID); n==-1 {
			client.CreateSeriesEntry(serie)
		} else {
			client.UpdateSeriesEntry(serie)
		}
	}
}

func(client *SQLClient)CreateSeriesEntry(series apiClient.Series) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create SeriesEntry", series.Name)
	sqlstr := fmt.Sprintf("INSERT INTO Series(id, title, overview, popularity, seasons, episodes, posterPath, voteCount, voteAvg, firstAir, lastAir, tagline) VALUES(%v,'%v','%v',%v,%v,%v,'%v',%v, %v,'%v', '%v', '%v')", series.ID, series.Name, series.Overview, series.Popularity, series.Number_of_seasons, series.Number_of_episodes, series.Poster_Path, series.Vote_Count, series.Vote_Average, series.First_air_date, series.Last_air_date, series.Tagline)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
	for _,genre := range series.Genres {

		client.CreateSeriesGenreEntry(series.ID, genre.ID)
	}
	for _,country := range series.Production_countries {

		client.CreateSeriesCountryEntry(series.ID, country.ISO_3166_1)
	}
	for _,network := range series.Networks {

		client.CreateSeriesNetworkEntry(series.ID, network)
	}
	for i,person := range series.Cast {
		if i > 8 {
			break
		}

		client.CreateTVPersonEntry(series.ID, person.ID, "Actor")
	}
	for i,person := range series.Crew {
		if i > 8 {
			break
		}

		client.CreateTVPersonEntry(series.ID, person.ID, person.Job)
	}
	for _,provider := range series.WatchProviders.Buy {

		client.CreateSeriesProviderEntry(series.ID, provider.Provider_id, "buy")
	}
	for _,provider := range series.WatchProviders.Rent {

		client.CreateSeriesProviderEntry(series.ID, provider.Provider_id, "rent")
	}
	for _,provider := range series.WatchProviders.Flatrate {

		client.CreateSeriesProviderEntry(series.ID, provider.Provider_id, "flat")
	}
}

func(client *SQLClient)ExtendTVTrends(series []apiClient.Series) {
	for _, serie := range series {
		if client.CheckIfTVTrendEntryExist(serie.ID, apiClient.WeekNr) == false {
			client.WriteTVTrendToSQL(serie, apiClient.WeekNr)
		}
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

func(client *SQLClient)GetSeriesByID(id int) (Series, int){
	sqlstr := fmt.Sprintf("SELECT * FROM Series WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var series Series
	err := row.Scan(&series.ID, &series.Title, &series.Overview, &series.Popularity, &series.Seasons, &series.Episodes, &series.Poster, &series.VoteAVG, &series.VoteCount, &series.FirstAit, &series.LastAir, &series.Tagline)
	if err == sql.ErrNoRows {
		return Series{}, -1
	}
	if err != nil {
		panic(err)
	}
	return series, 1
}

func(client *SQLClient)UpdateSeriesEntry(series apiClient.Series) {
	sqlmovie,_ := client.GetSeriesByID(series.ID)
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

	sqlstr := fmt.Sprintf("INSERT INTO SeriesNetwork(seriesId, networkId) VALUES(%v,'%v')", series, network.ID)
	_, err := client.Exec(sqlstr)
	if err != nil {
		fmt.Println(network)
		panic(err)
	}
}

func(client *SQLClient)CreateTVPersonEntry(series int, person int, job string) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create SeriesCreditsEntry ")
	job = strings.Replace(job, "'", "\\'", -1)
	sqlstr := fmt.Sprintf("INSERT INTO SeriesCredits(seriesId, personId, job) VALUES(%v,%v,'%v')", series, person, job)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)CreateSeriesProviderEntry(series int, provider int, service string) {
	fmt.Println("Create SeriesProvider Entry", series, provider, series)
	service = strings.Replace(service, "'", "\\'", -1)
	sqlstr := fmt.Sprintf("INSERT INTO SeriesProvider(seriesId, provider, service) VALUES(%v,%v,'%v')", series, provider, service)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient) WriteTVTrendToSQL(series apiClient.Series, weekNr int) {
	fmt.Println("Create SeriesTrend Entry", series.Name, weekNr)
	sql := fmt.Sprintf("INSERT INTO SeriesWeekPopularity(seriesId, weekNr, popularity, voteAVG, voteCount) VALUES ('%v', %v, %v, %v, %v)", series.ID, weekNr, series.Popularity, series.Vote_Average, series.Vote_Count)
	_, err := client.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	//WriteSQLToFile(sql)
}

func(client *SQLClient) CheckIfTVTrendEntryExist(id int, weekNr int) bool{
	sqlstr := fmt.Sprintf("SELECT seriesid FROM SeriesWeekPopularity WHERE seriesid=%v AND weekNr=%v", id, weekNr)
	row := client.DB.QueryRow(sqlstr)
	var found_id int
	err := row.Scan(&found_id)

	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}