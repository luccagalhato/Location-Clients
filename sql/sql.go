package sql

import (
	"GoogleMAPS/models"
	"GoogleMAPS/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"sort"
	"time"

	_ "github.com/denisenkom/go-mssqldb" //bblablalba
)

// SQLStr ...
type SQLStr struct {
	url *url.URL
	db  *sql.DB
}

type distClient struct {
	Client   models.Client `json:"client,omitempty"`
	Distance float64       `json:"distance,omitempty"`
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
	rows, err := s.db.QueryContext(context.Background(), `SELECT TOP(10)A.CLIENTE_VAREJO AS NOME ,A.CODIGO_CLIENTE,A.ENDERECO,A.NUMERO,A.BAIRRO,A.CIDADE,A.UF,A.CEP,A.PAIS,LTRIM(RTRIM(B.LAT)) AS LAT, LTRIM(RTRIM(B.LONG)) AS LONG FROM CLIENTES_VAREJO A 
    LEFT JOIN
	(SELECT CODIGO_CLIENTE, CAST("01203" AS FLOAT) AS LAT, CAST("01204" AS FLOAT) AS LONG
	FROM
	(
	  SELECT CODIGO_CLIENTE, VALOR_PROPRIEDADE, PROPRIEDADE
	  FROM LINX_TBFG..PROP_CLIENTES_VAREJO
	  WHERE PROPRIEDADE IN ('01203', '01204')
	) d
	PIVOT
	(
	  max(VALOR_PROPRIEDADE)
	  FOR PROPRIEDADE IN ("01203", "01204")
	) piv) B ON A.CODIGO_CLIENTE=B.CODIGO_CLIENTE
    WHERE A.PF_PJ= '0' AND A.CODIGO_CLIENTE != '' AND  NOT(B.LAT IS NULL OR B.LONG IS NULL)`, nil)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	distClients := make(distClients, 0)

	for rows.Next() {
		client := models.Client{}

		if err := rows.Scan(&client.Nome, &client.CodClient, &client.Endereco, &client.Numero, &client.Bairro, &client.Cidade, &client.Uf, &client.Cep, &client.Pais, &client.Lat, &client.Long); err != nil {
			fmt.Println(err)
			continue
		}
		// rst = append(rst, &client)
		distClients = append(distClients, distClient{
			Client:   client,
			Distance: utils.CalcDistancia(lat, long, *client.Lat, *client.Long),
		})
	}

	sort.Sort(distClients)

	// for i := 0; i < len(distClients); i++ {
	// 	fmt.Printf("%+v\n", distClients[i])
	// }

	return distClients, nil

}

const (
	update = "UPDATE %s SET VALOR_PROPRIEDADE=%s WHERE PROPRIEDAD=%s AND CLIENTE_ATACADO=%s"
)
const table = "LINX_TBFG..PROP_CLIENTES_ATACADO"

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

// Ping ...
func (s *SQLStr) Ping() {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err := s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		log.Panicln("reconnecting")
		if err != s.connect() {
			log.Println(err)
		}
	}

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
