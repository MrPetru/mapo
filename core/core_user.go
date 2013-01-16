package core

import (
    "mapo/log"
    "mapo/objectspace"

    "net/http"
    "strings"
    "strconv"

    "labix.org/v2/mgo/bson"
)

// NewUser crea un nuovo utente di mapo
// func NewUser(inValues values) interface{} {
func NewUser(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing NewUser function")

    // un contenitore per i dati relative al utente che si deve creare e che
    // verrano inseriti nella database
    user := objectspace.NewUser()

    // contenitore per gli errori
    errors := NewCoreErr()

    // procedura obligatori che estrae i dati codificati nel url della richiesta
    // o nei dati trasmessi da una forma e li inserisce nel in.Form
    //in.ParseForm()

    // extract and set username value
    username := in.FormValue("username")
    err := user.SetUsername(username)
    errors.append("username", err)

    // verifica e inserimento della passowrd nel contenitore del utente
    password := in.FormValue("password")
    err = user.SetPassword(password)
    errors.append("password", err)

    // get and set firstname
    firstname := in.FormValue("firstname")
    err = user.SetFirstname(firstname)
    errors.append("firstname", err)

    // get and set lastname
    lastname := in.FormValue("lastname")
    err = user.SetLastname(lastname)
    errors.append("lastname", err)


    // get and set description
    description := in.FormValue("description")
    err = user.SetDescription(description)
    errors.append("description", err)

    // get and set user email
    email := in.FormValue("email")
    err = user.SetEmail(email)
    errors.append("email", err)

	// per il momlento l'id e la soma md5 del username
    err = user.SetId(username)
	// se il SetUsername non avra errori allora anche qui non avremo.
    // questa volta non usiamo il error on quest modo errors.append("id", err)

    // TODO: tutte le altre operazioni per necesari per la registrazione utente

    // se i dati in entratta sono considerati errati, significa che la variabile errors
    // non è vota. a questo punto inviamo gli errori al cliente e concludeamo l'esecuzione
    // della richiesta.
    if len(errors) > 0 {
        WriteJsonResult(out, errors, "error")
        return
    }

    // se i dati in entrata sono stati accetati allora slava l'utente
    err = user.Save()
    if err != nil {
        errors.append("on save", err)
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
    errors.append("id", err)
    if len(errors) > 0 {
        WriteJsonResult(out, errors, "error")
        return
    }

    filter := bson.M{"_id":id}

    // oteniamo il utende dal database che poi vera aggiornato
    err = user.Restore(filter)
    if err != nil {
        errors.append("on restore", err)
        WriteJsonResult(out, errors, "error")
        return
    }

    password := ExtractSingleValue(in.Form, "password")
    if len(password) > 0 {
        err = user.SetPassword(password)
        errors.append("password", err)
    }

    firstname := ExtractSingleValue(in.Form, "firstname")
    if len(firstname) > 0 {
        err = user.SetFirstname(firstname)
        errors.append("firstname", err)
    }

    lastname := ExtractSingleValue(in.Form, "lastname")
    if len(lastname) > 0 {
        err = user.SetLastname(lastname)
        errors.append("lastname", err)
    }

    description := ExtractSingleValue(in.Form, "description")
    if len(description) > 0 {
        err = user.SetDescription(description)
        errors.append("description", err)
    }

    email := ExtractSingleValue(in.Form, "email")
    if len(email) > 0 {
        err = user.SetEmail(email)
        errors.append("email", err)
    }

    // aggiorniamo il valore del rating del utente
    strRating := ExtractSingleValue(in.Form, "rating")
    if len(strRating) > 0 {
        intRating, err := strconv.Atoi(strRating)
        if err != nil {
            errors.append("rating", err)
        } else {
            err = user.SetRating(intRating)
            if err != nil {
                errors.append("rating", err)
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
        errors.append("on update", err)
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
        errors.append("id", err)
    }

    // fermiamo l'esecuzione se fino a questo momento abbiamo incontrato qualche errore
    if len(errors) > 0{
        WriteJsonResult(out, errors, "error")
        return
    }

    filter := bson.M{"_id":id}

    // ricavare i dati del utente dalla database
    err = user.Restore(filter)
    if err != nil {
        errors.append("on restore", err)
        WriteJsonResult(out, errors, "error")
        return
    }

    log.Debug("%s", user.GetId())

    WriteJsonResult(out, user, "ok")
}

// GetUserAll restituisce una lista di tutti utenti nel database
// TODO: posibilita' di applicare dei filtri
// TODO: verificare she c'è bisogno di una funzione simile,
// ci sono casi quando serve avere una lista di utenti? sicuramente si, ma e questa l'utilizo giusto??
func GetUserAll(out http.ResponseWriter, in *http.Request){
    log.Msg("executing GetUserAll function")

    userList := objectspace.NewUserList()

    err := userList.Restore()
    if err != nil {
        errors := NewCoreErr()
        errors.append("on restore", err)
        WriteJsonResult(out, errors, "error")
        return
    }

    WriteJsonResult(out, userList, "ok")
}
