package core

import (
    "net/http"
    "mapo/objectspace"
    "labix.org/v2/mgo/bson"
    "strings"
)

func NewProject(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    in.ParseForm()

    project := objectspace.NewProject()

    name := in.FormValue("name")
    err := project.SetName(name)
    errors.append("name", err)

    description := in.FormValue("description")
    err = project.SetDescription(description)
    errors.append("description", err)

    //var studioID string // get studio ID
    sidCookie, err := in.Cookie("sid")
    studioID := sidCookie.Value

    studio := objectspace.NewStudio()
    err = studio.SetId(studioID)
    errors.append("studioid", err)

    err = studio.Restore(bson.M{"_id":studioID})
    errors.append("on studio restore", err)

    if len(errors) > 0 {
        WriteJsonResult(out, errors, "error")
        return
    }

    id := objectspace.Md5sum(studioID + name)
    err = project.SetId(id)
    errors.append("id", err)

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

func GetProjectAll(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    sidCookie, err := in.Cookie("sid")
    if err != nil {
        errors.append("studio", "no active studio in current session")
        WriteJsonResult(out, errors, "error")
        return
    }
    studioID := sidCookie.Value

    studio := objectspace.NewStudio()
    studio.SetId(studioID)

    studio.Restore(bson.M{"_id":studioID})

    projects := studio.Projects

    filter := bson.M{"_id":bson.M{"$in":projects}}

    projectlist, err := objectspace.ProjectRestorList(filter)

    if err != nil {
        WriteJsonResult(out, err, "error")
    }

    WriteJsonResult(out, projectlist, "ok")
}

func GetProject(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    id := strings.Split(in.URL.Path[1:], "/")[2]
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

    studioFilter := bson.M{"projects":id}

    studio := objectspace.NewStudio()
    err = studio.Restore(studioFilter)
    if err != nil {
        return
    }

    if studio.Id != sid {
        return
    }

    project := objectspace.NewProject()
    err = project.Restore(bson.M{"_id":id})
    if err != nil {
        return
    }

    WriteJsonResult(out, project, "ok")
}
