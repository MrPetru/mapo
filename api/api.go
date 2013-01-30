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
Package api implements the API framework used by addons.
*/
package api

import (
    "mapo/log"
    "mapo/admin"
    "mapo/objectspace"

    "net/http"
    "strings"
    "labix.org/v2/mgo/bson"
)

// apiData e' il contenitore dei dati che vengono inviati verso la funzione
// che processa la richiesta.
type apiData struct {
    Method string
    ProjectId string
    Resource string
    ResourceId string
    ResourceFunction string

    StudioId string
    ExtraData map[string][]string
}

// NewApiData crea un nuovo oggetto apiData
func NewApiData() *apiData{
    a := new(apiData)

    return a
}

// ApiRouter identifica e esegue la funzione del addon che deve essere eseguita
// per la risorsa richiesta.
func ApiRouter(data *apiData) (*apiData, error) {

    var err error

    // crea la lista dei addon per il progetto corrente
    {
        studios, err := objectspace.StudioRestoreAll(bson.M{"_id":data.StudioId,"projects":data.ProjectId})
        if err != nil || len(studios) != 1 {
            return nil, err
        }
    }

    project := objectspace.NewProject()
    project.SetId(data.ProjectId)
    err = project.Restore()
    if err != nil {
        return nil, err
    }

    addons := project.GetAddonList()
    addons = addons

    // identifica la funzione da eseguire

    // avvia la funzione

    // ritorna il risultato al cliente

    return data, nil
}

// e' la prima funzione chiamata da mapo che avvia il router delle api
// un convertitore da una richiesta http in una forma piu' generale per mapo
// potrebbero essere vari questi wrapper, per esempio per una richiesta ftp o via email.
func HttpWrapper(out http.ResponseWriter, in *http.Request) {

    //pathPattern := "method:/api/{projectId}/{resource}/{resourceId}/{function}"

    urlValues := make([]string, 0)
    {
        values := strings.Split(in.URL.Path, "/")
        for i :=0; i<6; i++ {
            if i <= (len(values) - 1) {
                urlValues = append(urlValues, values[i])
            } else {
                urlValues = append(urlValues, "")
            }
        }
    }

    data := NewApiData()
    data.Method = in.Method
    data.ProjectId = urlValues[2]
    data.Resource = urlValues[3]
    data.ResourceId = urlValues[4]
    data.ResourceFunction = urlValues[5]

    {
        eData := make(map[string][]string)
        for i,v := range(in.Form) {
            eData[i] = v
        }
        if c, err := in.Cookie("sid"); err == nil {
            data.StudioId = c.Value
        }
        if c, err := in.Cookie("uid"); err == nil {
            v := c.Value
            if v != in.FormValue("currentuid") {
                admin.WriteJsonResult(out, "authentication don't match", "error")
                return
            }
        }
        if c, err := in.Cookie("pid"); err != nil {
            v := c.Value
            if v != data.ProjectId {
                admin.WriteJsonResult(out, "project don't match", "error")
                return
            }
        }

        data.ExtraData = eData
    }

    result, err := ApiRouter(data)
    if err != nil {
        admin.WriteJsonResult(out, err.Error(), "error")
        return
    }

    log.Debug("api data = %v", data)
    log.Debug("result = %v", result)

    admin.WriteJsonResult(out, result, "ok")
}
