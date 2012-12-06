/*
database contiene la definizione dei oggetti che poi verranno usati per scambio
con il database. Probabilmente questo non è il modulo adatto ma al momento è
il pacchetto che si assomiglia di più alla caratteristiche del pacchetto interno
definito nel progetto example.
TODO: non è del tutto chiaro che l'interazione con questo modulo avviene tramite
le interfacce.
*/
package database

import (
    "mapo/log"
)

// TODO: definire una funzione che si ocupa con la rezione e gestione della
// connessione verso un database.
func NewConnection() {
    log.Debug("executing NewConnection function")
}
