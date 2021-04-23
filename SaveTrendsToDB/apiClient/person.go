package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Person struct {
	ID                   int     `json:id`
	Name                 string  `json:name`
	Birthday             string  `json:birthday`
	Deathday             string  `json:deathday`
	Known_for_department string  `json:known_for_department`
	Gender               int     `json:gender`
	Biography            string  `json:biography`
	Popularity           float64 `json:popularity`
	Profile_path         string  `json:profile_path`
	Job string `json:job`
	Character string `json:character`
}

type PersonJob struct {
	ID int `json:id`
	Job string `json:job`
	Character string `json:character`
}

type PersonRedux struct {
	ID                   int     `json:id`
}

type TrendResultPerson struct {
	Page          int `json:page`
	Results       []PersonRedux `json:results`
}

func (client *APIClient)GetPersonObjects(movies []Movie) []Person {
	var persons []Person
	fmt.Print("Anaysing Person-Informations")
	for _, movie := range movies {
		fmt.Print(".")
		//Cast zur Personen-Liste hinzufügen, falls noch nicht vorhanden
		for _,cast := range movie.Cast {
			if findPersonInSlice(persons, cast.ID) == -1 {
				newPerson := client.GetPersonByID(cast.ID)
				persons = append(persons, newPerson)
			}
		}
		for _, crew := range movie.Crew {
			if findPersonInSlice(persons, crew.ID) == -1 {
				newPerson := client.GetPersonByID(crew.ID)
				persons = append(persons, newPerson)
			}
		}
	}

	return persons
}

func (client *APIClient)GetPersonObjectsTV(series []Series) []Person {
	var persons []Person
	fmt.Print("Anaysing Person-Informations")
	for _, serie := range series {
		fmt.Print(".")
		//Cast zur Personen-Liste hinzufügen, falls noch nicht vorhanden
		for _,cast := range serie.Cast {
			if findPersonInSlice(persons, cast.ID) == -1 {
				newPerson := client.GetPersonByID(cast.ID)
				persons = append(persons, newPerson)
			}
		}
		for _, crew := range serie.Crew {
			if findPersonInSlice(persons, crew.ID) == -1 {
				newPerson := client.GetPersonByID(crew.ID)
				persons = append(persons, newPerson)
			}
		}
	}

	return persons
}

func findPersonInSlice(arr []Person, id int) int {
	for i,p := range arr {
		if p.ID == id {
			return i
		}
	}
	return -1
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

func(client *APIClient)GetPersonTrends() []PersonRedux{
	fmt.Println("Retrieving Person-Trend information from TMDb-Api...")
	var persons []PersonRedux
	for i:=1 ; i <=5 ; i++ {
		persons = append(persons, client.GetPersonTrendPage(i)...)
	}
	return persons
}

func(client *APIClient)GetPersonTrendPage(n int) []PersonRedux{
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

func(client *APIClient)GetPersons(ids []PersonRedux) []Person{
	var persons []Person
	fmt.Print("Analysing Person-Trends")
	for _,id := range ids {
		fmt.Print(".")
		persons = append(persons, client.GetPersonByID(id.ID))
	}

	return persons
}