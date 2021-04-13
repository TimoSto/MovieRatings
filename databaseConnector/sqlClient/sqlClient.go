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

type Movie struct {
	ID          int             `json:id`
	Title       sql.NullString  `json:title`
	Overview    sql.NullString  `json:overview`
	Popularity  sql.NullFloat64 `json:popularity`
	Revenue     sql.NullString  `json:revenue`
	PosterPath  sql.NullString  `json:posterPath`
	ReleaseDate sql.NullString  `json:release_date`
	VoteAVG     sql.NullFloat64 `json:voteAvg`
	VoteCount   sql.NullInt64   `json:VoteCount`
	Runtime     sql.NullInt64   `json:runtime`
	Tagline     sql.NullString  `json:tagline`
}

type Series struct {
	ID           int         `json:id`
	Title        sql.NullString  `json:title`
	Overview    sql.NullString  `json:overview`
	Popularity  sql.NullFloat64 `json:popularity`
	Seasons     sql.NullInt64   `json:seasons`
	Episodes    sql.NullInt64   `json:episodes`
	Poster      sql.NullString  `json:posterPath`
	ReleaseDate sql.NullString  `json:release_date`
	VoteAVG     sql.NullFloat64 `json:voteAvg`
	VoteCount   sql.NullInt64   `json:VoteCount`
	FirstAit    sql.NullString  `json:firstAir`
	LastAir     sql.NullString  `json:lastAir`
	Tagline     sql.NullString  `json:tagline`
}

type Person struct {
	ID sql.NullInt64 `json:id`
	Name sql.NullString `json:name`
	Birthday sql.NullString `json:birthday`
	Deathday sql.NullString `json:deathday`
	Popularity sql.NullFloat64 `json:popularity`
	ProfilePath sql.NullString `json:profilePath`
	Gender sql.NullInt64 `json:gender`
	Profession sql.NullString `json:profession`
}

type Country struct {
	ID sql.NullString `json:id`
	CName sql.NullString `json:cname`
}

type Network struct {
	ID            sql.NullInt64  `json:id`
	NName         sql.NullString `json:nname`
	Logo          sql.NullString `json:logo`
	OriginCountry sql.NullString `json:originCountry`
}

type Genre struct {
	ID sql.NullInt64 `json:id`
	Genre sql.NullString `json:genre`

}

func(client *SQLClient)EstablishConnectionToDB() {
	fmt.Println("Trying to connect to to mySQL-DB...")
	var err error
	client.DB, err = sql.Open("mysql", "root:Pa$$w0rd@tcp(127.0.0.1:3306)/movieratings")
	if err != nil {
		panic(err)
	}
}















