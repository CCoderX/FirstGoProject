package awesomeProject

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"taher_first_go/handlers"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initializeRoutes() {
	h := &handlers.Handler{}
	a.Router.HandleFunc("/getCars", h.GetCar).Methods("GET")
	a.Router.HandleFunc("/addCar", h.AddCar).Methods("POST")
	a.Router.HandleFunc("/getCarId/{id:[0-9]+}", h.GetCar).Methods("GET")
	a.Router.HandleFunc("/uodateCar/{id:[0-9]+}", h.UpdateCar).Methods("PUT")
	a.Router.HandleFunc("/deleteCar/{id:[0-9]+}", h.DeleteCar).Methods("DELETE")
}

func (a *App) getCars(writer http.ResponseWriter, request *http.Request) {
	count, _ := strconv.Atoi(request.FormValue("count"))
	start, _ := strconv.Atoi(request.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getCars(a.DB, start, count)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, products)
}

func (a *App) addCar(w http.ResponseWriter, r *http.Request) {
	var p Car
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createCar(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) getCar(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid Car ID")
		return
	}

	p := Car{ID: id}
	if err := p.getCar(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(writer, http.StatusNotFound, "Product not found")
		default:
			respondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(writer, http.StatusOK, p)
}

func (a *App) updateCar(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid Car ID")
		return
	}

	var car Car
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&car); err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer request.Body.Close()
	car.ID = id

	if err := car.updateCar(a.DB); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, car)
}

func (a *App) deleteCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := Car{ID: id}
	if err := p.deleteCar(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
