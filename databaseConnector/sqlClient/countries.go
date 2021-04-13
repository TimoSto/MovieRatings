package sqlClient

import (
	"database/sql"
	"dbconn.com/apiClient"
	"fmt"
)

func(client *SQLClient)GetCountryByID(id string) Country{
	sqlstr := fmt.Sprintf("SELECT * FROM Countries WHERE id='%v'", id)
	row := client.DB.QueryRow(sqlstr)
	var country Country
	err := row.Scan(&country.ID, &country.CName)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return country
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