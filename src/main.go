package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/NMEJIA93/gocourse_user/pkg/bootstrap"
	"github.com/NMEJIA93/gocourse_user/pkg/handler"
	"github.com/NMEJIA93/gocourse_user/src/user"
)

func main() {

	//router := mux.NewRouter()

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

	ctx := context.Background()

	userRepo := user.NewRepository(l, db)
	userService := user.NewService(l, userRepo)
	//userEnd := user.MakeEndpoints(userService)

	handleUser := handler.NewUserHTTPServer(ctx, user.MakeEndpoints(userService, user.Config{LimitPageDef: pagLimDef}))

	//router.HandleFunc("/user/{id}", userEnd.Get).Methods("GET")
	//router.HandleFunc("/user", userEnd.GetAll).Methods("GET")
	//router.HandleFunc("/user", userEnd.Create).Methods("POST")
	//router.HandleFunc("/user", userEnd.Update).Methods("PUT")
	//router.HandleFunc("/user/{id}", userEnd.Update).Methods("PATCH")
	//router.HandleFunc("/user/{id}", userEnd.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	address := fmt.Sprintf("localhost:%s", port)

	srv := &http.Server{
		Handler:      accessControl(handleUser),
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errChannel := make(chan error)
	go func() {
		l.Println(fmt.Sprintf("listening on port %s", port))
		errChannel <- srv.ListenAndServe()
	}()

	l.Println("starting server on", srv.Addr)

	err1 := <-errChannel
	if err1 != nil {
		log.Fatal(err1)
	}
}

func accessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
