/*
Copyright 2013 Petru Ciobanu, Francesco Paglia, Lorenzo Pierfederici

This file is part of Mapo.

Mapo is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Mapo is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Mapo.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
    "mapo/database"
    "mapo/addon"
    "mapo/log"
    "mapo/core"
    
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "sync"
    "time"
//    "fmt"
    "regexp"
    "strings"
)

// main risponde del avvio del'applicazione e della sua
// registrazione come server in ascolto su la rete.
func main() {

    // settiamo il livello generale dei messaggi da visualizzare
    log.SetLevel("DEBUG")
    
    // istruiamo la database di creare una nuova connessione.
    database.NewConnection()
    log.Msg("created a new database connection")
    
    // al avvio del'applicazione si verifica la disponibilità dei addon
    // e si crea una lista globale che sarà passata verso altri moduli
    addons := addon.GetAll()
    addons = addons
    log.Msg("load addons and generate a list")
    
    // al momento del spegnimeto del'applicazione potremo trovarci con delle
    // connessione attive dal parte del cliente. Il handler personalizzato usato
    // qui, ci permette di dire al server di spegnersi ma prima deve aspettare
    // che tutte le richieste siano processate e la connessione chiusa
    mapoMuxer := NewMapoMux()
    
    server := &http.Server {
        Addr:   ":8081",
        Handler: mapoMuxer,
    }
    
    // TODO: register this node to load-balancing service
    
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT)
    
    // aviamo in una nuova gorutine la funzione che ascoltera per il segnale di
    // spegnimento del server
    go mapoMuxer.getSignalAndClose(c)

    mapoMuxer.HandleFunc("POST", "/admin/user", core.NewUser)
    mapoMuxer.HandleFunc("GET", "/admin/user/{id}", core.GetUser)
    mapoMuxer.HandleFunc("GET", "/admin/user", core.GetUserAll)
    
    log.Info("start listening for requests")
    
    // avviamo il server che processerà le richieste
    log.Msg("close server with message: %v", server.ListenAndServe())
}

// handler, personalizzato per il server http che ci permetterà di spegnere
// l'applicazione senza rischi o corruzione dei dati.
type MapoMux struct {

    mu sync.RWMutex
    m map[string]Handler

    // il numero delle connessione attive in questo momento
    current_connections int
    lock sync.Mutex
    
    // il server è o no in fase di chiusura
    closing bool
}

func (mux *MapoMux) HandleFunc(method, Path string, handle func(http.ResponseWriter, *http.Request) ) {
    // set a handler
    
    handlerFunc := new(http.HandlerFunc)
    *handlerFunc = handle
    
    newPath := "(?i)^"
    
    if method != "" {
        newPath = newPath + method + ":/"
    } else {
        newPath = newPath + "(GET|POST)" + ":/"
    }
    
    pathVars := strings.Split(Path[1:], "/")
    for _, v := range(pathVars) {
        if v[0] == '{' {
            newPath = newPath + "[0-9a-z_\\.\\+\\-]*/"
        } else {
            newPath = newPath + v + "/"
        }
    }
    
    mux.m[newPath] = handlerFunc
}

func (mux *MapoMux) match(r *http.Request) Handler {
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

func NewMapoMux() *MapoMux {
    mux := new(MapoMux)
    mux.m = make(map[string]Handler, 0)
    
    return mux
}

type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}

// ServeHTTP e la funzione che vine eseguita come gorutine ogni volta che
// si deve processare qualche richiesta. Questa funzione soltanto si assicura
// che venga incrementato o decrementato il numero delle connessione attive e
// avvierà la funzione RequestHandler che processerà la richiesta del cliente.
// Comunque, il server http viene interrotto in maniera brutta ma senza alcun
// rischio. TODO: approfondire questa feature se servirà.
func (mux *MapoMux) ServeHTTP(out http.ResponseWriter, in *http.Request) {
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
func (mux *MapoMux) getSignalAndClose(c chan os.Signal) {

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
