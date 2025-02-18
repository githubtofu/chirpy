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
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
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
