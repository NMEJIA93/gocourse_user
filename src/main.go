package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/NMEJIA93/gocourse_user/pkg/bootstrap"
	"github.com/NMEJIA93/gocourse_user/src/user"
)

func main() {

	router := mux.NewRouter()

	_ = godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found or couldn't be loaded")
	}

	l := bootstrap.InitLogger()

	db, err := bootstrap.BDConnection()
	if err != nil {
		log.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	userRepo := user.NewRepository(l, db)
	userService := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userService, user.Config{LimitPageDef: pagLimDef})

	router.HandleFunc("/user/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/user", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/user", userEnd.Create).Methods("POST")
	router.HandleFunc("/user", userEnd.Update).Methods("PUT")
	router.HandleFunc("/user/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/user/{id}", userEnd.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	address := fmt.Sprintf("localhost:%s", port)

	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	l.Println("starting server on", srv.Addr)

	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err1)
	}
}
