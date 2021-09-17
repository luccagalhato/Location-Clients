package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"sort"

	_ "github.com/denisenkom/go-mssqldb" //bblablalba
)

// SQLStr ...
type SQLStr struct {
	url *url.URL
	db  *sql.DB
}

// ConnectLinx
func (s *SQLStr) ConnectLinx(callback func(client Clients, lat, long float64) error) error {

	// rst := make([]*Clients, 0)
	rows, err := s.db.QueryContext(context.Background(), `select LTRIM(RTRIM(NOME_CLIFOR)) AS NOME_CLIFOR, LTRIM(RTRIM(ENDERECO)) AS ENDERECO,LTRIM(RTRIM(COALESCE(NUMERO, ''))) AS NUMERO, LTRIM(RTRIM(BAIRRO)) AS BAIRRO, LTRIM(RTRIM(CIDADE)) AS CIDADE, LTRIM(RTRIM(CEP)) AS CEP, LTRIM(RTRIM(PAIS)) AS PAIS, LTRIM(RTRIM(CLIFOR)) from LINX_TBFG..CADASTRO_CLI_FOR
	WHERE INDICA_CLIENTE = '1'`, nil)
	if err != nil {
		// fmt.Println(err)
		return err
	}
	for rows.Next() {
		client := Clients{}

		if err := rows.Scan(&client.Nome, &client.Endereco, &client.Numero, &client.Bairro, &client.Cidade, &client.Cep, &client.Pais, &client.Clifor); err != nil {
			fmt.Println(err, client.Clifor)
			continue
		}
		// rst = append(rst, &client)
		lat, long, err := requestMaps(client)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := callback(client, lat, long); err != nil {
			fmt.Println(err)
		}
	}
	return nil

}

type distClient struct {
	Client   Clients `json:"client,omitempty"`
	Distance float64 `json:"distance,omitempty"`
}

type distClients []distClient

func (d distClients) Len() int {
	return len(d)
}

func (d distClients) Less(i, j int) bool {
	return d[i].Distance < d[j].Distance
}

func (d distClients) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (s *SQLStr) CompareRegion(lat, long float64) ([]distClient, error) {

	// rst := make([]*Clients, 0)
	rows, err := s.db.QueryContext(context.Background(), `SELECT LTRIM(RTRIM(NOME_CLIFOR)) AS NOME_CLIFOR,LTRIM(RTRIM(A.ENDERECO)) AS ENDERECO,LTRIM(RTRIM(A.NUMERO)) AS NUMERO,LTRIM(RTRIM(A.BAIRRO)) AS BAIRRO,LTRIM(RTRIM(A.CIDADE)) AS CIDADE,LTRIM(RTRIM(A.UF)) AS UF,LTRIM(RTRIM(A.CEP)) AS CEP,LTRIM(RTRIM(A.PAIS)) AS PAIS,LTRIM(RTRIM(A.CLIFOR)) AS CLIFOR, LTRIM(RTRIM(B.LAT)) AS LAT, LTRIM(RTRIM(B.LONG)) AS LONG
	FROM (SELECT * FROM LINX_TBFG..CADASTRO_CLI_FOR WHERE INDICA_CLIENTE='1' AND PJ_PF = '1') A 
	LEFT JOIN
	(SELECT CLIENTE_ATACADO, CAST("01203" AS FLOAT) AS LAT, CAST("01204" AS FLOAT) AS LONG, DATA_PARA_TRANSFERENCIA
	FROM
	(
	  SELECT CLIENTE_ATACADO, VALOR_PROPRIEDADE, PROPRIEDADE, DATA_PARA_TRANSFERENCIA
	  FROM LINX_TBFG..PROP_CLIENTES_ATACADO
	  WHERE PROPRIEDADE IN ('01203', '01204')
	) d
	PIVOT
	(
	  max(VALOR_PROPRIEDADE)
	  FOR PROPRIEDADE IN ("01203", "01204")
	) piv) B ON A.NOME_CLIFOR=B.CLIENTE_ATACADO
	WHERE NOT(B.LAT IS NULL OR B.LONG IS NULL)`, nil)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	distClients := make(distClients, 0)

	for rows.Next() {
		client := Clients{}

		if err := rows.Scan(&client.Nome, &client.Endereco, &client.Numero, &client.Bairro, &client.Cidade, &client.Uf, &client.Cep, &client.Pais, &client.Clifor, &client.Latitude, &client.Longitude); err != nil {
			fmt.Println(err)
			continue
		}
		// rst = append(rst, &client)
		distClients = append(distClients, distClient{
			Client:   client,
			Distance: calcDistancia(lat, long, client.Latitude, client.Longitude),
		})
	}

	sort.Sort(distClients)

	if len(distClients) > 10 {
		distClients = distClients[:10]
	}

	for i := 0; i < len(distClients); i++ {
		fmt.Printf("%+v\n", distClients[i])
	}

	return distClients, nil

}

// func (s *SQLStr) SearchClient(code string) {
// 	// rst := make([]*Clients, 0)
// 	rows, err := s.db.QueryContext(context.Background(), `SELECT NOME_CLIFOR,A.ENDERECO,A.NUMERO,A.BAIRRO,A.CIDADE,A.UF,A.CEP,A.PAIS,A.CLIFOR, B.LAT, B.LONG
// FROM (SELECT * FROM LINX_TBFG..CADASTRO_CLI_FOR WHERE INDICA_CLIENTE='1' AND PJ_PF = '1') A
// LEFT JOIN
// (SELECT CLIENTE_ATACADO, CAST("01203" AS FLOAT) AS LAT, CAST("01204" AS FLOAT) AS LONG, DATA_PARA_TRANSFERENCIA
// FROM
// (
//   SELECT CLIENTE_ATACADO, VALOR_PROPRIEDADE, PROPRIEDADE, DATA_PARA_TRANSFERENCIA
//   FROM LINX_TBFG..PROP_CLIENTES_ATACADO
//   WHERE PROPRIEDADE IN ('01203', '01204')
// ) d
// PIVOT
// (
//   max(VALOR_PROPRIEDADE)
//   FOR PROPRIEDADE IN ("01203", "01204")
// ) piv) B ON A.NOME_CLIFOR=B.CLIENTE_ATACADO
// WHERE NOT(B.LAT IS NULL OR B.LONG IS NULL)`, nil)
// 	if err != nil {
// 		// fmt.Println(err)
// 		return nil, err
// 	}
// 	rst := make(Clients, 0)
// 	for rows.Next() {
// 		client := Clients{}

// 		if err := rows.Scan(&client.Nome, &client.Endereco, &client.Numero, &client.Bairro, &client.Cidade, &client.Uf, &client.Cep, &client.Pais, &client.Clifor, &client.Latitude, &client.Longitude); err != nil {
// 			fmt.Println(err)
// 			continue

// 		}
// 		rst := append(rst, &client)

// 	}
// 	return
// }

//InsertSql ...
func (s *SQLStr) InsertSql(clifor string, nome string, endereco string, numero string, bairro string, cep string, cidade string, pais string, latitude float64, longitude float64) error {
	_, err := s.db.QueryContext(context.Background(), fmt.Sprintf(`INSERT INTO Manchester_Group..CADASTRO_CLIENTES (CLIFOR,NOME,ENDERECO,NUMERO,BAIRRO,CEP,CIDADE,PAIS,LATITUDE,LONGITUDE)
	VALUES ('%s','%s','%s','%s','%s','%s','%s','%s','%f','%f');`, clifor, nome, endereco, numero, bairro, cep, cidade, pais, latitude, longitude))
	return err
}

func MakeSQL(host, port, username, password string) (*SQLStr, error) {

	s := &SQLStr{}
	s.url = &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%s", host, port),
		RawQuery: url.Values{}.Encode(),
	}
	return s, s.connect()
}

func (s *SQLStr) connect() error {
	var err error
	if s.db, err = sql.Open("sqlserver", s.url.String()); err != nil {
		return err
	}
	return s.db.PingContext(context.Background())
}

// func (s *SQLStr) disconnect() error {
// 	return s.db.Close()
// }
