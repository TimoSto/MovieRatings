package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
)

type Network struct {
	ID            sql.NullInt64  `json:id`
	NName         sql.NullString `json:nname`
	Logo          sql.NullString `json:logo`
	OriginCountry sql.NullString `json:originCountry`
}

func(client *SQLClient)ExtendNetworkTable(networks []apiClient.Network) {

	for _,netw := range networks {
		if _,n := client.GetNetworkByID(netw.ID); n == -1 {
			client.CreateNetworkEntry(netw)
		}
	}
}

func(client *SQLClient)GetNetworkByID(id int) (Network, int){
	sqlstr := fmt.Sprintf("SELECT * FROM Networks WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var network Network
	err := row.Scan(&network.ID, &network.NName, &network.Logo, &network.OriginCountry)
	if err == sql.ErrNoRows {
		return Network{},-1
	}
	if err != nil {
		panic(err)
	}
	return network, 1
}

func(client *SQLClient)CreateNetworkEntry(network apiClient.Network) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create Networks-Entry ", network.ID)
	if _,n := client.GetCountryByID(network.Origin_country); n == -1 {
		client.CreateCountryEntry(apiClient.Country{network.Origin_country, ""})
	}
	sqlstr := fmt.Sprintf("INSERT INTO Networks(id, nname, logo, originCountry) VALUES('%v','%v', '%v', '%v')", network.ID, network.Name, network.Logo_Path, network.Origin_country)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}