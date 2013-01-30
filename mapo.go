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
	"mapo/log"
    "mapo/admin"
    "mapo/db"
    "mapo/addons"

    "net/http"
    "os"
    "os/signal"
    "syscall"
    "flag"
)

func main() {

	log.Info("Starting application")

	// parse flags
    var logLevel = flag.Int("log", 1, "set message level eg: 0 = DEBUG, 1 = INFO, 2 = ERROR")
    var confFilePath = flag.String("conf", "./conf.ini", "set path to configuration file")
    flag.Parse()

    // set log level
	log.SetLevel(*logLevel)
	log.Info("Setting log level to %d", *logLevel)

	// load config and setup application
	log.Info("Loading configuration from file")
    err := admin.ReadConfiguration(*confFilePath)
    if err != nil {
        log.Info("%s, no such file or directory", *confFilePath)
        return
    }

	// setup application

	// init db
	log.Info("Initializing db")
    /*
    in questa configurazione, connessione alla database viene attivata in un
    oggetto definito globalmente al interno del modulo db.
    L'idea originale per Mapo è di creare un oggetto che contenga la
    connessione attiva e passare questo aggetto a tutte le funzione che ne
    hanno bisogno di fare una richiesta alla database.

    Passare l'oggetto database da una funzione ad altra, potrebbe
    significare, creare una catena dalla prima funzione all'ultima. Che
    avvolte non fa niente altro che aumentare il numero dei parametri passati
    da una funzione ad altra. Per esempio, la connessione al database si usa
    nel modulo objectspace che viene chiamato dal modulo admin che al suo tempo
    viene chiamato da main. Inutile passare questo oggetto al modulo admin,
    visto che li lui non serve.

    NOTA: accesso ai oggetti globali deve essere in qualche modo sincronizzato
    per evitare i problemi di inconsistenza.

    NOTA: le osservazioni dimostrano che avendo una connessione attiva alla
    database che poi viene riutilizzata, diminuisce considerevolmente i tempi di
    interrogazione.
    */
    err = db.NewConnection("mapo")
    if err != nil {
        log.Error("%v", err)
        return
    }


	// load addons
	log.Info("Loading addons")
    /*
    anche qui il discorso è molto simile a quello della connessione alla
    database.
    Passare l'oggetto addons nella catena per arrivare al punto di destinazione
    potrebbe creare dei disagi.
    */
    addonList := addons.GetAll()
    addonList = addonList
    log.Info("load addons and generate a list")

    // al momento del spegnimento dell'applicazione potremo trovarci con delle
    // connessione attive dal parte del cliente. Il handler personalizzato usato
    // qui, ci permette di dire al server di spegnersi ma prima deve aspettare
    // che tutte le richieste siano processate e la connessione chiusa.
    //
    // Oltre al spegnimento sicuro il ServeMux permette di registra dei nuovi
    // handler usando come descrizione anche il metodo http tipo GET o POST.
    muxer := NewServeMux()

    // prepare server
    server := &http.Server {
        Addr:   ":8081",
        Handler: muxer,
    }

    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT)

    // aviamo in una nuova gorutine la funzione che ascolterà per il segnale di
    // spegnimento del server
    go muxer.getSignalAndClose(c)

	// register handlers
	log.Info("Registering handlers")

    // register with supervisor
	log.Info("Joining supervisor")

	// start server
	log.Info("Listening for requests")
    log.Info("close server with message: %v", server.ListenAndServe())

	// inform supervisor that we are up

	// for each request
		// check authentication/authorization

		// extract request operation

		// extract request arguments

		// pass operation and arguments to api.router

			// find function mapped to operation

			// call function with arguments

		// return result to user


	// close on signal
	log.Info("Closing application")
}
