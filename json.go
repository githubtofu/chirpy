package main


func respondWithError(w http.ResponseWriter, code int, msg string){
    type respBody struct{
        E string `json:"error"`
    }
    body := respBody{
        E: msg,
    }
    data, err := json.Marshal(&body)
    if err != nil {
        log.Printf("Error marshalling response body, %s", err)
        w.WriteHeader(500)
        return 
    }
    w.WriteHeader(400)
    w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
    
}

