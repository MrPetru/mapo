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

/*
Authenticate, un wrapper che aiuta a registrare dei handler che sono protetti,
i handler che hanno bisogno che i utenti siano autentificati.

Processo di autenticazione avviene attraverso dei cookie di sessione (validi
soltanto fino alla chiusura del browser). I cookie vengono creati alla fine
della procedura di autenticazione guidata dal "Identity Provider"

    oauth - un codice che si usa per dimostrare che l'utente che fa la richiesta e lui.
    uid - l'id del cliente corrente
    sid - l'id dello studio attivo
    pid - l'id del progetto attivo

questi dati poi vengono passati al handler insieme a tutti i dati della
richiesta (per il momento sono inseriti direttamente in Form)

PROCESSO DI AUTENTICAZIONE

Processo di autenticazione può avviarsi nelle condizioni:
    1. Quando l'utente fa click su il bottone login
        l'utente inizia il suo lavoro con l'autenticazione.
    2. Quando qualcuno prova ad accedere qualche risorsa
        se una risorsa è protetta, al utenti li sarà chiesto di autentificasi.
        sarà reindirizzato alla pagina de login dove seguirà il processi
        guidato. Alla fine della procedura di autenticazione verrà
        reindirizzato alla risorsa richiesta inizialmente.

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

        uid := uidCookie.Value
        authid := authidCookie.Value

        cookie_secret, err := GlobalConfiguration.GetString("default", "cookiesecret")
        if err != nil {
            log.Debug("error gettiong cookie secret value %v", err)
            Forbidden(out)
            return
        }

        if objectspace.Md5sum(uid+cookie_secret) == authid {

            // ora verifichiamo se nella database esiste un utente con questo ID
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
    // da questa funzione conterà una mappa che avrà la chiave "error"
    if value := in.FormValue("error"); len(value) > 0 {
        log.Debug("user authorisation result: %s", value)
        http.Redirect(out, in, "/", 302)
        return
    }

    // se l'autenticazione è avvenuta con successo allora tra i dati ricevuti
    // in questo punto abbiamo il campo "code" che è il authorisation code che
    // useremo per chiedere a google l'access_token
    code := in.FormValue("code")

    var client_id, client_secret, cookie_secret string

    // interroghiamo il file di configurazione
    client_id, err := GlobalConfiguration.GetString("googleoauth", "clientid")
    client_secret, err = GlobalConfiguration.GetString("googleoauth", "clientsecret")
    cookie_secret, err = GlobalConfiguration.GetString("default", "cookiesecret")
    if len(client_id) < 1 || len(client_secret) < 1 || len(cookie_secret) < 1 {
        log.Debug("invalid configuration for OAuth")
        return
    }

    // ora che abbiamo il permesso del utente chiediamo a google il acces_token per poter
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
        log.Debug("accessData json Unmarshal err: %v", err)
    }

    responseGet, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?access_token=%s", accessData["access_token"]))
    if err != nil {
        log.Debug("get userData error: %v", err)
    }

    user := objectspace.NewUser()

    /*
    per lo scenario di questa interrogazione i dati ricevuti sono:
    {
        "id": "101...",
        "email": "...@gmail.com",
        "verified_email": true,
        "name": "Petru Ciobanu",
        "given_name": "Petru",
        "family_name": "Ciobanu",
        "link": "https://plus.google.com/101...",
        "picture": "https://lh4.googleusercontent.com/.../photo.jpg",
        "gender": "male",
        "birthday": "0000-00-00",
        "locale": "en"
    }
    */
    responseGetBody, err := ioutil.ReadAll(responseGet.Body)
    userData := map[string]interface{}{}
    err = json.Unmarshal(responseGetBody, &userData)
    if err != nil {
        log.Debug("userData json Unmarshal err: %v", err)
    }

    user.Oauthid = userData["id"].(string)
    user.Email = userData["email"].(string)
    user.Name = userData["name"].(string)
    user.Avatar = userData["picture"].(string)
    user.AccessToken = accessData["access_token"].(string)
    user.Oauthprovider = "google.com"
    user.CreateId()

    // verifica se il utente esiste nella database
    if tmpUser := user; tmpUser.Restore() != nil {
        err := user.Save()
        if err != nil {
            log.Debug("on user save err = %v", err)
            http.Redirect(out, in, "/", 302)
            return
        }
        // TODO: user is loged in for first time

    } else {
        err := user.Update()
        if err != nil {
            log.Debug("on user update err = %v", err)
            http.Redirect(out, in, "/", 302)
            return
        }
        // wellcom back user
        // update user in database
    }

    // TODO: a valid value for authentication cookie
    authid := objectspace.Md5sum(user.Id+cookie_secret)
    http.SetCookie(out, &http.Cookie{Name:"authid", Value: authid, Path: "/"})

    http.SetCookie(out, &http.Cookie{Name:"uid", Value: user.Id, Path: "/"})
    http.Redirect(out, in, "/", 302)
}

/*
per iniziare la procedura di autenticazione con un Identity Provider serve un
link formattato in un modo specifico.
*/
func Login(out http.ResponseWriter, in *http.Request) {

    var oauthprovider = in.FormValue("oauthprovider")

    switch oauthprovider {
        case "goauth":
            // use google url
            url := "https://accounts.google.com/o/oauth2/auth?scope=https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile&state=profile&redirect_uri=http://localhost:8081/oauth2callback&response_type=code&client_id=60876467348.apps.googleusercontent.com"
            http.Redirect(out, in, url, 302)
            return
        default:
            // per ora non fa niente
    }
}
