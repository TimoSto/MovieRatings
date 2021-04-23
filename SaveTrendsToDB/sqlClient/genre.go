package sqlClient

import (
	"database/sql"
	"fmt"
	"savetrends.com/apiClient"
)

type Genre struct {
	ID sql.NullInt64 `json:id`
	Genre sql.NullString `json:genre`

}

func(client *SQLClient)ExtendGenresTable(genres []apiClient.Genre) {
	for _,genre := range genres {
		if _,n := client.GetGenreByID(genre.ID); n == -1 {
			client.CreateGenreEntry(genre)
		}
	}
}

func(client *SQLClient)CreateGenreEntry(genre apiClient.Genre) {
	//Eintrag für Film in SQL-DB hinzufügen
	fmt.Println("Create GenreEntry "+genre.Name)
	sqlstr := fmt.Sprintf("INSERT INTO Genres(id, genre) VALUES(%v,'%v')",genre.ID, genre.Name)
	_, err := client.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func(client *SQLClient)GetGenreByID(id int) (Genre, int){
	sqlstr := fmt.Sprintf("SELECT * FROM Genres WHERE id=%v", id)
	row := client.DB.QueryRow(sqlstr)
	var genre Genre
	err := row.Scan(&genre.ID, &genre.Genre)
	if err == sql.ErrNoRows {
		return Genre{}, -1
	} else if err != nil {
		panic(err)
	}
	return genre, 1
}