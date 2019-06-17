package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var mysqlsvc = os.Getenv("mysqlpolicysvc")

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

//Policy is the API output
type policy struct {
	PolicyNumber int64
}

type request struct {
	QuoteNumber, ReceiptNumber string
}

type eventpolicyissued struct {
	PolicyNumber               int64
	QuoteNumber, ReceiptNumber string
}

type erroresponse struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

//Error Code Enum
const (
	SystemErr = iota
	InputJSONInvalid
	AgeRangeInvalid
	RiskDetailsInvalid
	InvalidRestMethod
	InvalidContentType
)

func main() {
	log.Debug("server policy starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health/poicies", createPolicy)
	srv := http.Server{Addr: ":8000", Handler: mux}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Debug("shutting down policy server...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func validateReq(w http.ResponseWriter, req *http.Request) (*request, *erroresponse) {
	if req.Method != http.MethodPost {
		log.Debug("invalid method ", req.Method)
		return nil, &erroresponse{Code: InvalidRestMethod, Message: fmt.Sprintf("Invalid method %s", req.Method)}
	}

	if req.Header.Get("Content-Type") != "application/json" {
		log.Debug("invalid content type ", req.Header.Get("Content-Type"))
		msg := fmt.Sprintf("Invalid content-type %s require %s", req.Header.Get("Content-Type"), "application/json")
		return nil, &erroresponse{Code: InvalidContentType, Message: msg}
	}

	body, _ := ioutil.ReadAll(req.Body)
	r, err := marshalPolicy(string(body))

	if err != nil {
		return nil, err
	}

	return r, nil
}

func createPolicy(w http.ResponseWriter, req *http.Request) {
	r, err := validateReq(w, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(err)
		fmt.Fprintf(w, "%s", data)
	} else {
		pnumber, err := save(r)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			data, _ := json.Marshal(policy{PolicyNumber: *pnumber})
			fmt.Fprintf(w, "%s", data)
		}
	}
}

func save(r *request) (*int64, error) {
	db, _ := sql.Open("mysql", "root:admin@tcp("+mysqlsvc+":3306)/policy")
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

func marshalPolicy(data string) (*request, *erroresponse) {
	var r request
	err := json.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, &erroresponse{Code: SystemErr, Message: "input invalid"}
	}

	if len(r.QuoteNumber) == 0 || len(r.ReceiptNumber) == 0 {
		return nil, &erroresponse{Code: InputJSONInvalid, Message: "Invalid Input"}
	}

	return &r, nil
}
