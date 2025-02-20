package main

import (
    "github.com/githubtofu/chirpy/internal/database"
    "github.com/githubtofu/chirpy/internal/auth"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
)

func (cfg *apiConfig)loginHandler(w http.ResponseWriter, req *http.Request) {
    type JsonParams struct{
        Password string `json:"password"`
		Email string `json:"email"`
    }
    params := JsonParams{}
    err := json.NewDecoder(req.Body).Decode(&params)
    log.Printf("[usersHandler] Decoded Request:%v", params.Password)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
	user, err := cfg.db.LookUpUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	matched := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if matched != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
    type RespBody struct{
        ID  string `json:"id"`
        CreatedAt   string `json:"created_at"`
        UpdatedAt   string `json:"updated_at"`
        Email   string `json:"email"`
    }
    respondWithJSON(w, http.StatusOK, RespBody{
        ID: fmt.Sprintf("%v", user.ID),
        CreatedAt: fmt.Sprintf("%v", user.CreatedAt),
        UpdatedAt: fmt.Sprintf("%v", user.UpdatedAt),
        Email: string(user.Email),
    })
}

func (cfg *apiConfig)usersHandler(w http.ResponseWriter, req *http.Request) {
    type JsonParams struct{
        Password string `json:"password"`
		Email string `json:"email"`
    }
    params := JsonParams{}
    err := json.NewDecoder(req.Body).Decode(&params)
    log.Printf("[usersHandler] Decoded Request:%v", params.Password)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
	hashed, err := auth.HashPassword(params.Password)
    if err != nil {
        log.Printf("[UsersHandler]Error hashing : %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashed,
	})
    if err != nil {
        log.Printf("[usersHandler] failed to create user:%v", params.Email)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    type RespBody struct{
        ID  string `json:"id"`
        CreatedAt   string `json:"created_at"`
        UpdatedAt   string `json:"updated_at"`
        Email   string `json:"email"`
    }
    respondWithJSON(w, http.StatusCreated, RespBody{
        ID: fmt.Sprintf("%v", user.ID),
        CreatedAt: fmt.Sprintf("%v", user.CreatedAt),
        UpdatedAt: fmt.Sprintf("%v", user.UpdatedAt),
        Email: string(user.Email),
    })
}
