/*
DESCRIZIONE DI MAPO
*/
package main

import (
    "mapo/database"
    "mapo/addon"
    "mapo/log"
    "mapo/core"
    "mapo/webui"
    "mapo/api"

    "net/http"
    "os"
    "os/signal"
    "syscall"
    "flag"
)

// main risponde del avvio dell'applicazione e della sua
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

    /*
    in questa configurazione, connessione alla database viene attivata in un
    oggetto definito globalmente al interno del modulo database.
    L'idea originale per Mapo è di creare un oggetto che contenga la
    connessione attiva e passare questo aggetto a tutte le funzione che ne
    hanno bisogno di fare una richiesta alla database.

    Passare l'oggetto database da una funzione ad altra, potrebbe
    significare, creare una catena dalla prima funzione all'ultima. Che
    avvolte non fa niente altro che aumentare il numero dei parametri passati
    da una funzione ad altra. Per esempio, la connessione al database si usa
    nel modulo objectspace che viene chiamato dal modulo core che al suo tempo
    viene chiamato da main. Inutile passare questo oggetto al modulo core,
    visto che li lui non serve.

    NOTA: accesso ai oggetti globali deve essere in qualche modo sincronizzato
    per evitare i problemi di inconsistenza.

    NOTA: le osservazioni dimostrano che avendo una connessione attiva alla
    database che poi viene riutilizzata, diminuisce considerevolmente i tempi di
    interrogazione.
    */
    err = database.NewConnection("mapo")
    if err != nil {
        log.Info("error connecting to database (%v)", err)
        return
    }
    log.Msg("created a new database connection")

    // al avvio del'applicazione si verifica la disponibilità dei addon
    // e si crea una lista globale che sarà passata verso altri moduli
    // TODO: modulo addon ancora da implementare

    /*
    anche qui il discorso è molto simile a quello della connessione alla
    database.
    Passare l'oggetto addons nella catena per arrivare al punto di destinazione
    potrebbe creare dei disagi.
    */
    addons := addon.GetAll()
    addons = addons
    log.Msg("load addons and generate a list")

    // al momento del spegnimento dell'applicazione potremo trovarci con delle
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

    // aviamo in una nuova gorutine la funzione che ascolterà per il segnale di
    // spegnimento del server
    go muxer.getSignalAndClose(c)

    muxer.HandleFunc("GET", "/admin/user/{uid}", core.Authenticate(core.GetUser))

    muxer.HandleFunc("POST", "/admin/studio", core.Authenticate(core.NewStudio))
    muxer.HandleFunc("GET", "/admin/studio", core.Authenticate(core.GetStudioAll))
    muxer.HandleFunc("GET", "/admin/studio/{sid}", core.Authenticate(core.GetStudio))
    muxer.HandleFunc("GET", "/admin/studio/{sid}/update", core.Authenticate(core.UpdateStudio))

    muxer.HandleFunc("POST", "/admin/project", core.Authenticate(core.NewProject))
    muxer.HandleFunc("GET", "/admin/project", core.Authenticate(core.GetProjectAll))
    muxer.HandleFunc("GET", "/admin/project/{pid}", core.Authenticate(core.GetProject))

    muxer.HandleFunc("GET", "/api/{pid}", core.Authenticate(core.GetProject))
    muxer.HandleFunc("GET", "/api/{pid}/.*", core.Authenticate(api.HttpWrapper))

    muxer.HandleFunc("GET", "/", webui.Root)

    muxer.HandleFunc("GET", "/login/{oauthprovider}", core.Login)
    //muxer.HandleFunc("GET", "/logout", core.Logout)

    jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir("/home/develop/go/src/mapo/webui/static/js")))
    muxer.Handle("GET", "/js/.*\\.js", jsHandler)

    cssHandler := http.StripPrefix("/css/", http.FileServer(http.Dir("/home/develop/go/src/mapo/webui/static/css")))
    muxer.Handle("GET", "/css/.*\\.css", cssHandler)

    icoHandler := http.StripPrefix("/", http.FileServer(http.Dir("/home/develop/go/src/mapo/webui/static/image")))
    muxer.Handle("GET", "/favicon\\.ico", icoHandler)

    // OAuth
    // su questo url viene reinderizato il cliente dopo che la procedura di authenticazione
    // sul server del servizio aviene con successo o meno.
    muxer.HandleFunc("GET", "/oauth2callback", core.OAuthCallBack)

    log.Info("start listening for requests")

    // avviamo il server che processerà le richieste
    log.Msg("close server with message: %v", server.ListenAndServe())
}

