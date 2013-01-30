/*
Copyright 2013 Petru Ciobanu, Francesco Paglia, Lorenzo Pierfederici

This file is part of Mapo.

Mapo is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Mapo is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Mapo.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
Package admin implements the API for Mapo's administration components.
*/
package admin

import (
    "gconf/conf"
    "encoding/json"
    "net/http"
    "fmt"
)

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
