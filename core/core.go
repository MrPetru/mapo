package core

import (
    "encoding/json"
    "net/http"
    "fmt"
)

// ExtractSingleValue è una funzione che aiuta a prendere un singolo valore
// dalla mappa di valori della forma (in.Form per esempio).
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

// statusResult aiuta a formattare i dati inviati verso il cliente
type statusResult struct {
    Status string `json:"status"`
    Data interface{} `json:"data"`
}

// WriteJsonResult è una scorciatoia per inviare il risultato verso il cliente
// in formato json.
func WriteJsonResult(out http.ResponseWriter, data interface{}, status string) {

    result := new(statusResult)

    result.Status = status
    result.Data = data

    jsonResult, _ := json.Marshal(result)

    out.Header().Set("Content-Type","application/json")
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
func (ce *coreErr) append(key string, err interface{}) {
    if err == nil {
        return
    }

    if e, ok := err.(error); ok {
        if e != nil {
            (*ce)[key] = append((*ce)[key], e.Error())
        }
    } else {
        (*ce)[key] = append((*ce)[key], err.(string))
    }
}

