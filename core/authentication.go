package core

import (
    "net/http"
    "mapo/log"
    "fmt"
    "strings"
    "encoding/base64"
    "mapo/objectspace"
    "labix.org/v2/mgo/bson"
    "bytes"
    "io"
)

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

    var resultBuf bytes.Buffer
    b64Decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(encodedString))
    io.Copy(&resultBuf, b64Decoder)
    decodedString := resultBuf.String()

    var login, password string
    if tmp := strings.Split(string(decodedString), ":"); len(tmp) == 2 {
        login = tmp[0]
        password = tmp[1]
    } else {
        RequestAuth(out)
        return out, in, false
    }

    user := objectspace.NewUser()
    filter := bson.M{"login":login}
    err := user.Restore(filter)
    if err != nil {
        log.Debug("user restore error = %v", err)
        RequestAuth(out)
        return out, in, false
    }

    if user.Password == password {//[:len(password)-1] {
        // TODO:verificare se questo ParseForm non crea conflitto con altri ParseForm
        // che vengono chiamati nelle funzioni seguenti.
        in.ParseForm()
        in.Form["currentuid"] = []string{user.GetId()}
        return out, in, true
    }

    RequestAuth(out)
    return out, in, false
}

func Login(out http.ResponseWriter, in *http.Request) {

    fmt.Fprint(out, "login")
}

func Logout(out http.ResponseWriter, in *http.Request) {

    fmt.Fprint(out, "logout")
}
