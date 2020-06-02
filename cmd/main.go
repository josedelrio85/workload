package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	workload "github.com/bysidecar/workload/pkg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	log.Println("Workload API started")

	port := getSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	database := workload.Database{
		Host:      getSetting("DB_HOST"),
		Port:      portInt,
		User:      getSetting("DB_USER"),
		Pass:      getSetting("DB_PASS"),
		Dbname:    getSetting("DB_NAME"),
		ParseTime: "True",
		Charset:   "utf8",
		Loc:       "Local",
	}

	handler := workload.Handler{
		DbHandler: &database,
		ActiveModels: []workload.Modelable{
			&workload.Project{
				Type: 1,
			},
			&workload.Work{
				Type: 2,
			},
			&workload.Status{
				Type: 3,
			},
			&workload.Person{
				Type: 4,
			},
		},
	}

	if err := database.Open(); err != nil {
		log.Fatalf("Error opening database %v", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Error creating tables, err :%v", err)
	}

	router := mux.NewRouter()
	sub := router.PathPrefix("/workload").Subrouter()
	sub.Handle("/get", handler.Get()).Methods(http.MethodPost)
	sub.Handle("/create", handler.Put()).Methods(http.MethodPost)
	sub.Handle("/update", handler.Update()).Methods(http.MethodPost)
	sub.Handle("/delete", handler.Delete()).Methods(http.MethodDelete)

	log.Println("starting web server...")
	log.Fatal(http.ListenAndServe(":9001", cors.Default().Handler(sub)))
}

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init erro, %s setting was not found", setting)
	}
	return value
}
