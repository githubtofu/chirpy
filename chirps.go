package main

import (
    "net/http"
    "github.com/githubtofu/chirpy/internal/auth"
    "encoding/json"
    "log"
    "github.com/githubtofu/chirpy/internal/database"
    "github.com/google/uuid"
    "fmt"
)
func (cfg *apiConfig)getAChirpHandler(w http.ResponseWriter, req *http.Request) {
	pValue := fmt.Sprintf("%s", req.PathValue("chirpID"))
	log.Printf("[AChirp Handler] Got Path from request:%s", pValue)
	chirp_id, err := uuid.Parse(pValue)
    if err != nil {
        log.Printf("Error getting a chirp id (UUID): %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    chirp, err := cfg.db.GetAChirp(req.Context(), chirp_id)
    if err != nil {
        log.Printf("Error getting a chirp: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    type RespItem struct{
        UserID string `json:"user_id"`
        ID string `json:"id"`
        CreatedAt string `json:"created_at"`
        UpdatedAt string `json:"updated_at"`
        Body string `json:"body"`
    }
	response_json := RespItem{
		UserID: fmt.Sprintf("%s", chirp.UserID),
		ID:fmt.Sprintf("%s",  chirp.ID),
		CreatedAt: fmt.Sprintf("%s", chirp.CreatedAt),
		UpdatedAt: fmt.Sprintf("%s", chirp.UpdatedAt),
		Body: chirp.Body,
	}
    respondWithJSON(w, http.StatusOK, response_json)
}

func (cfg *apiConfig)getChirpsHandler(w http.ResponseWriter, req *http.Request) {
    chirps, err := cfg.db.GetChirps(req.Context())
    if err != nil {
        log.Printf("Error getting chirps: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    type RespItem struct{
        UserID string `json:"user_id"`
        ID string `json:"id"`
        CreatedAt string `json:"created_at"`
        UpdatedAt string `json:"updated_at"`
        Body string `json:"body"`
    }
    response_json := []RespItem{}
    for _, a_chirp := range(chirps) {
        response_json = append(response_json, RespItem{
            UserID: fmt.Sprintf("%s", a_chirp.UserID),
            ID:fmt.Sprintf("%s",  a_chirp.ID),
            CreatedAt: fmt.Sprintf("%s", a_chirp.CreatedAt),
            UpdatedAt: fmt.Sprintf("%s", a_chirp.UpdatedAt),
            Body: a_chirp.Body,
        })
    }
    respondWithJSON(w, http.StatusOK, response_json)
}

func (cfg *apiConfig)chirpsHandler(w http.ResponseWriter, req *http.Request) {
	btoken, err := auth.GetBearerToken(req.Header)
    if err != nil {
        log.Printf("Error getting token: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
	log.Println("[chirpsHandler] btoken:", btoken)
    const max_chirp_length = 140
    type JsonParams struct{
        Body string `json:"body"`
        JWT string `json:"jwt"`
    }
    params := JsonParams{}
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    log.Println("[chirpsHandler] Decoded Request")
    log.Printf("Body:%v", params.Body)
    log.Printf("Params all:%v", params)
	uid, err := auth.ValidateJWT(btoken, cfg.SECRET)
    if err != nil {
        log.Printf("Error validating: %s", err)
        w.WriteHeader(http.StatusUnauthorized)
        return 
    }
    if len(params.Body) > max_chirp_length {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
    }else {
        //user_id_uuid, err := uuid.Parse(uid)
        if err != nil {
            log.Printf("Error decoding user id: %s", err)
            w.WriteHeader(http.StatusInternalServerError)
            return 
        }
        chirp, err := cfg.db.CreateChirp(req.Context(), 
            database.CreateChirpParams{ Body: params.Body,
                UserID: uid,
            })
        if err != nil {
            log.Printf("[chirpsHandler] failed to create chirp:%v", params.Body)
            w.WriteHeader(http.StatusInternalServerError)
            return 
        }
        type RespBody struct{
            UserID string `json:"user_id"`
            ID string `json:"id"`
            CreatedAt string `json:"created_at"`
            UpdatedAt string `json:"updated_at"`
            Body string `json:"body"`
        }
        respondWithJSON(w, http.StatusCreated, RespBody{
            ID: fmt.Sprintf("%v", chirp.ID),
            CreatedAt: "tesetd date",//fmt.Sprintf("%v", chirp.CreatedAt),
            UpdatedAt: fmt.Sprintf("%v", chirp.UpdatedAt),
            Body: string(chirp.Body),
            UserID: fmt.Sprintf("%v", chirp.UserID),
        })
    }
}

