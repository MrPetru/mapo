package core

import (
    "mapo/log"
    "mapo/objectspace"

    "net/http"
    "labix.org/v2/mgo/bson"
    "strings"
)

func NewStudio(out http.ResponseWriter, in *http.Request) {
    // create new studio
    log.Msg("executing NewStudio function")

    //in.ParseForm()
    errors := NewCoreErr()

    // creamo un nuovo contenitore di tipo studio
    studio := objectspace.NewStudio()

    name := in.FormValue("name")
    err := studio.SetName(name)
    errors.append("name", err)

    currentuid := ExtractSingleValue(in.Form, "currentuid")
    err = studio.AppendOwner(currentuid)
    errors.append("ownerid", err)

    id := ExtractSingleValue(in.Form, "studioid")
    err = studio.SetId(id)
    errors.append("studioid", err)

    description := ExtractSingleValue(in.Form, "description")
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

func GetStudio(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    id := strings.Split(in.URL.Path[1:], "/")[2]
    if len(id) == 0 {
        errors.append("id", "no studio id was provided")
        WriteJsonResult(out, errors, "error")
        return
    }

    currentuid := ExtractSingleValue(in.Form, "currentuid")

    studio, err := objectspace.StudioRestoreAll(bson.M{"owners":currentuid, "_id":id})

    if err != nil || len(studio) != 1 {
        errors.append("on restore", "error on studio restore from database")
        WriteJsonResult(out, errors, "error")
        return
    }

    WriteJsonResult(out, studio[0], "ok")
}

func GetStudioAll(out http.ResponseWriter, in *http.Request) {
    // create new studio
    currentuid := ExtractSingleValue(in.Form, "currentuid")

    studios, err := objectspace.StudioRestoreAll(bson.M{"owners":currentuid})

    if err != nil {
        WriteJsonResult(out, err, "error")
    }
    WriteJsonResult(out, studios, "ok")
}
