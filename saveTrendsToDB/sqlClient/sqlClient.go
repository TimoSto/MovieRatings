package sqlClient

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"strings"
)

type SQLClient struct {
	DB *sql.DB
}

type Config struct {
	MySqlUser string
	MySqlPassword string
	MySqlPort string
}

func(client *SQLClient)Exec(sqlstr string) (sql.Result,error) {
	res, err := client.DB.Exec(sqlstr)
	if err == nil {
		WriteSQLToFile(sqlstr)
	}
	return res, err
}

func WriteSQLToFile(sqlstr string){
	f, err := os.OpenFile("../database/FILLDB.sql",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//sqlstr = strings.Replace(sqlstr,"'","\\'",-1)
	fc, err := ioutil.ReadFile("../database/FILLDB.sql")
	if err != nil {
		panic(err)
	}
	index := strings.Index(string(fc), sqlstr)
	if index == -1 {
		if _, err := f.WriteString(sqlstr +";\n"); err != nil {
			panic(err)
		}
	}
}

func(client *SQLClient)EstablishConnectionToDB(config Config) {
	fmt.Println("Trying to connect to to mySQL-DB...")
	var err error
	client.DB, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/movieratings", config.MySqlUser, config.MySqlPassword, config.MySqlPort))
	if err != nil {
		panic(err)
	}
}