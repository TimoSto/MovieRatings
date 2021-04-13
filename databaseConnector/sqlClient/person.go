package sqlClient

import (
	"database/sql"
	"dbconn.com/apiClient"
	"fmt"
	"strings"
)

func(client *SQLClient)GetPersonByID(id int) Person{
	sqlstr := fmt.Sprintf("SELECT * FROM Personen WHERE id='%v'", id)
	row := client.DB.QueryRow(sqlstr)
	var person Person
	err := row.Scan(&person.ID, &person.Name, &person.Birthday, &person.Deathday, &person.Popularity, &person.ProfilePath, &person.Gender, &person.Profession)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return person
}

func(client *SQLClient)CreatePersonEntry(person apiClient.Person) {
	//Eintrag für Film in SQL-DB hinzufügen
	person.Name = strings.Replace(person.Name, "'","\\'", -1)
	fmt.Println("Create Person-Entry "+ person.Name)
	sqlstr := fmt.Sprintf("INSERT INTO Personen(id, name, birthday, deathday, popularity, profilePath, gender, profession) VALUES(%v,'%v','%v','%v',%v,'%v',%v,'%v')", person.ID, person.Name, person.Birthday, person.Deathday, person.Popularity, person.Profile_path, person.Gender, person.Known_for_department)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)ExtendOrUpdatePersonTable(persons []apiClient.Person) {
	fmt.Println("Check for changes in persons and extend or update database accordingly...")
	for _, person := range persons {
		person.Name = strings.Replace(person.Name,"'","\\'",-1)
		person.Job = strings.Replace(person.Job,"'","\\'",-1)
		if client.GetPersonByID(person.ID).ID.Valid == true {
			client.UpdatePersonEntry(person)
		} else {
			client.CreatePersonEntry(person)
		}
	}
}

func(client *SQLClient)UpdatePersonEntry(person apiClient.Person) {

	sqlperson := client.GetPersonByID(person.ID)
	different := sqlperson.Popularity.Float64 != person.Popularity ||
		strings.Compare(sqlperson.Deathday.String, person.Deathday) != 0
	if different {
		fmt.Println("Update Person", person.ID)
		sqlstr := fmt.Sprintf("UPDATE TABLE Personen set deathday='%v', popularity=%v where id=%v",person.Deathday, person.Popularity, person.ID)
		_,err := client.Exec(sqlstr)
		if err != nil {
			panic(err)
		}
	}
}

func(client *SQLClient) WritePersonTrendsToSQL(trends []apiClient.Person, week int) {
	fmt.Println("Write Person-Trends to database...")
	for _, person := range trends {

		client.CheckIfPersonTrendEntryExist(person, week)
	}
}

func(client *SQLClient) CheckIfPersonTrendEntryExist(trend apiClient.Person, weekNr int) {
	sqlstr := fmt.Sprintf("SELECT personId FROM PersonWeek WHERE personId=%v AND weekNr=%v", trend.ID, weekNr)
	row := client.DB.QueryRow(sqlstr)
	var found_id int
	switch err := row.Scan(&found_id); err {
	case sql.ErrNoRows:
		fmt.Println("Create PersonWeekPopularity-Entry")
		client.WritePersonTrendToSQL(trend, weekNr)
	case nil:

	default:
		panic(err)
	}
}

func(client *SQLClient) WritePersonTrendToSQL(person apiClient.Person, weekNr int) {
	sql := fmt.Sprintf("INSERT INTO PersonWeek(personId, weekNr) VALUES (%v, %v)", person.ID, weekNr)
	_, err := client.Exec(sql)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	//WriteSQLToFile(sql)
}