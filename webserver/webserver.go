package webserver

import (
	"context"
	"github.com/gonamore/fxbd/config/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type WebServer struct {
	applicationConfig *models.ApplicationConfig
}

func NewWebServer(applicationConfig *models.ApplicationConfig) *WebServer {
	return &WebServer{applicationConfig: applicationConfig}
}

func (rcv *WebServer) Start() {
	//http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
	//	http.ServeFile(res, req, path.Join(applicationConfig.StatsDir, "index.html"))
	//})
	//log.Fatal(http.ListenAndServe(":" + strconv.Itoa(rcv.applicationConfig.Port), nil))

	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("webserver/assets/"))))
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(rcv.applicationConfig.StatsDir))))

	server := &http.Server{Addr: ":" + strconv.Itoa(rcv.applicationConfig.Port), Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
