package main

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/rcrowley/go-metrics"
	"github.com/akhenakh/statgo"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	s := statgo.NewStat()


	hitCounter := metrics.NewMeter()
	metrics.Register("hits", hitCounter)

	adminHandler := http.NewServeMux()
	adminHandler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			hitCounter.Mark(1)
		})
	adminHandler.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			hitCounter.Mark(1)
			fmt.Fprint(w, "{\"system\":")
			bytes, _ := json.Marshal(s.HostInfos())
			w.Write(bytes)
			fmt.Fprint(w, ",\"memory\":")
			bytes, _ = json.Marshal(s.MemStats())
			w.Write(bytes)
			fmt.Fprint(w, ",\"network\":")
			bytes, _ = json.Marshal(s.NetIOStats())
			w.Write(bytes)
			fmt.Fprint(w, ",\"metrics\":")
			metrics.WriteJSONOnce(metrics.DefaultRegistry, w)
			fmt.Fprint(w, "}")
		})
	admin := &http.Server{
		Addr: ":" + port,
		Handler: adminHandler,
	}
	log.Fatal(admin.ListenAndServe())
}