package core

import (
    "mapo/log"
    "mapo/objectspace"
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
func NewUser(inValues values) interface{} {

    log.Msg("executing NewUser function")
    
    // un contenitore per i dati relative al utente che si deve creare e che
    // verrano inseriti nella database
    user := objectspace.NewUser()
    
    // qui aviene anche la verifica che il nome abbia certe carateristiche
    login, err := inValues.GetSingleValue("login")
    if err != nil {
        inValues.SetError(err)
    }
    
    err = user.SetLogin(login)
    if err != nil {
        inValues.SetError(err)
    }
    
    // verifica e inserimento della passowrd nel contenitore del utente
    password, err := inValues.GetSingleValue("password")
    if err != nil {
        inValues.SetError(err)
    }
    
    err = user.SetPassword(password)
    if err != nil {
        inValues.SetError(err)
    }
    
    // TODO: get and set name
    // TODO: tutte le altre operazioni per necesari per la registrazione utente
    
    // se i dati in entratto sono considerati errati, esiste la chiave error
    // allora rimanda i dati indietro con l'errore
    if _, ok := inValues["error"]; ok {
        return inValues
    }
    
    user.SetId(user.GetLogin())
    
    // se i dati in entrata sono stati accetati allora slava l'utente
    err = user.Save()
    if err != nil {
        inValues.SetError(err)
        delete(inValues, "password")
        return inValues
    }
    
    // trasforma l'utente in un dato simile a quello in entratta
    userMap := user.ToMap()
    
    // escludiamo il password dal resultato che verà ristituita al cliente
    delete(userMap, "password")
    
    // ritorna i dati che sono stati salvati nella database
    return userMap
}


// UpdateUser aggiorna il valori di un utenti nella database
func UpdateUser(inValues values) interface{} {
    log.Msg("executing UpdateUser function")
    
    user := objectspace.NewUser()
    id, err := inValues.GetSingleValue("id")
    if err != nil {
        inValues.SetError(err)
        return inValues
    }
    
    user.SetId(id)
    err = user.Restore()
    if err != nil {
        inValues.SetError(err)
    }
    
    // update values for user
    // name
    // contacts
    // ...
    
    user.SaveUpdate()
    return user
}

// GetUser restituisce un utente che è gia salvato nella database
func GetUser(inValues values) interface{} {
    log.Msg("executing GetUser function")
    
    // cearmo un nuovo ogetto/contenitore per il utente richiesto
    user := objectspace.NewUser()
    
    // aggiorniamo il valore del id del utente, che servirà per ricavare l'utente
    // dal database
    id, err := inValues.GetSingleValue("id")
    if err != nil {
        inValues.SetError(err)
        return inValues
    }
    
    user.SetId(id)
    
    // ricavare i dati del utente dalla database
    err = user.Restore()
    if err != nil {
        inValues.SetError(err)
        return inValues
    }
    
    log.Debug("%s", user.GetId())
    
    // transformiamo tutto in un tipo di dato simile a quello in entrata
    // se anche non necessario, questa operazione rende più uniforma l'interazione
    // tra il module superiore è questo.
    userMap := user.ToMap()
    
    // canceliamo la password dal oggetto ritornato al client
    // questo dato non serve per cliente
    delete(userMap, "password")
    
    return userMap
}

// GetUserAll restituisce una lista di untenti
func GetUserAll() interface{} {
    log.Msg("executing GetUserAll function")
    
    userList := objectspace.NewUserList()
    userList.Restore()
    
    return userList
}
