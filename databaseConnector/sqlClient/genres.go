package sqlClient

import (
	"database/sql"
	"dbconn.com/apiClient"
	"fmt"
)

func(client *SQLClient)CreateGenreEntry(genre apiClient.Genre) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create GenreEntry "+genre.Name)
	sqlstr := fmt.Sprintf("INSERT INTO Genres(id, genre) VALUES(%v,'%v')",genre.ID, genre.Name)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)GetGenreByID(id int) Genre{
	sqlstr := fmt.Sprintf("SELECT * FROM Genres WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var genre Genre
	err := row.Scan(&genre.ID, &genre.Genre)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return genre
}