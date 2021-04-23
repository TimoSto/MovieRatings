package apiClient

type Network struct {
	Name           string `json:name`
	ID             int    `json:id`
	Logo_Path      string `json:logo_path`
	Origin_country string `json:origin_country`
}

func(client *APIClient)GetNetworksForTVTrends(series []Series) []Network{
	var networks []Network

	for _,s := range series {
		for _,netw := range s.Networks {
			if findNetworkInSlice(networks, netw.ID) == -1 {
				networks = append(networks, netw)
			}
		}
	}

	return networks
}

func findNetworkInSlice(arr []Network, id int) int {
	for i,n := range arr {
		if n.ID == id {
			return i
		}
	}
	return -1
}