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
    "github.com/maponet/utils/log"
    "mapo/db"

    "errors"
    "time"
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

    Created time.Time
    Addons map[string][]string `json:"-"`
}

func NewProject() project {
    p := new(project)
    p.Admins = make([]string, 0)
    p.Supervisors = make([]string, 0)
    p.Artists = make([]string, 0)
	p.Addons = make(map[string][]string, 0)

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
    err := db.Store(p, "projects")
    return err
}

func (p *project) Update() error {
    log.Debug("update project to database")
    err := db.Update(p, p.Id, "projects")
    return err
}

func ProjectRestorList(filter bson.M) ([]project, error) {
    p := make([]project, 0)

    err := db.RestoreList(&p, filter, "projects")

    if err != nil {
        return nil, err
    }

    return p, nil
}

func (p *project) Restore() error {

    err := db.RestoreOne(&p, bson.M{"_id":p.Id}, "projects")

    if err != nil {
        return err
    }

    return nil
}

func (p *project) AddAddon(entity, addonId string) {
	if p.Addons == nil {
		p.Addons = make(map[string][]string)
	}
	ent, ok := p.Addons[entity]
	if !ok {
		return
	}
	for _, aId := range(ent) {
		if aId == addonId {
			return
		}
	}
	ent = append(ent, addonId)
	p.Addons[entity] = ent
}

// ritorna dal database la lista dei addon attivi per il progetto
func (p *project) GetAddonList(entity string) []string {

    return p.Addons[entity]
}
