package main

import (
    "net/http"
    "sync/atomic"
    "strconv"
    "fmt"
)

type apiConfig struct {
    fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    cfg.fileserverHits.Add(1)
    return next
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("Hits: "))
    w.Write([]byte(strconv.Itoa(0)))
    cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) countHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte(fmt.Sprintf("Hit is: %d",cfg.fileserverHits.Load())))
}

func hHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}

func main(){
    sMux := http.NewServeMux()
    hServer := http.Server{
        Addr: ":8080",
        Handler: sMux,
    }
    cfg := apiConfig{}
    sMux.HandleFunc("/healthz", hHandler)
    sMux.HandleFunc("/metrics", cfg.countHandler)
    sMux.HandleFunc("/reset", cfg.resetHandler)
    sMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
    hServer.ListenAndServe()
}
