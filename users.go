package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
)
func (cfg *apiConfig)usersHandler(w http.ResponseWriter, req *http.Request) {
    type JsonParams struct{
        Body string `json:"email"`
    }
    params := JsonParams{}
    err := json.NewDecoder(req.Body).Decode(&params)
    log.Printf("[usersHandler] Decoded Request:%v", params.Body)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    user, err := cfg.db.CreateUser(req.Context(), params.Body)
    if err != nil {
        log.Printf("[usersHandler] failed to create user:%v", params.Body)
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

