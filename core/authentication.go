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

// RequestAuth richiede al client di autenticarsi
func RequestAuth(out http.ResponseWriter) {
    out.Header().Set("WWW-Authenticate", "Basic realm='mapomapo'")
    out.WriteHeader(401)
    fmt.Fprint(out, "not authorized!")
}

func Forbidden(out http.ResponseWriter) {
    out.WriteHeader(403)
    fmt.Fprint(out, "not authorized!")
}

// Authenticator, se attivo, verifica l'entita del utente che richiede una
// risorsa. La verifica aviene atraverso il processo di login o procedimenti
// simile che restano da concordare, come per esempio OAuth.
func Authenticator(out http.ResponseWriter, in *http.Request) (http.ResponseWriter, *http.Request, bool) {

    log.Msg("authenticate for %v", in.URL.Path)

    if c, err := in.Cookie("authid"); err == nil {
        authid := c.Value
        log.Debug("authid = %v", authid)
        user := objectspace.NewUser()
        filter := bson.M{"_id":authid}
        err = user.Restore(filter)
        if err == nil {
            in.ParseMultipartForm(0)
            in.Form["currentuid"] = []string{user.GetId()}
            log.Debug("form = %v", in.Form)
            return out, in, true
        }
    }

    Forbidden(out)
    return out, in, false

}

// dipende dalla procedura usata, questa funzione potrebbe non
// essere indinspensabile.
func Login(out http.ResponseWriter, in *http.Request) {

    errors := NewCoreErr()

    in.ParseMultipartForm(0)
    username := in.FormValue("username")
    password := in.FormValue("password")

    md5password := objectspace.Md5sum(password)

    user := objectspace.NewUser()
    filter := bson.M{"username":username}
    err := user.Restore(filter)
    if err != nil || user.Password != md5password {
        errors.append("login", "wrong credentiales")
        WriteJsonResult(out, errors, "error")
        return
    }

    // TODO: a valid value for authentication cookie
    authid := user.Id

    http.SetCookie(out, &http.Cookie{Name:"authid", Value: authid, Path: "/"})

    WriteJsonResult(out, nil, "ok")
}

// la funzione di deautenticazione non e' cosi semplice quando si usa il metodo
// Authorization, il header WWW-Authenticate deve essere cancellato da un script
// da parte del cliente.
func Logout(out http.ResponseWriter, in *http.Request) {

    fmt.Fprint(out, "logout")
}

// l'utente viene reindirizato verso questa funzione dopo la procedura
// di autenticazione guidata da google.
func OAuthCallBack(out http.ResponseWriter, in *http.Request) {

    in.ParseMultipartForm(0)

    if value := in.FormValue("error"); value == "" {

        code := in.FormValue("code")

        var client_id, client_secret string
        client_id, err := GlobalConfiguration.GetString("googleoauth", "clientid")
        log.Debug("client_id %v err %v", client_id, err)

        client_secret, err = GlobalConfiguration.GetString("googleoauth", "clientsecret")
        log.Debug("client_secret %v err %v", client_secret, err)

        if len(client_id) < 1 || len(client_secret) < 1 {
            log.Debug("invalid configuration for OAuth")
            return
        }

        data := url.Values{"code":{code}, "client_id":{client_id}, "client_secret":{client_secret}, "redirect_uri":{"http://localhost:8081/oauth2callback"}, "grant_type":{"authorization_code"}}
        resp, err := http.PostForm("https://accounts.google.com/o/oauth2/token", data)
        if err != nil {
            log.Debug("get token error %v", err)
        }
        defer resp.Body.Close()

        rbody, err := ioutil.ReadAll(resp.Body)

        v := map[string]interface{}{}
        err = json.Unmarshal(rbody, &v)
        log.Debug("response body %v, %v", v, err)

        // get user data
        userData := map[string]interface{}{}
        responseGet, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?access_token=%s", v["access_token"]))
        responseGetBody, err := ioutil.ReadAll(responseGet.Body)
        err = json.Unmarshal(responseGetBody, &userData)

        log.Debug("user data = %v", userData)

        return
    }

    log.Debug("form google: %v", in.Form)
}
