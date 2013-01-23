package core

import (
    "mapo/log"
    "mapo/objectspace"

    "net/http"
)

// GetUser restituisce un utente che è gia salvato nella database
// func GetUser(inValues values) interface{} {
func GetUser(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing GetUser function")

    errors := NewCoreErr()

    // cearmo un nuovo ogetto/contenitore per il utente richiesto
    user := objectspace.NewUser()

    // aggiorniamo il valore del id del utente, che servirà per ricavare l'utente
    // dal database
    id := in.FormValue("uid")
    err := user.SetId(id)
    if err != nil {
        errors.append("id", err)
    }

    // fermiamo l'esecuzione se fino a questo momento abbiamo incontrato qualche errore
    if len(errors) > 0{
        WriteJsonResult(out, errors, "error")
        return
    }

    // ricavare i dati del utente dalla database
    err = user.Restore()
    if err != nil {
        errors.append("on restore", err)
        WriteJsonResult(out, errors, "error")
        return
    }

    log.Debug("%s", user.GetId())

    WriteJsonResult(out, user, "ok")
}
