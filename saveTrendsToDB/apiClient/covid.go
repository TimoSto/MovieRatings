package apiClient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CovidResult struct {
	Cases int `json:cases`
	Deaths int `json:deaths`
	Recovered int `json:recovered`
}

type CovidResultSet struct {
	Current_totals CovidResult `json:current_totals`
}

func(client *APIClient) GetCovidStats() CovidResult{
	resp, err := http.Get("https://covid19-germany.appspot.com/now")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var res CovidResultSet
	err = json.Unmarshal(data, &res)
	if err != nil {
		panic(err)
	}
	return res.Current_totals
}