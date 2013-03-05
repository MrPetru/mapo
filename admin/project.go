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
	"labix.org/v2/mgo/bson"
	"net/http"
)

/*
NewProject crea un nuovo progetto.
*/
func httpNewProject(out http.ResponseWriter, in *http.Request) {

	errors := NewCoreErr()

	project := NewProject()

	name := in.FormValue("name")
	err := project.SetName(name)
	errors.append("name", err)

	description := in.FormValue("description")
	err = project.SetDescription(description)
	errors.append("description", err)

	sidCookie, err := in.Cookie("sid")
	studioID := sidCookie.Value

	studio := NewStudio()
	err = studio.SetId(studioID)
	errors.append("studioid", err)

	err = studio.Restore()
	errors.append("on studio restore", err)

	if len(errors) > 0 {
		WriteJsonResult(out, errors, "error")
		return
	}

	id := Md5sum(studioID + name)
	err = project.SetId(id)
	errors.append("id", err)

	project.SetStudioId(studioID)

	_ = studio.AppendProject(id)
	err = studio.Update()
	errors.append("update studio", err)

	if len(errors) > 0 {
		WriteJsonResult(out, errors, "error")
		return
	}

	err = project.Save()
	if err != nil {
		errors.append("on save", err)
		WriteJsonResult(out, errors, "error")
		return
	}

	WriteJsonResult(out, project, "ok")
}

/*
GetProjectAll restituisce al cliente una lista di progetti per il studio attivo
nella sessione del utente.
*/
func httpGetProjectAll(out http.ResponseWriter, in *http.Request) {

	errors := NewCoreErr()

	sidCookie, err := in.Cookie("sid")
	if err != nil {
		errors.append("studio", "no active studio in current session")
		WriteJsonResult(out, errors, "error")
		return
	}
	studioID := sidCookie.Value

	filter := bson.M{"studioid": studioID}

	projectlist, err := ProjectRestorList(filter)

	if err != nil {
		WriteJsonResult(out, err, "error")
	}

	WriteJsonResult(out, projectlist, "ok")
}

/*
GetProject restituisce al utente le informazioni di un singolo progetto.
*/
func httpGetProject(out http.ResponseWriter, in *http.Request) {

	errors := NewCoreErr()

	id := in.FormValue("pid")
	if len(id) == 0 {
		errors.append("id", "no project id was provided")
		WriteJsonResult(out, errors, "error")
		return
	}

	sidCookie, err := in.Cookie("sid")
	if err != nil {
		errors.append("studioid", "no studio id was provided")
		WriteJsonResult(out, errors, "error")
		return
	}

	sid := sidCookie.Value

	studio := NewStudio()
	studio.SetId(sid)
	err = studio.Restore()
	if err != nil {
		return
	}

	if studio.Id != sid {
		return
	}

	project := NewProject()
	project.SetId(id)
	err = project.Restore()
	if err != nil {
		return
	}

	WriteJsonResult(out, project, "ok")
}

func httpAppendAddon(out http.ResponseWriter, in *http.Request) {
	errors := NewCoreErr()

	id := in.FormValue("pid")
	if len(id) == 0 {
		errors.append("id", "no project id was provided")
		WriteJsonResult(out, errors, "error")
		return
	}

	addonId := in.FormValue("addonid")
	entityName := in.FormValue("entityname")

	project := NewProject()
	project.SetId(id)
	err := project.Restore()
	if err != nil {
		errors.append("on restore", err)
		WriteJsonResult(out, errors, "error")
		return
	}

	if len(addonId) < 1 {
		errors.append("addonid", "id del addon troppo corto")
	}

	if len(entityName) < 1 {
		errors.append("entityname", "entityname troppo corto")
	}

	if len(errors) > 0 {
		WriteJsonResult(out, errors, "error")
		return
	}

	project.AddAddon(entityName, addonId)

	err = project.Update()
	if err != nil {
		errors.append("un restore", err)
		WriteJsonResult(out, errors, "error")
		return
	}

	WriteJsonResult(out, project, "ok")
}
