package database

import (
    "mapo/log"
    
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "errors"
    
    "fmt"
)

var dbSession *mgo.Session

// TODO: c'è bisogno di una funzione che faccia il lavoro inverso. trasforma da un
// formatto adatto per database in un formatto di lavoro per mapo.
type Storer interface {
    //ToStoreFormat() map[string]interface{}
    ToMap() map[string]interface{}
    FillWithResult(map[string]interface{})
}

type StorerList interface {
    //ToStoreFormat() map[string]interface{}
    ToMap() []map[string]interface{}
    FillWithResult([]map[string]interface{})
}

// TODO: definire una funzione che si occupa con la creazione e gestione della
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

func RestoreOne(inData Storer) error {
    // interroga la database per un oggetto che è descritto nel inData
    
    c := dbSession.DB("mapo").C("users")

    fmt.Printf("solo un elemento\n")
    
    object := inData.ToMap()

    object["_id"] = object["id"]
    
    err := c.Find(bson.M{"_id": object["_id"]}).One(object)
    if err != nil {
        return err
    }
    
    object["id"] = object["_id"]
    delete(object, "_id")
    
    inData.FillWithResult(object)
    
    log.Debug("restored %v", object )

    return nil
}

func RestoreList(inData StorerList) error {
    // interroga la database per un oggetto che è descritto nel inData
    
    c := dbSession.DB("mapo").C("users")

    fmt.Printf("tanti elementi\n")
    
    object := make([]map[string]interface{}, 0)
    
    err := c.Find(bson.M{}).All(&object)
    if err != nil {
        return err
    }
    
    inData.FillWithResult(object)

    return nil
}

func Update(inData Storer) {
    // store a object to database
    object := inData.ToMap()
    
    log.Debug("updated %v", object )
}
