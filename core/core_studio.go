package core

import (
    "mapo/log"
    "mapo/objectspace"
    
    "net/http"
    
//    // utilizo di questo paccheto e soltanto temporaniamente
//    // per creare un id che poi avrà una funzione specifica
//    "labix.org/v2/mgo/bson"
)

func NewStudio(out http.ResponseWriter, in *http.Request) {
    // create new studio
    log.Msg("executing NewStudio function")
    
    in.ParseForm()
    errors := NewCoreErr()
    
    // creamo un nuovo contenitore di tipo studio
    studio := objectspace.NewStudio()
    
    name := ExtractSingleValue(in.Form, "name")
    err := studio.SetName(name)
    if err != nil {
        errors.append("name", err)
    }
    
    userid := ExtractSingleValue(in.Form, "userid")
    err = studio.SetUserid(userid)
    if err != nil {
        errors.append("userid", err)
    }
    
    id := name
    err = studio.SetId(id)
    if err != nil {
        errors.append("id", err)
    }
    
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
    
    // update user
    user := objectspace.NewUser()
    err = user.SetId(userid)
    if err != nil {
        errors.append("userid", err)
    }
    
    err = user.Restore()
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

//func GetStudio(inValues values) interface{} {
//    // create new studio
//    
//    return nil
//}

//func GetStudioAll(inValues values) interface{} {
//    // create new studio
//    
//    return nil
//}

//func UpdateStudio(inValues values) interface{} {
//    // create new studio
//    
//    return nil
//}
