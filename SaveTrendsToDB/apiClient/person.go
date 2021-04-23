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

func (client *APIClient)GetPersonObjects(movies []Movie) []Person {
	var persons []Person

	for _, movie := range movies {
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