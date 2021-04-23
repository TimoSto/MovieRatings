package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
)

type Country struct {
	ID sql.NullString `json:id`
	CName sql.NullString `json:cname`
}

func(client *SQLClient)ExtendCountriesTable(countries []apiClient.Country) {
	fmt.Println("Countries")
	for _, c := range countries {
		if _,n := client.GetCountryByID(c.ISO_3166_1); n == -1 {
			client.CreateCountryEntry(c)
		}
	}
}

func(client *SQLClient)GetCountryByID(id string) (Country, int){
	sqlstr := fmt.Sprintf("SELECT * FROM Countries WHERE id='%v'", id)
	row := client.DB.QueryRow(sqlstr)
	var country Country
	err := row.Scan(&country.ID, &country.CName)
	if err == sql.ErrNoRows {
		return Country{}, -1
	}
	if err != nil {
		panic(err)
	}
	return country, 1
}

func(client *SQLClient)CreateCountryEntry(country apiClient.Country) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create CountryEntry "+country.Name)
	sqlstr := fmt.Sprintf("INSERT INTO Countries(id, cname) VALUES('%v','%v')",country.ISO_3166_1, country.Name)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}