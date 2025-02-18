package main

import (
    "net/http"
    "sync/atomic"
    "strconv"
    "fmt"
    "encoding/json"
    "log"
)

type apiConfig struct {
    fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    //HandlerFunc(f)  -> Handler that calls f
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("Hits: "))
    w.Write([]byte(strconv.Itoa(0)))
    cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) countHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,cfg.fileserverHits.Load())))
}

func validateHandler(w http.ResponseWriter, req *http.Request) {
    const max_chirp_length = 140
    type JsonParams struct{
        Body string `json:"body"`
    }
    params := JsonParams{}
    err := json.NewDecoder(req.Body).Decode(&params)
    log.Printf("[validateHandler] Decoded Request:%v", params.Body)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(500)
        return 
    }
    if len(params.Body) > max_chirp_length {
        type respBody struct{
            E string `json:"error"`
        }
        body := respBody{
            E: "Chirp is too long",
        }
        data, err := json.Marshal(&body)
        if err != nil {
            log.Printf("Error marshalling response body, %s", err)
            w.WriteHeader(500)
            return 
        }
        w.WriteHeader(400)
        w.Write(data)
    }else {
        type RespBody struct{
            V bool `json:"valid"`
        }
        body := RespBody{
            V: true,
        }
        data, err := json.Marshal(&body)
        if err != nil {
            log.Printf("Error marshalling response body, %s", err)
            w.WriteHeader(500)
            return 
        }
        w.Write(data)
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    //w.WriteHeader(200)
}

func hHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}

func main(){
    mux := http.NewServeMux()
    hServer := http.Server{
        Addr: ":8080",
        Handler: mux,
    }
    cfg := apiConfig{}
    mux.HandleFunc("POST /api/validate_chirp", validateHandler)
    mux.HandleFunc("GET /api/healthz", hHandler)
    mux.HandleFunc("GET /admin/metrics", cfg.countHandler)
    mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
    mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    hServer.ListenAndServe()
}
