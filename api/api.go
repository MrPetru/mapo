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
    "github.com/maponet/utils/log"
    "mapo/admin"
    "mapo/objectspace"

    "net/http"
    "strings"
	"errors"
    "labix.org/v2/mgo/bson"

	"mapo/addons/repo/go/scene"
	"mapo/addons/repo/go/shot"
	"mapo/addons/repo/go/shotpatch"
)

// apiData e' il contenitore dei dati che vengono inviati verso la funzione
// che processa la richiesta.
type apiData struct {
    Method string
    ProjectId string
    ResourceType string
    ResourceId string
    ResourceFunction string

    StudioId string
    ExtraData map[string][]string
}

func (data *apiData) GetValue(name string) string {
	if name == "id" {
		return data.ResourceId
	}

	dataElement, ok := data.ExtraData[name]
	if !ok {
		//panic("key not found in api data")
		return ""
	}
	return dataElement[0]
}

func RegisterAddons() {
	newAddonContainer()
	scene.Register(&Addons)
	shot.Register(&Addons)
	shotpatch.Register(&Addons)
}

// NewApiData crea un nuovo oggetto apiData
func NewApiData() *apiData{
    a := new(apiData)

    return a
}

// ApiRouter identifica e esegue la funzione del addon che deve essere eseguita
// per la risorsa richiesta.
func ApiRouter(data *apiData) (interface{}, error) {//(*apiData, error) {

    var err error

    // verifica se il progetto e lo studio sono collegati
    {
        studios, err := objectspace.StudioRestoreAll(bson.M{"_id":data.StudioId,"projects":data.ProjectId})
        if err != nil {
			return nil, err
		}
		if len(studios) != 1 {
			log.Error("incorect query to database")
            return nil, errors.New("cant't find asociated projects")
        }
    }

    project := objectspace.NewProject()
    project.SetId(data.ProjectId)
    err = project.Restore()
    if err != nil {
        return nil, err
    }

	// creare il path della funzione da eseguire
	fPath := ""

	fPath = fPath + data.ResourceType
	if len(data.ResourceType) > 0 {
		fPath = "/"
	} else {
		// operation not o entity
		// run defalut project function
		return nil, errors.New("project handler not found")
	}

	fPath = fPath + data.ResourceId
	if fPath[len(fPath)-1] != '/' {
		fPath = fPath + "/"
	}

	fPath = fPath + data.ResourceFunction
	if fPath[len(fPath)-1] != '/' {
		fPath = fPath + "/"
	}

	log.Debug("requested function is %v", fPath)

	// get entity type
    addonsId := project.GetAddonList(data.ResourceType)
	if len(addonsId) < 1 {
		// run default action
		// return result, err
		return nil, errors.New("no active addons was found")
	}

	// construct Entities
	entitiesList := NewEntitiesList()

	// create a ordered addons dependency list
	orderedAddons := orderByDependency(addonsId, Addons)
	for _, a := range(orderedAddons) {
		constructors := Addons[a].Constructors
		for _, c := range(constructors) {
			c(entitiesList)
		}
	}
	for _, e := range(*entitiesList) {
		e.projectId = data.ProjectId
	}

	resultCompEntity, err := entitiesList.Run(data.ResourceType, data.Method, fPath, data)
	if err == nil {
		compE, _ := resultCompEntity.(*composedEntity)
		return compE.ToMap(), nil
	}
	return nil, err
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
    data.ResourceType = urlValues[3]
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
        if c, err := in.Cookie("pid"); err == nil {
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

    log.Debug("err = %v", err)
    log.Debug("api data = %v", data)
    log.Debug("result = %v", result)

    admin.WriteJsonResult(out, result, "ok")
}
