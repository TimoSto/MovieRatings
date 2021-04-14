package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func(client *APIClient)GetPersonTrends() []Person{
	var persons []Person
	for i:=1 ; i <=5 ; i++ {
		persons = append(persons, client.GetPersonTrendPage(i)...)
	}
	return persons
}

func(client *APIClient)GetPersonTrendPage(n int) []Person{
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/person/popular?api_key=%v&page=%v", client.APIKey, n))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var trendingResultSet TrendResultPerson
	err = json.Unmarshal(res,&trendingResultSet)
	if err != nil {
		panic(err)
	}

	return trendingResultSet.Results
}

func(client *APIClient)GetPersons(trends []Person) []Person {
	var persons []Person
	for _, trend := range trends {
		persons = append(persons, client.GetPersonByID(trend.ID))
	}
	return persons
}

func(client *APIClient)GetPersonByID(id int) Person {
	resp, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/person/%v?api_key=%v", id, client.APIKey))
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var person Person
	err = json.Unmarshal(res, &person)
	if err != nil {
		panic(err)
	}

	return person
}