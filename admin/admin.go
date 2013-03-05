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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
)

// statusResult aiuta a formattare i dati inviati verso il cliente
type statusResult struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// WriteJsonResult è una scorciatoia per inviare il risultato verso il cliente
// in formato json.
// TODO: in caso di errore che codice dobbiamo ritornare? 412? 424?
func WriteJsonResult(out http.ResponseWriter, data interface{}, status string) {

	result := new(statusResult)

	result.Status = status
	result.Data = data

	jsonResult, _ := json.Marshal(result)

	out.Header().Set("Content-Type", "application/json;charset=UTF-8")
	fmt.Fprint(out, string(jsonResult))
}

// coreErr è un contenitore per gli errori.
type coreErr map[string][]string

// NewCoreErr crea un nuovo oggetto di tipo coreErr
func NewCoreErr() coreErr {
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

// crea la soma md5 di una stringa
func Md5sum(value string) string {
	sum := md5.New()
	sum.Write([]byte(value))

	result := fmt.Sprintf("%x", sum.Sum(nil))

	return result
}
