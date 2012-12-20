package objectspace

import (
    "mapo/database"
    "mapo/log"
    
    "errors"
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

func (s *studio) SetUserid(value string) error {
    if len(value) != 24 {
        return errors.New("troppo corto")
    }
    
    s.Owners = append(s.Owners, value)
    return nil
}

func (s *studio) Save() error {
    log.Debug("save studio to database")
    err := database.Store(s, "studios")
    return err
}
