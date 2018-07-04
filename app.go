package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	. "movies-restapi-master/config"
	. "movies-restapi-master/dao"
	. "movies-restapi-master/models"

	"github.com/gorilla/mux"
)

var config = Config{}
var dao = MoviesDAO{}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/movies", findMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", findMovie).Methods("GET")
	router.HandleFunc("/movies", createMovies).Methods("POST")
	router.HandleFunc("/movies", updateMovies).Methods("PUT")
	router.HandleFunc("/movies", deleteMovies).Methods("DELETE")

	err := http.ListenAndServe(":3000", router)

	if err != nil {
		log.Fatal(err)
	}
}

func findMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := dao.FindAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movies)
}

func findMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.FindById(params["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID do filme Inválido!")
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func createMovies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Requisição inválida")
		return
	}

	movie.ID = bson.NewObjectId()
	if err := dao.Insert(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, movie)
}

func updateMovies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie

	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Requisição invalida")
		return
	}

	if err := dao.Update(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "sucess"})
}

func deleteMovies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie

	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := dao.Delete(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func init() {
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}
