package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
	"strings"
)

type Person struct {
	ID sql.NullInt64 `json:id`
	Name sql.NullString `json:name`
	Birthday sql.NullString `json:birthday`
	Deathday sql.NullString `json:deathday`
	Popularity sql.NullFloat64 `json:popularity`
	ProfilePath sql.NullString `json:profilePath`
	Gender sql.NullInt64 `json:gender`
	Profession sql.NullString `json:profession`
}

func(client *SQLClient)GetPersonByID(id int) (Person, int){
	sqlstr := fmt.Sprintf("SELECT * FROM Personen WHERE id='%v'", id)
	row := client.DB.QueryRow(sqlstr)
	var person Person
	err := row.Scan(&person.ID, &person.Name, &person.Birthday, &person.Deathday, &person.Popularity, &person.ProfilePath, &person.Gender, &person.Profession)
	if err == sql.ErrNoRows {
		return Person{}, -1
	} else if err != nil {
		panic(err)
	}
	return person, 1
}

func(client *SQLClient)ExtendOrUpdatePersonTable(persons []apiClient.Person) {

	for _, person := range persons {
		sqlperson, exists := client.GetPersonByID(person.ID)
		if exists == -1 {
			client.CreatePersonEntry(person)
		} else {
			client.UpdatePersonEntry(person, sqlperson)
		}
	}
}

func(client *SQLClient)CreatePersonEntry(person apiClient.Person) {
	fmt.Println("Create Person "+person.Name)
	person.Name = strings.Replace(person.Name, "'","\\'", -1)
	sqlstr := fmt.Sprintf("INSERT INTO Personen(id, name, birthday, deathday, popularity, profilePath, gender, profession) VALUES(%v,'%v','%v','%v',%v,'%v',%v,'%v')", person.ID, person.Name, person.Birthday, person.Deathday, person.Popularity, person.Profile_path, person.Gender, person.Known_for_department)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)UpdatePersonEntry(person apiClient.Person, sqlperson Person) {

	different := strings.Compare(person.Birthday, sqlperson.Birthday.String) != 0 ||
		strings.Compare(person.Deathday, sqlperson.Deathday.String) != 0 ||
		person.Popularity != sqlperson.Popularity.Float64

	if different {
		fmt.Println("Update Person "+person.Name)
		sqlstr := fmt.Sprintf("UPDATE Personen set birthday='%v', deathday='%v', popularity=%v where id=%v",person.Birthday, person.Deathday, person.Popularity, person.ID)
		_,err := client.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
	}
}