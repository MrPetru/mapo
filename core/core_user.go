package core

import (
    "mapo/log"
    "mapo/objectspace"
    
    "net/http"
    "fmt"
)

/*
TODO: dichiarere le interface da usare per user
type SetterGetter interface {
    SetLogin(string) error
    GetLogin() string
    ...
}
*/

//NewUser crea un nuovo utente di mapo
//func NewUser(inValues values) interface{} {
func NewUser(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing NewUser function")
    
    // un contenitore per i dati relative al utente che si deve creare e che
    // verrano inseriti nella database
    user := objectspace.NewUser()
    
    errors := make(map[string][]string)
    
    in.ParseForm()
    
    login := ExtractSingleValue(in.Form, "login")
    err := user.SetLogin(login)
    if err != nil {
        errors["login"] = append(errors["login"], err.Error())
    }
    
    // verifica e inserimento della passowrd nel contenitore del utente
    password := ExtractSingleValue(in.Form, "password")
    err = user.SetPassword(password)
    if err != nil {
        errors["password"] = append(errors["password"], err.Error())
    }
    
    // get and set name
    name := ExtractSingleValue(in.Form, "name")
    err = user.SetName(name)
    if err != nil {
        errors["name"] = append(errors["name"], err.Error())
    }
    
    // TODO: tutte le altre operazioni per necesari per la registrazione utente
    
    // se i dati in entratto sono considerati errati, esiste la chiave error
    // allora rimanda i dati indietro con l'errore
    for _, _ = range(errors) {
        fmt.Fprint(out, errors)
        delete(in.Form, "password")
        fmt.Fprint(out, in.Form)
        return
    }
    
    user.SetId(user.GetLogin())
    
    // se i dati in entrata sono stati accetati allora slava l'utente
    err = user.Save()
    if err != nil {
        errors["on save"] = append(errors["on save"], err.Error())
        fmt.Fprint(out, errors)
        delete(in.Form, "password")
        fmt.Fprint(out, in.Form)
        return
    }
    
    // trasforma l'utente in un ogetto più dinamico
    userMap := user.ToMap()
    
    // escludiamo il password dal resultato che verà ristituita al cliente
    delete(userMap, "password")
    
    // ritorna i dati che sono stati salvati nella database
    fmt.Fprint(out, userMap)
}


//// UpdateUser aggiorna il valori di un utenti nella database
//func UpdateUser(inValues values) interface{} {
//    log.Msg("executing UpdateUser function")
//    
//    user := objectspace.NewUser()
//    id, err := inValues.GetSingleValue("id")
//    if err != nil {
//        inValues.SetError(err)
//        return inValues
//    }
//    
//    user.SetId(id)
//    err = user.Restore()
//    if err != nil {
//        inValues.SetError(err)
//    }
//    
//    // update values for user
//    // name
//    // contacts
//    // ...
//    
//    user.SaveUpdate()
//    return user
//}

// GetUser restituisce un utente che è gia salvato nella database
//func GetUser(inValues values) interface{} {
func GetUser(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing GetUser function")
    
    in.ParseForm()
    
    errors := make(map[string][]string)
    
    // cearmo un nuovo ogetto/contenitore per il utente richiesto
    user := objectspace.NewUser()
    
    // aggiorniamo il valore del id del utente, che servirà per ricavare l'utente
    // dal database
    id := ExtractSingleValue(in.Form, "id")
    err := user.SetId(id)
    if err != nil {
        errors["id"] = append(errors["id"], err.Error())
    }
    
    for _, _ = range(errors) {
        fmt.Fprint(out, errors)
        delete(in.Form, "password")
        fmt.Fprint(out, in.Form)
        return
    }
    
    // ricavare i dati del utente dalla database
    err = user.Restore()
    if err != nil {
        errors["on restore"] = append(errors["on restore"], err.Error())
        fmt.Fprint(out, errors)
        delete(in.Form, "password")
        fmt.Fprint(out, in.Form)
    }
    
    log.Debug("%s", user.GetId())
    
    // transformiamo tutto in un tipo di dato simile a quello in entrata
    // se anche non necessario, questa operazione rende più uniforma l'interazione
    // tra il module superiore è questo.
    userMap := user.ToMap()
    
    // canceliamo la password dal oggetto ritornato al client
    // questo dato non serve per cliente
    delete(userMap, "password")
    
    fmt.Fprint(out, userMap)
}

// GetUserAll restituisce una lista di tutti utenti nel database
// TODO: posibilita' di applicare dei filtri
func GetUserAll(out http.ResponseWriter, in *http.Request){
    log.Msg("executing GetUserAll function")
    
    userList := objectspace.NewUserList()
    userList.Restore()
    
    userMapList := make([]map[string]interface{}, 0)
    
    for _, u := range(userList) {
        userMap := u.ToMap()
        userMapList = append(userMapList, userMap)
    }
    
    fmt.Fprint(out, userMapList)
}
