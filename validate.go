package main

import (
    "net/http"
    "encoding/json"
    "log"
)

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
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    if len(params.Body) > max_chirp_length {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
    }else {
        type RespBody struct{
            V string `json:"cleaned_body"`
        }
        respondWithJSON(w, http.StatusOK, RespBody{
            V: replacePatterns( string(params.Body),
                []string{ "kerfuffle", "sharbert", "fornax" },
                "****") })
//func replacePatterns(base string, patterns []string, profane_mask string) string{
    }
    //w.WriteHeader(200)
}
