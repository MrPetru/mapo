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

/*
Package db contains a data abstraction layer and underlying facilities
to store entities in a database.
*/
package db

import (
    "github.com/maponet/utils/log"

    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
)

// un oggetto globale che contiene una connessione attiva con la database.
var database *mgo.Database

// TODO: definire una funzione che si occupa con la creazione e gestione della
// connessione verso un database.
func NewConnection(databaseName string) error {
    log.Info("executing NewConnection function")

    session, err := mgo.Dial("localhost")
    if err != nil {
        return err
    }

    database = session.DB(databaseName)
    return nil
    // connessione alla data base avvenne usando diversi livelli di autenticazione
}

// Store salva nella database un singolo oggetto
func Store(data interface{}, table string) error {

    c := database.C(table)

    err := c.Insert(data)

    return err
}

// RestoreOne riprende dalla database un singolo oggetto identificato da un id
func RestoreOne(data interface{}, filter bson.M, table string) error {

    c := database.C(table)

    err := c.Find(filter).One(data)

    return err
}

// RestoreList riprende dalla database una lista (tutti) di oggetti, senza alcun filtro
func RestoreList(data interface{}, filter bson.M, table string) error {

    c := database.C(table)

    err := c.Find(filter).All(data)

    return err
}

// Update aggiorna i valori di un oggetto nella database, identificato da un id
func Update(data interface{}, id string, table string) error {

    c := database.C(table)

    err := c.Update(bson.M{"_id": id}, data)

    return err
}
