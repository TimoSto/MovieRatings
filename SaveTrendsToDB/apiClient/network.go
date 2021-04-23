package apiClient

type Network struct {
	Name           string `json:name`
	ID             int    `json:id`
	Logo_Path      string `json:logo_path`
	Origin_country string `json:origin_country`
}