package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var mysqlsvc = os.Getenv("mysqlpolicysvc")

//Policy is the API output
type policy struct {
	PolicyNumber int64
}

type request struct {
	QuoteNumber, ReceiptNumber string
}

func main() {
	log.Println("server policy starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health/poicies", createPolicy)
	srv := http.Server{Addr: ":8080", Handler: mux}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Print("shutting down policy server...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
	//log.Fatal(http.ListenAndServe(":8000", mux))
}

func validateReq(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodPost {
		log.Println("invalid method ", req.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Invalid method %s", req.Method)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		log.Println("invalid content type ", req.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("Invalid content-type require %s", "application/json")
	}
	return nil
}

func createPolicy(w http.ResponseWriter, req *http.Request) {

	if err := validateReq(w, req); err != nil {
		return
	}

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
			data, _ := json.Marshal(policy{PolicyNumber: *pnumber})
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
		return nil, err
	}
	return &r, nil
}
