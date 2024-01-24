package main

import (
	"encoding/json"
	"handlar_tes/model"
	"handlar_tes/storge"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {

	mux := mux.NewRouter()
    mux.HandleFunc("/user/create", CreateUser).Methods("POST")
	mux.HandleFunc("/user/update", UpdateUser).Methods("PUT")
	mux.HandleFunc("/user/delete", DeleteUser).Methods("DELETE")
	mux.HandleFunc("/user/get", GetUser).Methods("GET")
	mux.HandleFunc("/user/all", GetAllUsers).Methods("GET")
	log.Println("Server is running...")
	if err := http.ListenAndServe("localhost:7777", mux); err != nil {
		log.Println("Error sever is running!")
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user *model.User

	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.NewString()
	user.ID = id

	respUser, err := storge.CreateUser(user)
	if err != nil {
		log.Println("error while creating body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println("Error while converting page")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error while converting page")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := storge.GetAll(intPage, intLimit)
	if err != nil {
		log.Println("Error while getting all users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(users)
	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	user, err := storge.Get(id)

	if err != nil {
		log.Println("Error while getting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(user)

	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := storge.DeleteUser(id)

	if err != nil {
		log.Println("Error while deleting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u model.User

	err = json.Unmarshal(bodyByte, &u)

	user, err := storge.UptadeUser(id, &u)

	if err != nil {
		log.Println("Error not update user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(user)

	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}
