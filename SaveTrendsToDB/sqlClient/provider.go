package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
)

type Provider struct {
	ID sql.NullInt64 `json:id`
	PName sql.NullString `json:pname`
	Service sql.NullString `json:service`
}

func(client *SQLClient)ExtendProviderTable(providers []apiClient.Provider) {
	for _,p := range providers {
		if _,n := client.GetProviderByID(p.Provider_id); n == -1 {
			client.CreateProviderEntry(p)
		}
	}
}

func(client *SQLClient)GetProviderByID(id int) (Provider, int){
	sqlstr := fmt.Sprintf("SELECT * FROM Providers WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var provider Provider
	err := row.Scan(&provider.ID, &provider.PName)
	if err == sql.ErrNoRows {
		return Provider{}, -1
	}
	if err != nil {
		panic(err)
	}
	return provider, 1
}

func(client *SQLClient)CreateProviderEntry(provider apiClient.Provider) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create ProviderEntry "+provider.Provider_name)
	sqlstr := fmt.Sprintf("INSERT INTO Providers(id, pname) VALUES(%v,'%v')",provider.Provider_id, provider.Provider_name)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}