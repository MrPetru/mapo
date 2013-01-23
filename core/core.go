package core

import (
    "encoding/json"
    "net/http"
    "fmt"
    "gconf/conf"
)

// statusResult aiuta a formattare i dati inviati verso il cliente
type statusResult struct {
    Status string `json:"status"`
    Data interface{} `json:"data"`
}

// WriteJsonResult è una scorciatoia per inviare il risultato verso il cliente
// in formato json.
// TODO: in caso di errore che codice dobbiamo ritornare? 412? 424?
func WriteJsonResult(out http.ResponseWriter, data interface{}, status string) {

    result := new(statusResult)

    result.Status = status
    result.Data = data

    jsonResult, _ := json.Marshal(result)

    out.Header().Set("Content-Type","application/json;charset=UTF-8")
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

/*
GlobalConfiguration, il oggetto globale per l'accesso ai dati contenuti nel
file di configurazione.
*/
var GlobalConfiguration *conf.ConfigFile

/*
ReadConfiguration, attiva il GlobalConfiguration.
*/
func ReadConfiguration(filepath string) error {

    c, err := conf.ReadConfigFile(filepath)
    if err == nil {
        GlobalConfiguration = c
    }

    return err
}
