package main

import (
    "mapo/log"

    "os"
    "sync"
    "time"
    "regexp"
    "strings"
    "net/http"
)

/*
ServeMux, nasce dalla necessita di registrare dei handler differenziati anche
dal metodo http usato durante la richiesta dal parte del utente. Cosi lo stesso
url usato con il metodo POST ha un funzionamento diverso da una richiesta dove
si usa il metodo GET.

Un altra possibilità che ci offre questo handler personalizzato è di poter
interrompere il server usando la combinazione dei tasti CTRL+C
*/
type ServeMux struct {

    // lista dei handler registrati, con o senza autenticazione
    m map[string]Handler
    mVars map[string]map[int]string

    // il numero delle connessione attive in questo momento
    current_connections int
    lock sync.Mutex

    // il server è o no in fase di chiusura
    closing bool
}

func (mux *ServeMux) HandleFunc(method, path string, handle func(http.ResponseWriter, *http.Request)) {
    handlerFunc := new(http.HandlerFunc)
    *handlerFunc = handle

    pattern := createPattern(method, path)

    mux.m[pattern] = handlerFunc
    mux.mVars[pattern] = createUrlVars(path)
}

func (mux *ServeMux) Handle(method, path string, handler Handler) {
    pattern := createPattern(method, path)
    mux.m[pattern] = handler
    mux.mVars[pattern] = createUrlVars(path)
}

/*
createPattern, crea l'espressione regulare che si usa più tardi per trovare
il handler corretto per il path/risorsa richiesta.
*/
func createPattern(method, path string) string {
    pattern := "(?i)^"

    if method != "" {
        pattern = pattern + method + ":/"
    } else {
        pattern = pattern + "(GET|POST)" + ":/"
    }

    if len(path) > 1 {
        pathSlice := strings.Split(path[1:], "/")
        for _, v := range(pathSlice) {
            if v[0] == '{' {
                pattern = pattern + "[0-9a-z_\\ \\.\\+\\-]*/"
            } else {
                pattern = pattern + v + "/"
            }
        }
    }
    pattern = pattern + "$"
    return pattern
}

/*
createUrlVars, mappa le variabili inserite del url inserite al momento della
registrazione del handler. Le variabili sono segnati con le parentesi graffe
al interno delle quali si trova il nome della variabile. Questa mappa sarà usata
più tardi per passare i dati a forma di copie (chiave:valore) ai handler.
*/
func createUrlVars(path string) map[int]string {
    vlist := strings.Split(path, "/")

    data := make(map[int]string,0)

    for i, v := range(vlist) {
        if len(v) < 3 {
            continue
        }
        if v[0] == '{' && v[len(v)-1] == '}' {
            data[i] = v[1:len(v)-1]
        }
    }

    return data
}

/*
match, è usata per identificare quale dei handler corrisponde per un certo url.
Questa funzione fa utilizzo dei pattern (espressioni regolari).
*/
func (mux *ServeMux) match(r *http.Request) (Handler, string) {
    method := r.Method
    url := r.URL.Path

    if url[len(url)-1] != '/' {
        url = url + "/"
    }

    var handler Handler
    var pattern string

    for k, v := range(mux.m) {
        matching, _ := regexp.MatchString(k, method + ":" + url)
        if matching {
            handler = v
            pattern = k
            break
        }
    }

    if handler != nil {
        return handler, pattern
    }
    return http.NotFoundHandler(), ""
}

/*
NewServeMux, restituisce un nuovo miltiplixier personalizzato.
*/
func NewServeMux() *ServeMux {
    mux := new(ServeMux)
    mux.m = make(map[string]Handler, 0)
    mux.mVars = make(map[string]map[int]string, 0)

    return mux
}

/*
Handler, è un interfaccia che come funzionalità e struttura non è diversa dal
handler originale del modulo http di go.
TODO: Probabilmente è più corretto usare il http.Handler, resta da verificare.
*/
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}

// ServeHTTP e la funzione che vine eseguita come gorutine ogni volta che
// si deve processare qualche richiesta. Questa funzione soltanto si assicura
// che venga incrementato o decrementato il numero delle connessione attive e
// avvierà la funzione RequestHandler che processerà la richiesta del cliente.
// Comunque, il server http viene interrotto in maniera brutta ma senza alcun
// rischio. TODO: approfondire questa feature se servirà.
func (mux *ServeMux) ServeHTTP(out http.ResponseWriter, in *http.Request) {
    if !mux.closing {
        start := time.Now()
        defer func() {
            log.Info("time: %v for %s", time.Since(start), in.URL.Path)
        }()

        mux.lock.Lock()
        mux.current_connections++
        mux.lock.Unlock()

        defer func() {
            mux.lock.Lock()
            mux.current_connections--
            mux.lock.Unlock()
        }()

        handle, pattern := mux.match(in)
        if len(pattern) > 0 {
            in.ParseMultipartForm(0)
            urlValues := strings.Split(in.URL.Path, "/")
            for k, v := range(mux.mVars[pattern]) {
                in.Form[v] = []string{urlValues[k]}
            }
        }
        handle.ServeHTTP(out, in)
    }
}

// se viene richiesto che l'applicazione si deve chiudere, in questo momento si
// parla del commando CTRL+C dal terminale, potremmo corrompere i dati a colpa
// del'interruzione in maniera incorretta delle richieste in corso. La presente
// Funzione sta in ascolto per il segnale SIGINT dopo di che si assicura che il
// server venga chiuso non appena le connessione attive saranno zero.
func (mux *ServeMux) getSignalAndClose(c chan os.Signal) {

    _ = <-c
    log.Info("closing ...")
    mux.closing = true

    // TODO: send notification to load balancing that this node is unavailable

    for {
        if mux.current_connections == 0 {
            log.Info("bye ... :)")
            os.Exit(1)
        } else {
            log.Info("waiting for %d opened connections", mux.current_connections)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
