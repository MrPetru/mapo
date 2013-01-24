package objectspace

import (
    "mapo/log"
    "mapo/database"

    "errors"
    "labix.org/v2/mgo/bson"
)

type user struct {
    Id string `bson:"_id"`
    Name string
    Email string
    Oauthid string
    Oauthprovider string
    Avatar string `json:"picture"`

    AccessToken string `json:"-"`
}

func NewUser() user {
    u := new(user)
    return *u
}

func (u *user) CreateId() {
    u.Id = Md5sum(u.Oauthprovider + u.Oauthid)
}

func (u *user) SetId(id string) error {
    if len(id) != 32 {
        return errors.New("invalid user id")
    }

    u.Id = id
    return nil
}

func (u *user) GetId() string {
    return u.Id
}

func (u *user) Restore() error {
    id := u.Id
    err := database.RestoreOne(u, bson.M{"_id":id}, "users")
    return err
}

func (u *user) Save() error {
    log.Debug("save user to database")
    err := database.Store(u, "users")
    return err
}

func (u * user) Update() error {
    log.Debug("update user to database")
    err := database.Update(u, u.Id, "users")
    return err
}
