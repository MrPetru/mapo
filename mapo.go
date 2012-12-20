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
)

// main risponde del avvio del'applicazione e della sua
// registrazione come server in ascolto su la rete.
func main() {

    // settiamo il livello generale dei messaggi da visualizzare
    log.SetLevel("DEBUG")
    
    // istruiamo la database di creare una nuova connessione.
    // specificandoli a quale database si deve collegare
    database.NewConnection("mapo")
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

    muxer.HandleFunc("POST", "/admin/user", core.NewUser)
    muxer.HandleFunc("GET", "/admin/user/{id}", core.GetUser)
    muxer.HandleFunc("GET", "/admin/user", core.GetUserAll)
    muxer.HandleFunc("POST", "/admin/user/{id}", core.UpdateUser)
    
    muxer.HandleFunc("POST", "/admin/studio", core.NewStudio)
    
    log.Info("start listening for requests")
    
    // avviamo il server che processerà le richieste
    log.Msg("close server with message: %v", server.ListenAndServe())
}


