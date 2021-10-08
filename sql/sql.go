package sql

import (
	maps "GoogleMAPS/googlemaps"
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
	rows, err := s.db.QueryContext(context.Background(), `SELECT LTRIM(RTRIM(COALESCE(NOME_CLIFOR,''))) AS NOME_CLIFOR,LTRIM(RTRIM(COALESCE(A.ENDERECO,''))) AS ENDERECO,LTRIM(RTRIM(COALESCE(A.NUMERO,''))) AS NUMERO,LTRIM(RTRIM(COALESCE(A.BAIRRO,''))) AS BAIRRO,LTRIM(RTRIM(COALESCE(A.CIDADE,''))) AS CIDADE,LTRIM(RTRIM(COALESCE(A.UF,''))) AS UF,LTRIM(RTRIM(COALESCE(A.CEP,''))) AS CEP,LTRIM(RTRIM(COALESCE(A.PAIS,''))) AS PAIS,LTRIM(RTRIM(COALESCE(A.CLIFOR,''))) AS CLIFOR, LTRIM(RTRIM(B.LAT)) AS LAT, LTRIM(RTRIM(B.LONG)) AS LONG
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
		client := models.Client{}

		if err := rows.Scan(&client.Nome, &client.Endereco, &client.Numero, &client.Bairro, &client.Cidade, &client.Uf, &client.Cep, &client.Pais, &client.Clifor, &client.Lat, &client.Long); err != nil {
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

func (s *SQLStr) SearchNewClient() error {
	rows, err := s.db.QueryContext(context.Background(), `SELECT LTRIM(RTRIM(COALESCE(NOME_CLIFOR,''))) AS NOME_CLIFOR,LTRIM(RTRIM(COALESCE(A.ENDERECO,''))) AS ENDERECO,LTRIM(RTRIM(COALESCE(A.NUMERO,''))) AS NUMERO,LTRIM(RTRIM(COALESCE(A.BAIRRO,''))) AS BAIRRO,LTRIM(RTRIM(COALESCE(A.CIDADE,''))) AS CIDADE,LTRIM(RTRIM(COALESCE(A.UF,''))) AS UF,LTRIM(RTRIM(COALESCE(A.CEP,''))) AS CEP,LTRIM(RTRIM(COALESCE(A.PAIS,''))) AS PAIS,LTRIM(RTRIM(COALESCE(A.CLIFOR,''))) AS CLIFOR, LTRIM(RTRIM(B.LAT)) AS LAT, LTRIM(RTRIM(B.LONG)) AS LONG, b.DATA_PARA_TRANSFERENCIA
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
	WHERE A.DATA_PARA_TRANSFERENCIA>B.DATA_PARA_TRANSFERENCIA OR B.DATA_PARA_TRANSFERENCIA IS NULL`, nil)
	if err != nil {
		// fmt.Println(err)
		return err
	}
	for rows.Next() {
		client := models.Client{}
		if err := rows.Scan(&client.Nome, &client.Endereco, &client.Numero, &client.Bairro, &client.Cidade, &client.Uf, &client.Cep, &client.Pais, &client.Clifor, &client.Lat, &client.Long, &client.Data); err != nil {
			fmt.Println(err)
			continue
		}
		lat, long, err := maps.RequestMapsNewclientRoutine(client)
		if err != nil || lat == 0.00 || long == 0.00 {
			log.Println(err, "cliente:", client.Nome)
			continue
		}
		if client.Data != nil {
			s.UpdateRow(client.Nome, fmt.Sprintf("%f", lat), "01203")
			s.UpdateRow(client.Nome, fmt.Sprintf("%f", long), "01204")
			continue
		}
		s.InsertRow(fmt.Sprintf("%f", lat), fmt.Sprintf("%f", long), client.Nome)
	}
	return nil
}
func (s *SQLStr) InsertRow(lat, long string, nome string) {
	query := fmt.Sprintf(`insert into LINX_TBFG..PROP_CLIENTES_ATACADO (PROPRIEDADE,CLIENTE_ATACADO,ITEM_PROPRIEDADE, VALOR_PROPRIEDADE)
	VALUES 
	('01203', '%s', 1, '%s'),
	('01204', '%s', 1, '%s');`, nome, lat, nome, long)
	_, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		fmt.Println(err, query)
		return
	}
}
func (s *SQLStr) UpdateRow(cliente, value, condition string) {
	query := fmt.Sprintf(update, table, value, condition, cliente)
	_, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		fmt.Println(err, query)
		return
	}
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

//Ping ...
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
