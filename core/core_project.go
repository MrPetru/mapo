package core

import (
    "net/http"
    "mapo/objectspace"
    "labix.org/v2/mgo/bson"
)

func NewProject(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    in.ParseForm()
    project := objectspace.NewProject()
    studio := objectspace.NewStudio()
//    user := objectspace.NewUser()

//    currentuid := ExtractSingleValue(in.Form, "currentuid")
//    err := user.SetId(currentuid)
//    errors.append("current user", err)

    name := ExtractSingleValue(in.Form, "name")
    err := project.SetName(name)
    errors.append("name", err)

    studioid := ExtractSingleValue(in.Form, "studio")
    err = project.SetStudio(studioid)
    errors.append("studio", err)

    id := bson.NewObjectId().Hex()
    err = project.SetId(id)
    errors.append("id", err)

    //err = studio.SetId(studioid)
    filter := bson.M{"_id":studioid}
    err = studio.Restore(filter)
    errors.append("on studio restore", err)

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

    studio.AppendProject(id)
    err = studio.Update()
    if err != nil {
        errors.append("on studio update", err)
        WriteJsonResult(out, errors, "error")
        return
    }

    WriteJsonResult(out, project, "ok")
}

