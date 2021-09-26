package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/emicklei/go-restful"
	"github.com/mhelmeck/RailAPI/dbutils"
)

var DB *sql.DB

type TrainResource struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}
type StationResource struct {
	ID          int
	Name        string
	OpeningTime time.Time
	ClosingTime time.Time
}
type ScheduleResource struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime time.Time
}

func (t *TrainResource) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/v1/trains").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))

	container.Add(ws)
}

func (t TrainResource) getTrain(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("train-id")

	err := DB.QueryRow("select ID, DRIVER_NAME, OPERATING_STATUS FROM train where id=?", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		log.Println(err)

		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusNotFound, "Train could not befound.")
	} else {
		res.WriteEntity(t)
	}
}

func (t TrainResource) createTrain(req *restful.Request, res *restful.Response) {
	log.Println(req.Request.Body)

	decoder := json.NewDecoder(req.Request.Body)
	var tr TrainResource
	decoder.Decode(&tr)

	log.Println(tr.DriverName, tr.OperatingStatus)
	statement, _ := DB.Prepare("insert into train (DRIVER_NAME, OPERATING_STATUS) values (?, ?)")
	result, err := statement.Exec(tr.DriverName, tr.OperatingStatus)
	if err == nil {
		newID, _ := result.LastInsertId()
		tr.ID = int(newID)

		res.WriteHeaderAndEntity(http.StatusCreated, tr)
	} else {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (t TrainResource) removeTrain(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("train-id")

	statement, _ := DB.Prepare("delete from train where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		res.WriteHeader(http.StatusOK)
	} else {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func main() {
	var err error
	DB, err = sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("Driver creation failed!")
	}

	dbutils.Initialize(DB)

	rsContainer := restful.NewContainer()
	rsContainer.Router(restful.CurlyRouter{})

	t := TrainResource{}
	t.Register(rsContainer)

	log.Printf("start listening on localhost:8000")
	server := &http.Server{Addr: ":8000", Handler: rsContainer}

	log.Fatal(server.ListenAndServe())
}
