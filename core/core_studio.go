package core

import (
    "mapo/log"
    "mapo/objectspace"

    "net/http"
    "labix.org/v2/mgo/bson"
    "strings"

//    // utilizo di questo paccheto e soltanto temporaniamente
//    // per creare un id che poi avrà una funzione specifica
//    "labix.org/v2/mgo/bson"
)

func NewStudio(out http.ResponseWriter, in *http.Request) {
    // create new studio
    log.Msg("executing NewStudio function")

    //in.ParseForm()
    errors := NewCoreErr()

    // creamo un nuovo contenitore di tipo studio
    studio := objectspace.NewStudio()

    name := ExtractSingleValue(in.Form, "name")
    err := studio.SetName(name)
    errors.append("name", err)

    currentuid := ExtractSingleValue(in.Form, "currentuid")
    err = studio.AppendOwner(currentuid)
    errors.append("ownerid", err)

    id := name
    err = studio.SetId(id)
    errors.append("id", err)

    // update user
    user := objectspace.NewUser()
    err = user.SetId(currentuid)
    errors.append("userid", err)

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

    filter := bson.M{"_id":currentuid}
    err = user.Restore(filter)
    if err != nil {
        errors.append("on restore", err)
        WriteJsonResult(out, errors, "error")
        return
    }

    err = user.AppendStudioId(studio.GetId())
    errors.append("studioid", err)

    err = user.Update()
    errors.append("on user update", err)

    if len(errors) > 0 {
        WriteJsonResult(out, errors, "error")
        return
    }

    WriteJsonResult(out, studio, "ok")
}

func GetStudio(out http.ResponseWriter, in *http.Request) {
    // create new studio

    errors := NewCoreErr()

    id := strings.Split(in.URL.Path[1:], "/")[2]
    if len(id) == 0 {
        errors.append("id", "no studio id was provided")
        WriteJsonResult(out, errors, "error")
        return
    }

    currentuid := ExtractSingleValue(in.Form, "currentuid")

    user := objectspace.NewUser()

    filter := bson.M{"_id":currentuid}
    user.Restore(filter)

    studiosid := user.GetStudiosId()
    for _, sid := range(studiosid) {
        if sid == id {
            studio := objectspace.NewStudio()
            filter := bson.M{"_id":id}
            err := studio.Restore(filter)
            if err != nil {
                errors.append("on restore", err)
                WriteJsonResult(out, errors, "error")
                return
            }
            WriteJsonResult(out, studio, "ok")
            return
        }
    }
    errors.append("on restore", "no studio was found")
    WriteJsonResult(out, errors, "error")
}

func GetStudioAll(out http.ResponseWriter, in *http.Request) {
    // create new studio
    currentuid := ExtractSingleValue(in.Form, "currentuid")

    user := objectspace.NewUser()

    filter := bson.M{"_id":currentuid}
    user.Restore(filter)

    studiosid := user.GetStudiosId()

    WriteJsonResult(out, studiosid, "ok")
}

//func UpdateStudio(inValues values) interface{} {
//    // create new studio
//    
//    return nil
//}
