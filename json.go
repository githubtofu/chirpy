package main

import (
    "encoding/json"
    "net/http"
    "log"
)

func respondWithError(w http.ResponseWriter, code int, msg string){
    type respBody struct{
        E string `json:"error"`
    }
    respondWithJSON(w, code, respBody{E: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling response body, %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    w.WriteHeader(code)
    w.Write(data)
}

