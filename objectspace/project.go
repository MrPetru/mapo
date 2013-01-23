package objectspace

import (
    "mapo/log"
    "mapo/database"

    "errors"
    "labix.org/v2/mgo/bson"
)

type project struct {
    Id string `bson:"_id"`
    Name string
    Description string
    StudioId string
    Admins []string
    Supervisors []string
    Artists []string

    Addons []string `json:"-"`
}

func NewProject() project {
    p := new(project)
    p.Admins = make([]string, 0)
    p.Supervisors = make([]string, 0)
    p.Artists = make([]string, 0)

    return *p
}

func (p *project) SetName(value string) error {
    if len(value) > 6 {
        p.Name = value
        return nil
    }

    return errors.New("nome progetto tropo corto")
}

func (p *project) SetDescription(value string) error {
    p.Description = value

    return nil
}

func (p *project) SetId(value string) error {

    if len(value) < 32 {
        return errors.New("troppo corto")
    }
    p.Id = value
    return nil
}

func (p *project) SetStudioId(value string) error {
    p.StudioId = value

    return nil
}

func (p *project) Save() error {
    log.Debug("save project to database")
    err := database.Store(p, "projects")
    return err
}

func ProjectRestorList(filter bson.M) ([]project, error) {
    p := make([]project, 0)

    err := database.RestoreList(&p, filter, "projects")

    if err != nil {
        return nil, err
    }

    return p, nil
}

func (p *project) Restore() error {

    err := database.RestoreOne(&p, bson.M{"_id":p.Id}, "projects")

    if err != nil {
        return err
    }

    return nil
}

// ritorna dal database la lista dei addon attivi per il progetto
func (p *project) GetAddonList() []string {

    return p.Addons
}
