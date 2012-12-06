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
    "fmt"
    "strings"
    "encoding/json"
)

// main risponde del avvio del'applicazione e della sua
// registrazione come server in ascolto su la rete.
func main() {

    // settiamo il livello generale dei messaggi da visualizzare
    log.SetLevel("DEBUG")
    
    // istruiamo la database di creare una nuova connessione.
    //
    // NOTE: metto alla valutazione il fatto che l'attivazione del dattabase
    // deve 
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
    h := new(handler)
    
    // TODO: register this node to load-balancing service
    
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT)
    
    // aviamo in una nuova gorutine la funzione che ascoltera per il segnale di
    // spegnimento del server
    go h.getSignalAndClose(c)
    
    log.Info("start listening for requests")
    
    // avviamo il server che processerà le richieste
    log.Msg("close server with message: %v", http.ListenAndServe(":8081", h))
}

// handler, personalizzato per il server http che ci permetterà di spegnere
// l'applicazione senza rischi o corruzione dei dati.
type handler struct {

    // il numero delle connessione attive in questo momento
    current_connections int
    lock sync.Mutex
    
    // il server è o no in fase di chiusura
    closing bool
}

/*
RequestHandler processa in maniera separata ogni richiesta verso il server.
Questo è il primo passaggio che colleziona i dati della richiesta che poi
vengono indirizzati verso il modulo giusto. La risposta del modulo vera inviata
al cliente ma prima trasformerà il risultato il un formato conosciuto al cliente
, es: json

TODO: decidere se il processo di autenticazione deve essere qui o da un altra
parte
*/
func (h *handler) RequestHandler(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing RequestHandler function")
    
    // i dati ottenuti dal ParseForm e ParseMultipartForm sono passati ai moduli
    // specifici come core, api, webui.
    // al momento la richiesta del utente può essere ridotta a una seria di dati del tipo
    // chiave=valore, e una lista di file. I questa situazione i dati sono
    // contenuti in una delle due variabile: in.Form e in.MultipartForm
    in.ParseForm()
    in.ParseMultipartForm(0)
    
    // un passaggio importante qui è il modulo di autenticazione. Ogni richiesta
    // passere prima una verifica del utente
    // TODO: il modulo di autenticazione
    
    resourcePath := strings.Split(in.URL.Path, "/")
    log.Debug("%s, %d", resourcePath, len(resourcePath))
    
    requestMethod := in.Method
    
    // Qui si identifica il modulo a cui li si deve passare il controllo.
    switch m := resourcePath[1]; m {
        case "admin":
            result := core.Start(resourcePath[2:], requestMethod, in.Form)
            jsonResult, _ := json.Marshal(result)
            out.Header().Set("Content-Type","text/x-json")
            fmt.Fprint(out, string(jsonResult))
        case "api":
            // avvia modulo api
        case "webui":
            // avvia modulo webui
        default:
            // probabilmente qui servirà un altro modulo o possiamo ritornare
            // l'errore pagenotfound
            http.Error(out, "404 page not found", http.StatusNotFound)
            
    }
    return
}

// ServeHTTP e la funzione che vine eseguita come gorutine ogni volta che
// si deve processare qualche richiesta. Questa funzione soltanto si assicura
// che venga incrementato o decrementato il numero delle connessione attive e
// avvierà la funzione RequestHandler che processerà la richiesta del cliente.
// Comunque, il server http viene interrotto in maniera brutta ma senza alcun
// rischio. TODO: approfondire questa feature se servirà.
func (h *handler) ServeHTTP(out http.ResponseWriter, in *http.Request) {
    
    if !h.closing {
        h.lock.Lock()
        h.current_connections++
        h.lock.Unlock()
        
        defer func() {
            h.lock.Lock()
            h.current_connections--
            h.lock.Unlock()
        }()
        
        h.RequestHandler(out, in)
    }
}

// se viene richiesto che l'applicazione si deve chiudere, in questo momento si
// parla del commando CTRL+C dal terminale, potremmo corrompere i dati a colpa
// del'interruzione in maniera incorretta delle richieste in corso. La presente
// Funzione sta in ascolto per il segnale SIGINT dopo di che si assicura che il
// server venga chiuso non appena le connessione attive saranno zero.
func (h *handler) getSignalAndClose(c chan os.Signal) {

    _ = <-c
    log.Info("closing ...")
    h.closing = true
    
    // TODO: send notification to load balancing that this node is unavailable
    
    for {
        if h.current_connections == 0 {
            log.Info("bye ... :)")
            os.Exit(1)
        } else {
            log.Info("waiting for %d opened connections", h.current_connections)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
