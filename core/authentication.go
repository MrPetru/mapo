package core

import (
    "net/http"
    "mapo/log"
    "fmt"
    "strings"
    "encoding/base64"
    "mapo/objectspace"
    "labix.org/v2/mgo/bson"
)

// RequestAuth richiede al client di autenticarsi
func RequestAuth(out http.ResponseWriter) {
    out.Header().Set("WWW-Authenticate", "Basic realm='mapomapo'")
    out.WriteHeader(401)
    fmt.Fprint(out, "not authorized!")
}

// Authenticator, se attivo, verifica l'entita del utente che richiede una
// risorsa. La verifica aviene atraverso il processo di login o procedimenti
// simile che restano da concordare, come per esempio OAuth.
func Authenticator(out http.ResponseWriter, in *http.Request) (http.ResponseWriter, *http.Request, bool) {

    log.Msg("authenticate for %v", in.URL.Path)

    authHeader, ok := in.Header["Authorization"]
    if !ok {
        RequestAuth(out)
        return out, in, false
    }
    encodedString := strings.Split(authHeader[0], " ")[1]

    encoder := base64.StdEncoding
    dec, _ := encoder.DecodeString(encodedString)
    decodedString := string(dec)

    var username, password string
    if tmp := strings.Split(string(decodedString), ":"); len(tmp) == 2 {
        username = tmp[0]
        password = tmp[1]
    } else {
        RequestAuth(out)
        return out, in, false
    }

    user := objectspace.NewUser()
    filter := bson.M{"username":username}
    err := user.Restore(filter)
    if err != nil {
        log.Debug("user restore error = %v", err)
        RequestAuth(out)
        return out, in, false
    }

    md5password := objectspace.Md5sum(password)
    if user.Password == md5password {
        // TODO:verificare se questo ParseForm non crea conflitto con altri ParseForm
        // che vengono chiamati nelle funzioni seguenti.
        // in.ParseForm() - questa funzione vera chiamata in automatico da ParseMultipartForm
        // se sara bisogno.
        in.ParseMultipartForm(0)
        in.Form["currentuid"] = []string{user.GetId()}
        return out, in, true
    }

    RequestAuth(out)
    return out, in, false
}

// in dipendenza dalla procedura usata, questa funzione potrebbe non
// essere indinspensabile.
func Login(out http.ResponseWriter, in *http.Request) {

    fmt.Fprint(out, "login")
}

// la funzione di deautenticazione non e' cosi semplice quando si usa il metodo
// Authorization, il header WWW-Authenticate deve essere cancellato da un script
// da parte del cliente.
func Logout(out http.ResponseWriter, in *http.Request) {

    fmt.Fprint(out, "logout")
}
