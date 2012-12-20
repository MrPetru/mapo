package core

import (
    "encoding/json"
    "net/http"
    "fmt"
)

// ExtractSingleValue è una funzione che aiuta a prendere un singolo valore
// dalla mappa di valori della forma.
func ExtractSingleValue(data map[string][]string, name string) string {
    v, ok := data[name]
    if !ok {
        return ""
    }
    
    if len(v) < 1 {
        return ""
    }
    
    if len(v) > 1 {
        return ""
    }
    
    return v[0]
}

// statusResult aiuta a formattare i dati in uscita
type statusResult struct {
    Status string `json:"status"`
    Data interface{} `json:"data"`
}

// WriteJsonResult è una scorciatoia per inviare il risultato verso il cliente
// in formato json.
func WriteJsonResult(out http.ResponseWriter, data interface{}, status string) {

    result := new(statusResult)

    // else send the result
    result.Status = status
    result.Data = data
    
    jsonResult, _ := json.Marshal(result)

    out.Header().Set("Content-Type","text/x-json")
    fmt.Fprint(out, string(jsonResult))
}

// coreErr è un contenitore per gli errori.
type coreErr map[string][]string

// NewCoreErr crea un nuovo oggetto di tipo coreErr
func NewCoreErr() coreErr{
    ce := make(coreErr, 0)
    return ce
}

// append aggiunge una nuovo elemento alla lista di errori per una chiave specifica.
func (ce *coreErr) append(key string, err error) {
    if err != nil {
        (*ce)[key] = append((*ce)[key], err.Error())
    }
}

func Authenticator(out http.ResponseWriter, in *http.Request) (http.ResponseWriter, *http.Request, bool) {

    fmt.Printf("authenticate for %v\n", in.URL.Path)
    
    return out, in, false
}

