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

package admin

import (
    "github.com/maponet/utils/log"
    "mapo/db"

    "errors"
    "time"
    "labix.org/v2/mgo/bson"
)

type user struct {
    Id string `bson:"_id"`
    Name string
    Email string
    Oauthid string
    Oauthprovider string
    Avatar string `json:"picture"`

    Registered time.Time `json:"-"`
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
    err := db.RestoreOne(u, bson.M{"_id":id}, "users")
    return err
}

func (u *user) Save() error {
    log.Debug("save user to database")
    err := db.Store(u, "users")
    return err
}

func (u * user) Update() error {
    log.Debug("update user to database")
    err := db.Update(u, u.Id, "users")
    return err
}
