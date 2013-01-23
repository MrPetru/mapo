package core

import (
    "mapo/objectspace"
	"mapo/log"

    "net/http"
    "net/url"
    "fmt"
    "encoding/json"
    "io/ioutil"
)

/*
Forbidden, una scorciatoia usata per ritornare al cliente il messaggio che lui
non e' autorizzato ad accedere questa risorsa o probabilmente che lui non ha
fatto login.
*/
func Forbidden(out http.ResponseWriter) {
    out.Header().Set("Content-Type","application/json;charset=UTF-8")

    http.SetCookie(out, &http.Cookie{Name:"authid", Value: "", Path: "/"})
    http.SetCookie(out, &http.Cookie{Name:"uid", Value: "", Path: "/"})

    out.WriteHeader(403)
    WriteJsonResult(out, "not authorised", "error")
}

// Authenticator, se attivo, verifica l'entita del utente che richiede una
// risorsa. La verifica aviene atraverso il processo di login o procedimenti
// simile che restano da concordare, come per esempio OAuth.

/*
Authenticate, un wrapper che aiuta a registrare dei handler che sono protetti,
i handler che hanno bisogno che i utenti siano autentificati.

il processo di autenticazione avviene attraverso dei cookie di sessione (validi
soltanto fino alla chiusura del browser). I cookie vengono creati alla fine
della procedura di autenticazione guidata dal "Identity Provider"

oauth - un codice che si usa per dimostrare che l'utente che fa la richiesta e
lui.

uid - l'id del cliente corrente

sid - l'id dello studio attivo

pid - l'id del progetto attivo

questi dato poi vengono passati al handler insieme a tutti i dati della
richiesta (per il momento sono inseriti direttamente in Form)
*/
func Authenticate(handleFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

    return func(out http.ResponseWriter, in *http.Request) {
        log.Msg("authenticate for %v", in.URL.Path)

        var authidCookie, uidCookie *http.Cookie
        var err error
        if authidCookie, err = in.Cookie("authid"); err != nil {
            Forbidden(out)
            return
        }

        if uidCookie, err = in.Cookie("uid"); err != nil {
            Forbidden(out)
            return
        }

        log.Debug("authidCookie = %v \n uidCookie = %v", authidCookie, uidCookie)

        uid := uidCookie.Value
        authid := authidCookie.Value
        log.Debug("authid = %v \n uid = %v", authid, uid)

        cookie_secret, err := GlobalConfiguration.GetString("default", "cookiesecret")
        if err != nil {
            log.Debug("error gettiong cookie secret value %v", err)
            Forbidden(out)
            return
        }

        if objectspace.Md5sum(uid+cookie_secret) == authid {

            // ora verifchiamo se nella database esiste un utente con questo ID
            user := objectspace.NewUser()
            user.SetId(uid)
            err := user.Restore()
            if err == nil {

                // se fin qua tutt e' a posto allora...
                in.Form["currentuid"] = []string{uid}
                handleFunc(out, in)
                return
            }
        }

        Forbidden(out)
        return
    }

}

/*
l'utente viene reindirizzato verso questa funzione dopo la procedura
di autenticazione guidata da google.

TODO: raffinare questa funzione
*/
func OAuthCallBack(out http.ResponseWriter, in *http.Request) {

    // nel caso che l'utente non consente l'accesso ai suoi dati, il dati ricevuti
    // da questa funzione contera una mapa che avra la chiave "error"
    if value := in.FormValue("error"); value == "" {

        code := in.FormValue("code")

        var client_id, client_secret, cookie_secret string
        client_id, err := GlobalConfiguration.GetString("googleoauth", "clientid")
        client_secret, err = GlobalConfiguration.GetString("googleoauth", "clientsecret")
        cookie_secret, err = GlobalConfiguration.GetString("default", "cookiesecret")
        if len(client_id) < 1 || len(client_secret) < 1 {
            log.Debug("invalid configuration for OAuth")
            return
        }

        // ora che abbiamo il permesso del utente chediamo a google il acces_token per poter
        // accedere ai deti del utente
        postData := url.Values{"code":{code}, "client_id":{client_id},
                "client_secret":{client_secret},
                "redirect_uri":{"http://localhost:8081/oauth2callback"},
                "grant_type":{"authorization_code"}}

        response, err := http.PostForm("https://accounts.google.com/o/oauth2/token", postData)
        if err != nil {
            log.Debug("get token error %v", err)
            return
        }
        defer response.Body.Close()

        responseBody, _ := ioutil.ReadAll(response.Body)

        accessData := map[string]interface{}{}
        err = json.Unmarshal(responseBody, &accessData)
        if err != nil {
            log.Debug("access data json Unmarshal err: %v", err)
        }

        userData := objectspace.NewUser()
        responseGet, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?access_token=%s", accessData["access_token"]))
        if err != nil {
            log.Debug("error on get user data: %v", err)
        }

        responseGetBody, err := ioutil.ReadAll(responseGet.Body)
        err = json.Unmarshal(responseGetBody, &userData)
        if err != nil {
            log.Debug("user data json Unmarshal err: %v", err)
        }

        userData.AccessToken = accessData["access_token"].(string)

        userData.Oauthprovider = "google.com"
        userData.CreateId()

        // verifica se il utente esiste nella database
        if tmpUser := userData; tmpUser.Restore() != nil {
            err := userData.Save()
            if err != nil {
                log.Debug("on user save err = %v", err)
                return
            }
            // user is loged in for first time
        } else {
            err := userData.Update()
            if err != nil {
                log.Debug("on user update err = %v", err)
                return
            }
            // wellcom back user
            // update user in database
        }

        log.Debug("user data = %v", userData)

        // TODO: a valid value for authentication cookie
        authid := objectspace.Md5sum(userData.Id+cookie_secret)
        http.SetCookie(out, &http.Cookie{Name:"authid", Value: authid, Path: "/"})

        http.SetCookie(out, &http.Cookie{Name:"uid", Value: userData.Id, Path: "/"})
        http.Redirect(out, in, "/", 302)

        return
    }

    // TODO: cosa succede se l'utente non acceta di l'authenticazione?
    // redirect alla pagina di login o alla pagina / ?
    log.Debug("form google: %v", in.Form)
}

