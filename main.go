package main

import _ "github.com/lib/pq"

import (
	"log"
    "net/http"
    "sync/atomic"
    "strconv"
    "fmt"
    "strings"
    "github.com/githubtofu/chirpy/internal/database"
    "github.com/joho/godotenv"
    "os"
    "database/sql"
)

type apiConfig struct {
    fileserverHits atomic.Int32
    dbURL string
    db *database.Queries
    PLATFORM string
	SECRET string
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
    if cfg.PLATFORM != "dev" {
        respondWithError(w, http.StatusForbidden, "Access denied")
        return
    }
    cfg.db.DeleteAllUsers(req.Context())
    w.WriteHeader(http.StatusOK)
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

func hHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}

func replacePatterns(base string, patterns []string, profane_mask string) string{
    split_base := strings.Split(base, " ")
    masked := []string{}
    for _, a_word := range(split_base) {
        is_profane := false
        for _, a_pattern := range(patterns) {
            if strings.ToLower(a_word) == strings.ToLower(a_pattern) {
                masked = append(masked, profane_mask)
                is_profane = true
                break
            }
        }
        if !is_profane {
            masked = append(masked, a_word)
        }
    }
    return strings.Join(masked, " ")
}

func main(){
    godotenv.Load()
    dbURL := os.Getenv("DB_URL")
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return
    }
    dbQueries := database.New(db)
    fmt.Println(dbQueries)
    mux := http.NewServeMux()
    hServer := http.Server{
        Addr: ":8080",
        Handler: mux,
    }
	s := os.Getenv("SECRET")
	log.Println("[main] SECRET:", s)
    cfg := apiConfig{
        dbURL: dbURL,
        db : dbQueries,
        PLATFORM : os.Getenv("PLATFORM"),
		SECRET : os.Getenv("SECRET"),
    }
    mux.HandleFunc("POST /api/validate_chirp", validateHandler)
    mux.HandleFunc("POST /api/users", cfg.usersHandler)
    mux.HandleFunc("POST /api/chirps", cfg.chirpsHandler)
    mux.HandleFunc("POST /api/login", cfg.loginHandler)
    mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
    mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getAChirpHandler)
    mux.HandleFunc("GET /api/healthz", hHandler)
    mux.HandleFunc("GET /admin/metrics", cfg.countHandler)
    mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
    mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    hServer.ListenAndServe()
}
