package core

import (
    "mapo/objectspace"

    "net/http"
    "labix.org/v2/mgo/bson"
)

/*
NewProject crea un nuovo progetto.
*/
func NewProject(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    project := objectspace.NewProject()

    name := in.FormValue("name")
    err := project.SetName(name)
    errors.append("name", err)

    description := in.FormValue("description")
    err = project.SetDescription(description)
    errors.append("description", err)

    sidCookie, err := in.Cookie("sid")
    studioID := sidCookie.Value

    studio := objectspace.NewStudio()
    err = studio.SetId(studioID)
    errors.append("studioid", err)

    err = studio.Restore()
    errors.append("on studio restore", err)

    if len(errors) > 0 {
        WriteJsonResult(out, errors, "error")
        return
    }

    id := objectspace.Md5sum(studioID + name)
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
func GetProjectAll(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    sidCookie, err := in.Cookie("sid")
    if err != nil {
        errors.append("studio", "no active studio in current session")
        WriteJsonResult(out, errors, "error")
        return
    }
    studioID := sidCookie.Value

    filter := bson.M{"studioid":studioID}

    projectlist, err := objectspace.ProjectRestorList(filter)

    if err != nil {
        WriteJsonResult(out, err, "error")
    }

    WriteJsonResult(out, projectlist, "ok")
}

/*
GetProject restituisce al utente le informazioni di un singolo progetto.
*/
func GetProject(out http.ResponseWriter, in *http.Request) {

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

    studio := objectspace.NewStudio()
    studio.SetId(sid)
    err = studio.Restore()
    if err != nil {
        return
    }

    if studio.Id != sid {
        return
    }

    project := objectspace.NewProject()
    project.SetId(id)
    err = project.Restore()
    if err != nil {
        return
    }

    WriteJsonResult(out, project, "ok")
}
