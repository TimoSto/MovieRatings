package sqlClient

import (
	"database/sql"
	"dbconn.com/apiClient"
	"fmt"
)

func(client *SQLClient)GetNetworkByID(id int) Network{
	sqlstr := fmt.Sprintf("SELECT * FROM Networks WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var network Network
	err := row.Scan(&network.ID, &network.NName, &network.Logo, &network.OriginCountry)
	if err != nil && err != sql.ErrNoRows{
		panic(err)
	}
	return network
}

func(client *SQLClient)CreateNetworkEntry(network apiClient.Network) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create Networks-Entry ", network.ID)
	if client.GetCountryByID(network.Origin_country).ID.Valid == false {
		client.CreateCountryEntry(apiClient.Country{network.Origin_country, ""})
	}
	sqlstr := fmt.Sprintf("INSERT INTO Networks(id, nname, logo, originCountry) VALUES('%v','%v', '%v', '%v')", network.ID, network.Name, network.Logo_Path, network.Origin_country)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}