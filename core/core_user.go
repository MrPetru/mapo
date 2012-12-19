package core

import (
    "mapo/log"
    "mapo/objectspace"
    
    "net/http"
    "strings"
    "strconv"
    
    // utilizo di questo paccheto e soltanto temporaniamente
    // per creare un id che poi avrà una funzione specifica
    "labix.org/v2/mgo/bson"
)

// NewUser crea un nuovo utente di mapo
// func NewUser(inValues values) interface{} {
func NewUser(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing NewUser function")
    
    // un contenitore per i dati relative al utente che si deve creare e che
    // verrano inseriti nella database
    user := objectspace.NewUser()
    
    errors := NewCoreErr()
    
    in.ParseForm()
    
    login := ExtractSingleValue(in.Form, "login")
    err := user.SetLogin(login)
    if err != nil {
        errors.append("login", err.Error())
    }
    
    // verifica e inserimento della passowrd nel contenitore del utente
    password := ExtractSingleValue(in.Form, "password")
    err = user.SetPassword(password)
    if err != nil {
        errors.append("password", err.Error())
    }
    
    // get and set name
    name := ExtractSingleValue(in.Form, "name")
    err = user.SetName(name)
    if err != nil {
        errors.append("name", err.Error())
    }
    
    //id := name + "_" + login
    id := bson.NewObjectId().Hex()
    err = user.SetId(id)
    if err != nil {
        errors.append("id", err.Error())
    }
    
    // TODO: tutte le altre operazioni per necesari per la registrazione utente
    
    // se i dati in entratta sono considerati errati,
    // allora rimanda i dati indietro con l'errore
    if len(errors) > 0 {
        WriteJsonResult(out, errors, "error")
        return
    }
    
    // se i dati in entrata sono stati accetati allora slava l'utente
    err = user.Save()
    if err != nil {
        errors.append("on save", "data base error")
        log.Debug("%v", err)
        WriteJsonResult(out, errors, "error")
        return
    }
    
    WriteJsonResult(out, user, "ok")
}


// UpdateUser aggiorna il valori di un utenti nella database
func UpdateUser(out http.ResponseWriter, in *http.Request) {
    log.Msg("executing UpdateUser function")
    
    in.ParseForm()
    errors := NewCoreErr()
    
    user := objectspace.NewUser()

    id := strings.Split(in.URL.Path[1:], "/")[2]
    err := user.SetId(id)
    if err != nil {
        errors.append("id", err.Error())
    }
    
    // oteniamo il utende dal database che poi vera aggiornato
    err = user.Restore()
    if err != nil {
        errors.append("on restore", err.Error())
        WriteJsonResult(out, errors, "error")
        return
    }
    
    // aggiorniamo il valore del rating del utente
    strRating := ExtractSingleValue(in.Form, "rating")
    if len(strRating) > 0 {
        intRating, err := strconv.Atoi(strRating)
        if err != nil {
            errors.append("rating", err.Error())
        } else {
            err = user.SetRating(intRating)
            if err != nil {
                errors.append("rating", err.Error())
            }
        }
    }
    
    
    // update values for user
    // name
    // contacts
    // ...
    
    if len(errors) > 0 {
        WriteJsonResult(out, errors, "error")
        return
    }
    
    err = user.Update()
    if err != nil {
        errors.append("on update", err.Error())
        WriteJsonResult(out, errors, "error")
        return
    }
    
    WriteJsonResult(out, user, "ok")
}

// GetUser restituisce un utente che è gia salvato nella database
// func GetUser(inValues values) interface{} {
func GetUser(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing GetUser function")
    
    in.ParseForm()
    
    errors := NewCoreErr()
    
    // cearmo un nuovo ogetto/contenitore per il utente richiesto
    user := objectspace.NewUser()
    
    // aggiorniamo il valore del id del utente, che servirà per ricavare l'utente
    // dal database
    id := strings.Split(in.URL.Path[1:], "/")[2]
    err := user.SetId(id)
    if err != nil {
        errors.append("id", err.Error())
    }
    
    // fermiamo l'esecuzione se fino a questo momento abbiamo incontrato qualche errore
    if len(errors) > 0{
        WriteJsonResult(out, errors, "error")
        return
    }
    
    // ricavare i dati del utente dalla database
    err = user.Restore()
    if err != nil {
        errors.append("on restore", "data base error")
        WriteJsonResult(out, errors, "error")
        return
    }
    
    log.Debug("%s", user.GetId())

    WriteJsonResult(out, user, "ok")
}

// GetUserAll restituisce una lista di tutti utenti nel database
// TODO: posibilita' di applicare dei filtri
func GetUserAll(out http.ResponseWriter, in *http.Request){
    log.Msg("executing GetUserAll function")
    
    userList := objectspace.NewUserList()
    
    err := userList.Restore()
    if err != nil {
        errors := make(map[string][]string)
        errors["on restore"] = append(errors["on restore"], err.Error())
        WriteJsonResult(out, errors, "error")
        return
    }

    WriteJsonResult(out, userList, "ok")
}
