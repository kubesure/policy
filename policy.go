package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var mysqlsvc = os.Getenv("mysqlpolicysvc")

//Policy is the API output
type Policy struct {
	PolicyNumber int64
}

type request struct {
	QuoteNumber, ReceiptNumber string
}

func main() {
	log.Println("server policy starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health/poicies", createPolicy)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func createPolicy(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	r, err := marshalPolicy(string(body))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		pnumber, err := save(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			data, _ := json.Marshal(Policy{PolicyNumber: *pnumber})
			fmt.Fprintf(w, "%s", data)
		}
	}
}

func save(r *request) (*int64, error) {
	db, _ := sql.Open("mysql", "root:@tcp("+mysqlsvc+":3306)/test")
	defer db.Close()

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	statm, errstat := db.Prepare("INSERT policy set quote_no =?, receipt_no = ?")
	if errstat != nil {
		return nil, errstat
	}
	defer statm.Close()

	rs, errExec := statm.Exec(r.QuoteNumber, r.ReceiptNumber)
	if errExec != nil {
		return nil, err
	}

	polid, _ := rs.LastInsertId()
	return &polid, nil
}

func marshalPolicy(data string) (*request, error) {
	var r request
	err := json.Unmarshal([]byte(data), &r)
	if err != nil {
		log.Println("there was an error in marshalling request", err.Error())
		return nil, err
	}
	return &r, nil
}
