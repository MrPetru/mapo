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
    "mapo/webui"

    "net/http"
    "os"
    "os/signal"
    "syscall"
    "flag"
)

// main risponde del avvio del'applicazione e della sua
// registrazione come server in ascolto su la rete.
func main() {

    var confFilePath = flag.String("conf", "./conf.ini", "path to configuration file")
    var logLevel = flag.String("log", "DEBUG", "output log level NONE, INFO, MESSAGE, ERROR, DEBUG")

    flag.Parse()

    // livello generale del log, quantita dei messaggi da stampare
    log.SetLevel(*logLevel)

    err := core.ReadConfiguration(*confFilePath)
    if err != nil {
        log.Info("no valid configuration, details: %v", err)
        return
    }

    // istruiamo la database di creare una nuova connessione.
    // specificandoli a quale database si deve collegare
    // TODO: qunado il database viene creato per la prima volta, ci dobbiamo
    // assigurare che abbiamo descrito il campo usernaeme come un campo con dei valori
    // unici.
    err = database.NewConnection("mapo")
    if err != nil {
        log.Info("error connecting to database (%v)", err)
        return
    }
    log.Msg("created a new database connection")

    // al avvio del'applicazione si verifica la disponibilità dei addon
    // e si crea una lista globale che sarà passata verso altri moduli
    // TODO: modulo addon ancora da implementare
    addons := addon.GetAll()
    addons = addons
    log.Msg("load addons and generate a list")

    // al momento del spegnimento del'applicazione potremo trovarci con delle
    // connessione attive dal parte del cliente. Il handler personalizzato usato
    // qui, ci permette di dire al server di spegnersi ma prima deve aspettare
    // che tutte le richieste siano processate e la connessione chiusa.
    //
    // Oltre al spegnimento sicuro il ServeMux permette di registra dei nuovi
    // handler usando come descrizione anche il metodo http tipo GET o POST.
    muxer := NewServeMux()

    // qui si assegna al muxer la funzione che sara' usata per l'autenticazione
    muxer.SetAuthenticator(core.Authenticator)

    server := &http.Server {
        Addr:   ":8081",
        Handler: muxer,
    }

    // TODO: register this node to load-balancing service

    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT)

    // aviamo in una nuova gorutine la funzione che ascoltera per il segnale di
    // spegnimento del server
    go muxer.getSignalAndClose(c)

    muxer.HandleFuncNoAuth("POST", "/admin/user", core.NewUser)
    muxer.HandleFunc("GET", "/admin/user/{id}", core.GetUser)
    muxer.HandleFunc("GET", "/admin/user", core.GetUserAll)
    muxer.HandleFunc("POST", "/admin/user/{id}", core.UpdateUser)

    muxer.HandleFunc("POST", "/admin/studio", core.NewStudio)
    muxer.HandleFunc("GET", "/admin/studio", core.GetStudioAll)
    muxer.HandleFunc("GET", "/admin/studio/{id}", core.GetStudio)

    muxer.HandleFunc("POST", "/admin/project", core.NewProject)
//    muxer.HandleFunc("GET", "/admin/project", core.GetProjectAll)
//    muxer.HandleFunc("GET", "/admin/project/{id}", core.GetProject)

    muxer.HandleFuncNoAuth("GET", "/", webui.Root)

    muxer.HandleFuncNoAuth("POST", "/login", core.Login)
    muxer.HandleFuncNoAuth("GET", "/logout", core.Logout)

    jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir("/home/develop/go/src/mapo/webui/static/js")))
    muxer.Handle("GET", "/js/.*\\.js", jsHandler)

    cssHandler := http.StripPrefix("/css/", http.FileServer(http.Dir("/home/develop/go/src/mapo/webui/static/css")))
    muxer.Handle("GET", "/css/.*\\.css", cssHandler)

    // OAuth
    //oauth2callback
    muxer.HandleFuncNoAuth("GET", "/oauth2callback", core.OAuthCallBack)

    log.Info("start listening for requests")

    // avviamo il server che processerà le richieste
    log.Msg("close server with message: %v", server.ListenAndServe())
}


