package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	queryCreateTbl = "CREATE TABLE IF NOT EXISTS values (number INT);"
	querySelectAll = "SELECT DISTINCT number FROM values;"
	queryInsert = "INSERT INTO values(number) VALUES($1);"
)

const (
	dbhost     = "postgres"
	dbport     = 5432
	dbuser     = "postgres"
	dbpassword = "postgres_password"
	dbname   = "postgres"
)
type payload struct {
	index string `json:"index"`
}

func main() {
	redisClient := NewRedisClient()
	pgClient := NewPGClient()
	defer pgClient.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hi")
		fmt.Fprintf(w, "Hi")
	})
	http.HandleFunc("/values/all", func(w http.ResponseWriter, r *http.Request) {
		stmt, err := pgClient.Prepare(querySelectAll)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()
		rows, err := stmt.Query()
		if err != nil {
			panic(err)
		}
		values := make([]int, 0)
		for rows.Next() {
			var value int
			rows.Scan(&value)
			values = append(values, value)
		}
		jsonRes, _ := json.Marshal(values)
		w.Write(jsonRes)
	})
	http.HandleFunc("/values/current", func(w http.ResponseWriter, r *http.Request) {
		redisVals := redisClient.HGetAll("values")
		result, err := redisVals.Result()
		if err != nil {
			panic(err)
		}
		jsonRes, _ := json.Marshal(result)
		w.Write(jsonRes)
		//fmt.Println(redisVals.String())
		//fmt.Fprintf(w, redisVals.String())
	})

	http.HandleFunc("/values", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			fmt.Fprintf(w, "GET")
			return
		}
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var dat map[string]string

		if err := json.Unmarshal(reqBody, &dat); err != nil {
			panic(err)
		}

		if _, found := dat["index"]; !found {
			panic("index not sent")
		}
		index, err := strconv.Atoi(dat["index"])
		if err != nil {
			log.Fatal(err)
		}
		if index > 40 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("index too high"))
			return
		}
		redisClient.HSet("values", dat["index"], "Nothing yet")
		redisClient.Publish("insert", index)
		q, err := pgClient.Prepare(queryInsert)
		if err != nil {
			log.Fatal(err)
		}
		defer q.Close()
		q.Exec(index)
		//PG INSERT array of index
	})

	http.ListenAndServe(":5000", nil)
	fmt.Println("Listening...")
}

func NewPGClient() *sql.DB {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbuser,
		dbpassword,
		dbhost,
		dbport,
		dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	//defer db.Close()
	stmt, err := db.Prepare(queryCreateTbl)
	if err != nil {
		panic(err)
	}
	stmt.Exec()
	defer stmt.Close()
	return db
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}
