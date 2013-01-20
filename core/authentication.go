package core

import (
    "net/http"
    "net/url"
    "mapo/log"
    "fmt"
    "encoding/json"
    "mapo/objectspace"
    "labix.org/v2/mgo/bson"
    "io/ioutil"
)

func Forbidden(out http.ResponseWriter) {
    out.Header().Set("Content-Type","application/json;charset=UTF-8")

    http.SetCookie(out, &http.Cookie{Name:"authid", Value: "", Path: "/"})
    http.SetCookie(out, &http.Cookie{Name:"uid", Value: "", Path: "/"})

    out.WriteHeader(403)
    message := make(map[string][]string, 0)
    message["authentication"] = []string{"invalid user"}
    WriteJsonResult(out, message, "error")
}

// Authenticator, se attivo, verifica l'entita del utente che richiede una
// risorsa. La verifica aviene atraverso il processo di login o procedimenti
// simile che restano da concordare, come per esempio OAuth.
func Authenticator(out http.ResponseWriter, in *http.Request) (http.ResponseWriter, *http.Request, bool) {

    log.Msg("authenticate for %v", in.URL.Path)

    var authidCookie, uidCookie *http.Cookie
    var err error
    if authidCookie, err = in.Cookie("authid"); err != nil {
        Forbidden(out)
        return out, in, false
    }

    if uidCookie, err = in.Cookie("uid"); err != nil {
        Forbidden(out)
        return out, in, false
    }

    log.Debug("authidCookie = %v \n uidCookie = %v", authidCookie, uidCookie)

    uid := uidCookie.Value
    authid := authidCookie.Value
    log.Debug("authid = %v \n uid = %v", authid, uid)

    cookie_secret, err := GlobalConfiguration.GetString("default", "cookiesecret")
    if err != nil {
        log.Debug("error gettiong cookie secret value %v", err)
        Forbidden(out)
        return out, in, false
    }

    if objectspace.Md5sum(uid+cookie_secret) == authid {

        // ora verifchiamo se nella database esiste un utente con questo ID
        user := objectspace.NewUser()
        err := user.Restore(bson.M{"_id":uid})
        if err == nil {

            // se fin qua tutt e' a posto allora...
            in.ParseMultipartForm(0)
            in.Form["currentuid"] = []string{uid}
            return out, in, true
        }
    }

    Forbidden(out)
    return out, in, false

}

// l'utente viene reindirizato verso questa funzione dopo la procedura
// di autenticazione guidata da google.
func OAuthCallBack(out http.ResponseWriter, in *http.Request) {

    in.ParseMultipartForm(0)

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
        postData := url.Values{"code":{code}, "client_id":{client_id}, "client_secret":{client_secret}, "redirect_uri":{"http://localhost:8081/oauth2callback"}, "grant_type":{"authorization_code"}}
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
        if tmpUser := objectspace.NewUser(); tmpUser.Restore(bson.M{"_id":userData.Id}) != nil {
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

