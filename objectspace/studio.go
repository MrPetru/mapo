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

package objectspace

import (
    "mapo/db"
    "github.com/maponet/utils/log"

    "errors"
    "time"
    "labix.org/v2/mgo/bson"
)

type studio struct {
    Id string `bson:"_id"`
    Name string
    Description string
    Owners []string
    Projects []string

    Created time.Time
}

func NewStudio() studio {
    s := new(studio)
    s.Owners = make([]string, 0)
    s.Projects = make([]string, 0)

    return *s
}

func (s *studio) SetId(value string) error {
    if len(value) < 3 {
        return errors.New("troppo corto")
    }
    s.Id = value
    return nil
}

func (s *studio) GetId() string {
    return s.Id
}

func (s *studio) SetName(value string) error {
    if len(value) < 6 {
        return errors.New("troppo corto")
    }

    s.Name = value
    return nil
}

func (s *studio) SetDescription(value string) error {
    s.Description = value

    return nil
}

func (s *studio) AppendOwner(value string) error {
    if len(value) != 32 {
        return errors.New("troppo corto")
    }

    s.Owners = append(s.Owners, value)
    return nil
}

func (s *studio) AppendProject(value string) error {
    if len(value) != 32 {
        return errors.New("troppo corto")
    }

    s.Projects = append(s.Projects, value)
    return nil
}

func (s *studio) Save() error {
    log.Debug("save studio to database")
    err := db.Store(s, "studios")
    return err
}

// Restore interoga il database per le informazioni di un certo studio
func (s *studio) Restore() error {
    log.Debug("restoring user from database")

    err := db.RestoreOne(s, bson.M{"_id":s.Id}, "studios")

    return err
}

func StudioRestoreAll(filter bson.M) ([]studio, error) {
    studioList := make([]studio,0)

    err := db.RestoreList(&studioList, filter, "studios")

    return studioList, err
}

func (s *studio) Update() error {
    log.Debug("update studio to database")
    err := db.Update(s, s.Id, "studios")
    return err
}
