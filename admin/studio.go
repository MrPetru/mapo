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

package admin

import (
	"github.com/maponet/utils/log"

	"labix.org/v2/mgo/bson"
	"net/http"
)

// NewStudio crea un nuovo studio
func httpNewStudio(out http.ResponseWriter, in *http.Request) {
	// create new studio
	log.Info("executing NewStudio function")

	errors := NewCoreErr()

	// creamo un nuovo contenitore di tipo studio
	studio := NewStudio()

	name := in.FormValue("name")
	err := studio.SetName(name)
	errors.append("name", err)

	currentuid := in.FormValue("currentuid")
	err = studio.AppendOwner(currentuid)
	errors.append("ownerid", err)

	id := in.FormValue("studioid")
	err = studio.SetId(id)
	errors.append("studioid", err)

	description := in.FormValue("description")
	err = studio.SetDescription(description)
	errors.append("description", err)

	if len(errors) > 0 {
		WriteJsonResult(out, errors, "error")
		return
	}

	err = studio.Save()
	if err != nil {
		errors.append("on store", err)
		WriteJsonResult(out, errors, "error")
		return
	}

	WriteJsonResult(out, studio, "ok")
}

// GetStudio restituisce al utente le informazioni di un solo progetto
func httpGetStudio(out http.ResponseWriter, in *http.Request) {

	errors := NewCoreErr()

	id := in.FormValue("sid")
	if len(id) == 0 {
		errors.append("id", "no studio id was provided")
		WriteJsonResult(out, errors, "error")
		return
	}

	currentuid := in.FormValue("currentuid")

	studio, err := StudioRestoreAll(bson.M{"owners": currentuid, "_id": id})

	if err != nil || len(studio) != 1 {
		errors.append("on restore", "error on studio restore from database")
		WriteJsonResult(out, errors, "error")
		return
	}

	WriteJsonResult(out, studio[0], "ok")
}

/*
UpdateStudio riceve i dati dal cliente e aggiorna quelli che sono già nella
database.

le situazioni:
    a. update
        assolutamente tutti i valori sono inviati dal cliente, che quelli non
        modificati.
    b. patch
        vengono inviati soltanto i valori che sono stati modificati. i campi
        non ricevuti dovranno essere ignorati. Ecco il link a il draft del
        path+json http://tools.ietf.org/html/draft-ietf-appsawg-json-patch-10
    NOTA: in entrambe le situazioni, per cancellare un valore si deve inviare
    un dato nullo per quella chiave.
*/
func httpUpdateStudio(out http.ResponseWriter, in *http.Request) {
	// patch o normal_update?

	// proviamo a implementare questa funzione come patch: i campi non ricevuti
	// verranno ignorati.

	sid := in.FormValue("sid")
	studio := NewStudio()
	studio.SetId(sid)
	err := studio.Restore()
	if err != nil {
		return
	}

	//description
	if _, ok := in.Form["description"]; ok {
		studio.SetDescription(in.FormValue("description"))
	}

	//name
	if _, ok := in.Form["name"]; ok {
		studio.SetName(in.FormValue("name"))
	}

	//owners
	if _, ok := in.Form["owner"]; ok {
		studio.AppendOwner(in.FormValue("owner"))
		// dobbiamo anche cancellare dei owners dalla lista?
	}

	studio.Update()

	WriteJsonResult(out, studio, "ok")
}

// GetStudioAll restituisce al cliente le informazioni di piu' progetti in una
// lista
func httpGetStudioAll(out http.ResponseWriter, in *http.Request) {
	// create new studio
	currentuid := in.FormValue("currentuid")

	studios, err := StudioRestoreAll(bson.M{"owners": currentuid})

	if err != nil {
		WriteJsonResult(out, err, "error")
	}
	WriteJsonResult(out, studios, "ok")
}
