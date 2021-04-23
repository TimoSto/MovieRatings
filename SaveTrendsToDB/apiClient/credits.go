package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CreditsForMovieOrTV struct {
	ID int `json:id`
	Cast []PersonJob `json:cast`
	Crew []PersonJob `json:crew`
}

func(client *APIClient)GetCreditsForMovie(id int) CreditsForMovieOrTV {
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/%v/credits?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var credits CreditsForMovieOrTV
	err = json.Unmarshal(res, &credits)
	if err != nil {
		panic(err)
	}
	return credits
}