package objectspace

import (
    "mapo/database"
    "mapo/log"

    "errors"
    "labix.org/v2/mgo/bson"
)

type studio struct {
    Id string `bson:"_id"`
    Name string
    Owners []string
    Projects []string
}

func NewStudio() studio {
    s := new(studio)
    s.Owners = make([]string, 0)
    s.Projects = make([]string, 0)

    return *s
}

func (s *studio) SetId(value string) error {
    if len(value) < 4 {
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

func (s *studio) AppendOwner(value string) error {
    if len(value) != 24 {
        return errors.New("troppo corto")
    }

    s.Owners = append(s.Owners, value)
    return nil
}

func (s *studio) AppendProject(value string) error {
    if len(value) != 24 {
        return errors.New("troppo corto")
    }

    s.Projects = append(s.Projects, value)
    return nil
}

func (s *studio) Save() error {
    log.Debug("save studio to database")
    err := database.Store(s, "studios")
    return err
}

// Restore interoga il database per le informazioni di un certo studio
func (s *studio) Restore(filter bson.M) error {
    log.Debug("restoring user from database")

    err := database.RestoreOne(s, filter, "studios")

    return err
}

func (s *studio) Update() error {
    log.Debug("update studio to database")
    err := database.Update(s, s.Id, "studios")
    return err
}