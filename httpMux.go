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

// handler, personalizzato per il server http che ci permetterà di spegnere
// l'applicazione senza rischi o corruzione dei dati.
type ServeMux struct {

    // lista dei handler registrati, con o senza autenticazione
    m map[string]Handler

    // se non è nil, contiene la funzione autentifica il utente che fa la richiesta
    // questo identiticatore si usarà di seguito per registrare dei handeler che hanno
    // bisogno di autenticazione.
    authenticator func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request, bool)

    // il numero delle connessione attive in questo momento
    current_connections int
    lock sync.Mutex

    // il server è o no in fase di chiusura
    closing bool
}

type handlerWithAuthentication struct {
    f func(http.ResponseWriter, *http.Request)
    auth func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request, bool)
}
func (hwa handlerWithAuthentication) ServeHTTP(out http.ResponseWriter, in *http.Request) {
    if hwa.auth != nil {
        if o, i, ok := hwa.auth(out, in); ok {
            hwa.f(o, i)
        }
    } else {
        hwa.f(out, in)
    }
}

type handlerWithoutAuthentication func(http.ResponseWriter, *http.Request)
func (hwoa handlerWithoutAuthentication) ServeHTTP(out http.ResponseWriter, in *http.Request) {
    hwoa(out, in)
}

// SetAuthenticator allega al mux una funzione da noi definita che si userà
// per autenticare il utente che fa la richeista.
func (mux *ServeMux) SetAuthenticator(auth func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request, bool)) {
    mux.authenticator = auth
}

// Handler registra un handler che di defautl usa l'autenticazione.
func (mux *ServeMux) HandleFunc(method, path string, handle func(http.ResponseWriter, *http.Request) ) {

    handlerFunc := new(handlerWithAuthentication)
    handlerFunc.auth = mux.authenticator
    handlerFunc.f = handle

    pattern := createPattern(method, path)

    mux.m[pattern] = handlerFunc
}

// HandleFuncNoAuth registra in modo explicito un handler che non ha bisogno
// di un autente autenticato.
func (mux *ServeMux) HandleFuncNoAuth(method, path string, handle func(http.ResponseWriter, *http.Request) ) {

    handlerFunc := new(handlerWithoutAuthentication)
    *handlerFunc = handle

    pattern := createPattern(method, path)

    mux.m[pattern] = handlerFunc
}

func (mux *ServeMux) Handle(method, path string, handler Handler) {
    pattern := createPattern(method, path)
    mux.m[pattern] = handler
}

// createPattern, per il momento è una funzione limitata a creare la regola
// in base a quale si andrà a verificare se il path corisponde a un certo handler.
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

// match è usata nel processo di identificazione delle handler necessario da
// eseguire per una certa risorsa identificata dal url.
func (mux *ServeMux) match(r *http.Request) Handler {
    method := r.Method
    url := r.URL.Path

    if url[len(url)-1] != '/' {
        url = url + "/"
    }

    var handler Handler

    for k, v := range(mux.m) {
        matching, _ := regexp.MatchString(k, method + ":" + url)
        if matching {
            handler = v
            break
        }
    }

    if handler != nil {
        return handler
    }
    return http.NotFoundHandler()
}

// NewServeMux restituisce un nuovo mux, molto simile al mux originale del
// modulo http di go.
func NewServeMux() *ServeMux {
    mux := new(ServeMux)
    mux.m = make(map[string]Handler, 0)

    return mux
}

// Handler è la ridefenizione del Handler del modulo http di go.
// usato per provare a mantenere la logica del quel modulo.
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

        handle := mux.match(in)
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
