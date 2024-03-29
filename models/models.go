package models

// Client ...
type Client struct {
	CodClient    string   `json:"CODIGO_CLIENTE,omitempty"`
	Nome         string   `json:"NOME,omitempty"`
	Endereco     string   `json:"ENDERECO,omitempty"`
	Numero       string   `json:"NUMERO,omitempty"`
	Bairro       string   `json:"BAIRRO,omitempty"`
	Cidade       string   `json:"CIDADE,omitempty"`
	Uf           string   `json:"UF,omitempty"`
	Cep          string   `json:"CEP,omitempty"`
	Pais         string   `json:"PAIS,omitempty"`
	Lat          *float64 `json:"LAT,omitempty"`
	Long         *float64 `json:"LONG,omitempty"`
	IndicaFilial string   `json:"INDICA_FILIAL,omitempty"`
}
type AutoGenerated struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}
type Response struct {
	Status string
	Error  string
	Data   interface{}
}
type Street struct {
	Endereco string `json:"endereco,omitempty"`
	Numero   string `json:"numero,omitempty"`
	Bairro   string `json:"bairro,omitempty"`
	Cidade   string `json:"cidade,omitempty"`
	Uf       string `json:"uf,omitempty"`
	Cep      string `json:"cep,omitempty"`
	Pais     string `json:"pais,omitempty"`
}
