package apiClient

type Genre struct {
	ID   int    `json:id`
	Name string `json:name`
}

func (client *APIClient)GetGenres(movies []Movie) []Genre {
	var genres []Genre

	for _, movie := range movies {
		//Cast zur Personen-Liste hinzufügen, falls noch nicht vorhanden
		for _,genre := range movie.Genres {
			if findGenreInSlice(genres, genre.ID) == -1 {
				genres = append(genres, genre)
			}
		}
	}

	return genres
}

func findGenreInSlice(arr []Genre, id int) int {
	for i,g := range arr {
		if g.ID == id {
			return i
		}
	}
	return -1
}