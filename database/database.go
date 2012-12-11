package database

import (
    "mapo/log"
    
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "errors"
)

var dbSession *mgo.Session

// TODO: c'è bisogno di una funzione che faccia il lavoro inverso. trasforma da un
// formatto adatto per database in un formatto di lavoro per mapo.
type Storer interface {
    //ToStoreFormat() map[string]interface{}
    ToMap() map[string]interface{}
    FillWithResult(map[string]interface{})
}

// TODO: definire una funzione che si ocupa con la rezione e gestione della
// connessione verso un database.
func NewConnection() {
    log.Debug("executing NewConnection function")
    
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    
    dbSession = session
    
    // connessione alla data base avvenne usando diversi livelli di autenticazione
    // come admin, user, ... e probabile altri?
}

func Store(inData Storer) error {
    // store a object to database
    object := inData.ToMap()
    
    object["_id"] = object["id"]
    delete(object, "id")
    
    id, _ := object["_id"]
    idStr := id.(string)
    
    if len(idStr) < 1 {
        return errors.New("no id was submited")
    }
    
    c := dbSession.DB("mapo").C("users")
    
    err := c.Insert(&object)
    if err != nil {
            return err
    }
    
    log.Debug("stored %v", object )
    return nil
}

func Restore(inData Storer) error {
    // interroga la database per un oggetto che è descritto nel inData
    object := inData.ToMap()

    object["_id"] = object["id"]
    
    c := dbSession.DB("mapo").C("users")
    
    err := c.Find(bson.M{"_id": object["_id"]}).One(&object)
    if err != nil {
        return err
    }
    
    object["id"] = object["_id"]
    delete(object, "_id")
    
    inData.FillWithResult(object)
    
    log.Debug("restored %v", object )
    return nil
}

func Update(inData Storer) {
    // store a object to database
    object := inData.ToMap()
    
    log.Debug("updated %v", object )
}
