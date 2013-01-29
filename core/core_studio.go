package core

import (
    "mapo/log"
    "mapo/objectspace"

    "net/http"
    "labix.org/v2/mgo/bson"
)

// NewStudio crea un nuovo studio
func NewStudio(out http.ResponseWriter, in *http.Request) {
    // create new studio
    log.Msg("executing NewStudio function")

    errors := NewCoreErr()

    // creamo un nuovo contenitore di tipo studio
    studio := objectspace.NewStudio()

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
func GetStudio(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    id := in.FormValue("sid")
    if len(id) == 0 {
        errors.append("id", "no studio id was provided")
        WriteJsonResult(out, errors, "error")
        return
    }

    currentuid := in.FormValue("currentuid")

    studio, err := objectspace.StudioRestoreAll(bson.M{"owners":currentuid, "_id":id})

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
func UpdateStudio(out http.ResponseWriter, in *http.Request) {
    // patch o normal_update?

    // proviamo a implementare questa funzione come patch: i campi non ricevuti
    // verranno ignorati.

    sid := in.FormValue("sid")
    studio := objectspace.NewStudio()
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
func GetStudioAll(out http.ResponseWriter, in *http.Request) {
    // create new studio
    currentuid := in.FormValue("currentuid")

    studios, err := objectspace.StudioRestoreAll(bson.M{"owners":currentuid})

    if err != nil {
        WriteJsonResult(out, err, "error")
    }
    WriteJsonResult(out, studios, "ok")
}
