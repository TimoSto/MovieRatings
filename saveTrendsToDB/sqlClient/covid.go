package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
)

type Covid struct {
	WeekNr sql.NullInt64 `json:weekNr`
	Cases sql.NullInt64 `json:cases`
	Deaths sql.NullInt64 `json:deaths`
	Recovered sql.NullInt64 `json:recovered`
}

func(client *SQLClient)ExtendOfUpdateCovidTable(stats apiClient.CovidResult) {
	entry, n := client.GetCovidEntry(apiClient.WeekNr)
	if n == -1 {
		client.CreateCovidEntry(stats)
	} else if entry.Cases.Int64 != int64(stats.Cases) ||
		entry.Deaths.Int64 != int64(stats.Deaths) ||
		entry.Recovered.Int64 != int64(stats.Recovered) {
		client.UpdateCovidEntry(stats)
	}
}

func(client *SQLClient)GetCovidEntry(week int) (Covid, int) {
	sqlstr := fmt.Sprintf("SELECT cases, deaths, recovered FROM Covid WHERE weekNr=%v", week)
	row := client.DB.QueryRow(sqlstr)
	var covid Covid
	err := row.Scan(&covid.Cases, &covid.Deaths, &covid.Recovered)
	if err == sql.ErrNoRows {
		return Covid{}, -1
	}
	if err != nil {
		panic(err)
	}
	return covid, 1
}

func(client *SQLClient)CreateCovidEntry(stats apiClient.CovidResult) {
	fmt.Println("Create Covid-Entry")
	sqlstr := fmt.Sprintf("INSERT INTO Covid VALUES(%v, %v, %v, %v)", apiClient.WeekNr, stats.Cases, stats.Deaths, stats.Recovered)
	_,err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)UpdateCovidEntry(stats apiClient.CovidResult) {
	fmt.Println("Update Covid-Entry")
	sqlstr := fmt.Sprintf("UPDATE Covid SET cases=%v, deaths=%v, recovered=%v WHERE weekNr=%v", stats.Cases, stats.Deaths, stats.Recovered, apiClient.WeekNr)
	_,err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}