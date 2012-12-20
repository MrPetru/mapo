package objectspace

import (
    "mapo/log"
    "mapo/database"
    
    "errors"
)

// il contenitore base che si usa per transportare i dati di un utente verso
// il database e dat database.
// Accesso a questo contenitore avviene attraverso le funzioni definiti qui.
type user struct {
    Id string `bson:"_id"`
    Login string
    Name string
    Password string `json:"-"`
    Contacts []string
    Description string
    Rating int
    Studios []string
}

// una lista di utenti
type userList []user

func NewUserList() userList {
    ul := make(userList, 0)
    
    return ul
}

func (ul *userList) Restore() error {
    log.Debug("restore all users from database")
    
    err := database.RestoreList(ul, "users")
    
    return err
}

func NewUser() user {
    u := new(user)
    u.Contacts = make([]string,0)
    u.Rating = 0
    u.Studios = make([]string,0)
    
    return *u
}

func (u *user) SetId(value string) error {
    
    if len(value) < 24 {
        return errors.New("troppo corto")
    }
    u.Id = value
    return nil
}

func (u *user) GetId() string {
    return u.Id
}

func (u *user) SetLogin(value string) error {

    if len(value) < 4 {
        return errors.New("troppo corto")
    }
    
    u.Login = value
    return nil
}

func (u *user) GetLogin() string {

    return u.Login
}

func (u *user) SetPassword(value string) error {
    
    if len(value) < 6 {
        return errors.New("troppo corta") 
    }
    
    u.Password = value
    return nil
}

func (u *user) SetName(value string) error {
    
    if len(value) < 6 {
        return errors.New("troppo corto") 
    }
    
    u.Name = value
    return nil
}

func (u *user) SetRating(value int) error {
    if value > 100.0 || value < 0 {
        return errors.New("value out of range")
    }
    
    u.Rating = value
    return nil
}

func (u *user) AppendStudioId(value string) error {
    if len(value) < 4 {
        return errors.New("tropo corto")
    }
    
    u.Studios = append(u.Studios, value)
    return nil
}

// Reastore interoga il database per le informazioni di un certo utente
func (u *user) Restore() error {
    log.Debug("restoring user from database")
    
    err := database.RestoreOne(u, u.Id, "users")
    
    return err
}

// Save salva i dati contenuti nel contenitore user nella database
func (u *user) Save() error {
    log.Debug("save user to database")
    err := database.Store(u, "users")
    return err
}

func (u *user) Update() error {
    log.Debug("update user to database")
    err := database.Update(u, u.Id, "users")
    return err
}

