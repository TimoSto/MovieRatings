package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StreamingProvider struct {
	Provider_id int `json:provider_id`
	Provider_name string `json:provider_name`
}

type Provider struct {
	Provider_id int `json:provider_id`
	Provider_name string `json:provider_name`
}

type ProviderResultSetDE struct {
	Rent []Provider `json:rent`
	Buy []Provider `json:buy`
	Flatrate []Provider `json:flatrate`
}

type StreamingProviderResults struct {
	DE ProviderResultSetDE `json:DE`
}

type StreamingProvidersResultSet struct {
	id      int                      `json:id`
	Results StreamingProviderResults `json:results`
}

func(client *APIClient)GetStreamingProvidersForMovieTrends(trends *[]Movie) []Provider{
	var providers []Provider

	for _,movie := range *trends {
		provider := client.GetStreamingProvidersForMovie(movie.ID)
		movie.WatchProviders = provider
		for _,p := range provider.Rent {
			if findProviderInSlice(providers, p.Provider_id) < 0 {
				providers = append(providers, p)
			}
		}
		for _,p := range provider.Buy {
			if findProviderInSlice(providers, p.Provider_id) < 0 {
				providers = append(providers, p)
			}
		}
		for _,p := range provider.Flatrate {
			if findProviderInSlice(providers, p.Provider_id) < 0 {
				providers = append(providers, p)
			}
		}
	}
	return providers
}

func(client *APIClient)GetStreamingProvidersForMovie(id int) ProviderResultSetDE{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/movie/%v/watch/providers?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var resultSet StreamingProvidersResultSet
	err = json.Unmarshal(res, &resultSet)
	if err != nil {
		panic(err)
	}

	return resultSet.Results.DE
}

func findProviderInSlice(arr []Provider, id int) int {

	for i,p := range arr {
		if p.Provider_id == id {
			return i
		}
	}

	return -1
}