package objectspace

import (
    "mapo/log"
    "mapo/database"
    
    "errors"
    "fmt"
    "strconv"
//    "reflect"
)

// il contenitore base che si usa per transportare i dati di un utente verso
// il database e dat database.
// Accesso a questo contenitore avviene attraverso le funzioni definiti qui.
type user struct {
    id string
    login string
    name string
    password string
    contacts []string
    description string
    rating float32
    studios []string
}

// una lista di utenti
type userList []user

func (ul *userList) Restore() error {
    log.Debug("restored all users from database")
    
    err := database.RestoreList(ul)
    
    return err
}

func (ul userList) ToMap() []map[string]interface{} {
    
    return nil
}

func (ul *userList) FillWithResult(result []map[string]interface{}) {
    for _, v := range(result) {
        user := NewUser()
        user.FillWithResult(v)
        *ul = append(*ul, user)
    }
}

func NewUser() user {
    u := new(user)
    u.contacts = make([]string,0)
    u.rating = 0
    u.studios = make([]string,0)
    
    return *u
}

func NewUserList() userList {
    ul := make(userList, 0)
    
    return ul
}

func (u *user) SetId(value string) error {
    
    if len(value) < 4 {
        return errors.New("id: troppo corto")
    }
    u.id = value
    return nil
}

func (u *user) GetId() string {
    return u.id
}

func (u *user) SetLogin(value string) error {

    if len(value) < 4 {
        return errors.New("login: troppo corto")
    }
    
    u.login = value
    return nil
}

func (u *user) GetLogin() string {

    return u.login
}

func (u *user) SetPassword(value string) error {
    
    if len(value) < 6 {
        return errors.New("password: troppo corta") 
    }
    
    u.password = value
    return nil
}

func (u *user) SetName(value string) error {
    
    if len(value) < 6 {
        return errors.New("nome: troppo corto") 
    }
    
    u.name = value
    return nil
}

// Reastore interoga il database per le informazioni di un certo utente
func (u *user) Restore() error {
    log.Debug("restoring user from database")
    
    err := database.RestoreOne(u)
    
    return err
}

// Save salva i dati contenuti nel contenitore user nella database
func (u *user) Save() error {
    log.Debug("save user to database")
    err := database.Store(u)
    return err
}

func (u *user) SaveUpdate() {
    log.Debug("save user to database")
    database.Update(u)
}

// ToMap, trasforma il contenitore user in una ogetto di tipo mapo. Questa
// operazione permette di omogenizzare i dati restituiti dal pacchetto database
// ai pacchetti esterni.
func (u user) ToMap() map[string]interface{} {
    log.Msg("translate user struct to a map[] object")
    m := make(map[string]interface{})
    
    m["id"] = u.id
    m["login"] = u.login
    m["name"] = u.name
    m["password"] = u.password
    m["description"] = u.description
    m["contacts"] = u.contacts
    m["studios"] = u.studios
    return m
}

func (u *user) FillWithResult(result map[string]interface{}) {
    //
    
    tmp := make([]string, 0)
    
    if _, ok := result["id"]; ok {
        u.id = result["id"].(string)
    } else {
        u.id = result["_id"].(string)
    }
    
    u.login = result["login"].(string)
    u.name = result["name"].(string)

    for _, v := range(result["contacts"].([]interface{})) {
        tmp = append(tmp, fmt.Sprintf("%v", v))
    }
    u.contacts = tmp

    u.description = result["description"].(string)
    vfloat, _ := strconv.ParseFloat(fmt.Sprintf("%v",result["rating"]), 32)
    u.rating = float32(vfloat)

    tmp = make([]string, 0)
    for _, v := range(result["studios"].([]interface{})) {
        tmp = append(tmp, fmt.Sprintf("%v", v))
    }
    u.studios = tmp
}

