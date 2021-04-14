package sqlClient

import (
	"database/sql"
	"dbconn.com/apiClient"
	"fmt"
)

func(client *SQLClient)GetProviderByID(id int) Provider{
	sqlstr := fmt.Sprintf("SELECT * FROM Providers WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var provider Provider
	err := row.Scan(&provider.ID, &provider.PName)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return provider
}

func(client *SQLClient)CreateProviderEntry(provider apiClient.StreamingProvider) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create ProviderEntry "+provider.Provider_name)
	sqlstr := fmt.Sprintf("INSERT INTO Providers(id, pname) VALUES(%v,'%v')",provider.Provider_id, provider.Provider_name)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}