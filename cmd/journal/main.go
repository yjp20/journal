package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var (
	addr     = flag.String("addr", ":4000", "HTTP network address")
	password = flag.String("password", "password", "Login password")
	db       = flag.String("db", "dbname=journal sslmode=disable", "DsN string")
	cors     = flag.String("cors", "http://localhost:3000", "CORS origin")
)

func main() {
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("postgres", *db)
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	app := &App{
		Addr:     *addr,
		DB:       db,
		Password: *password,
		Models: Models{
			Todo:       &TodoModel{db},
			Media:      &MediaModel{db},
			FeedSource: &FeedSourceModel{db},
			FeedItem:   &FeedItemModel{db},
		},
		Origins:  strings.Split(*cors, ","),
		errorLog: errorLog,
	}

	router := httprouter.New()

	router.GET("/api/todo", app.listTodo)
	router.POST("/api/todo", app.createTodo)
	router.PUT("/api/todo/:id", app.updateTodo)
	router.DELETE("/api/todo/:id", app.deleteTodo)

	router.POST("/api/media/link", app.linkMedia)
	router.GET("/api/media", app.listMedia)
	router.POST("/api/media", app.createMedia)
	router.PUT("/api/media/:id", app.updateMedia)
	router.DELETE("/api/media/:id", app.deleteMedia)

	router.GET("/api/feedsource", app.listFeedSource)
	router.POST("/api/feedsource", app.subscribeFeedSource)
	router.DELETE("/api/feedsource/:id", app.unsubscribeFeedSource)

	router.GET("/api/feed", app.getFeed)
	router.POST("/api/feed/collect", app.collectFeed)

	router.POST("/api/token", app.confirmToken)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { app.notFoundResponse(w) })

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.enableCORS(router),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func (a *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")

		for i := range a.Origins {
			if a.Origins[i] == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				if r.Method == http.MethodOptions {
					w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
