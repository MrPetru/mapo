package objectspace

import (
    "mapo/log"
    "mapo/database"

    "errors"
)

type project struct {
    Id string `bson:"_id"`
    Name string
    Studio string
//    Owners []string
    Admins []string
}

func NewProject() project {
    p := new(project)
    //p.Owners = make([]string, 0)
    p.Admins = make([]string, 0)

    return *p
}

func (p *project) SetName(value string) error {
    if len(value) > 6 {
        p.Name = value
        return nil
    }

    return errors.New("nome progetto tropo corto")
}

func (p *project) SetStudio(value string) error {
    if len(value) > 6 {
        p.Studio = value
        return nil
    }

    return errors.New("nome studio tropo corto")
}

func (p *project) SetId(value string) error {

    if len(value) < 24 {
        return errors.New("troppo corto")
    }
    p.Id = value
    return nil
}

func (p *project) Save() error {
    log.Debug("save project to database")
    err := database.Store(p, "projects")
    return err
}
