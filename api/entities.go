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

package api

import (
	"github.com/maponet/utils/log"
	"mapo/db"
	"labix.org/v2/mgo/bson"

	"regexp"
	"strings"
)

type Entity struct {
	name string
	projectId string
	//addonsId []string
	Function map[string] func(*EntityContainer, Data) interface{}
	attributes map[string]attribute
}

type attribute struct {
	value string
	t int
}

func (e *Entity) AddAttribute(name string, t int) {
	a := new(attribute)
	a.value = ""
	a.t = t

	e.attributes[name] = *a
}

func (e *Entity) SetAttribute(name, value string) {
	attr, ok := e.attributes[name]
	if ok {
		attr.value = value
		e.attributes[name] = attr
	}
}

func (e *Entity) AddFunction(method, path string, f func(*EntityContainer, Data) interface{}) {
	if e.Function == nil {
		e.Function = make(map[string] func(*EntityContainer, Data) interface{})
	}
	pattern := createPattern(method, path)
	e.Function[pattern] = f
}

func (e *Entity) RunByPath(method, path string, entities *EntityContainer, data Data) interface{} {
	var f func(*EntityContainer, Data) interface{}
	for k, v := range(e.Function) {
		matching, _ := regexp.MatchString(k, method + ":" + path)
        if matching {
            f = v
            break
        }

	}

	if f != nil {
		result := f(entities, data)
		return result
	}

	return nil
}

func (e *Entity) ToMap() map[string]interface{} {
	m := make(map[string]interface{},0)
	for name, attr := range(e.attributes) {
		m[name] = attr.value
	}
	log.Debug("entity as map %v\n", m)

	return m
}

func (e *Entity) Restore(pid, id string) error {
	collection := "foraddon_"+pid+"_"+e.name
	//m := e.ToMap()
	m := make(map[string]interface{})
	m["_id"] = ""
	delete(m, "id")

	filter := bson.M{"_id":bson.ObjectIdHex(id)}
	err := db.RestoreOne(m, filter, collection)
	if err != nil {
		return err
	}

	m["id"] = ""

	for k, v := range(m) {
		if k == "_id"  || k == "id" {
			e.SetAttribute("id", id)
			continue
		}
		e.SetAttribute(k, v.(string))
	}
	delete(m, "_id")
	return nil
}

func (e *Entity) Store(pid string) (string, error) {
	collection := "foraddon_"+pid+"_"+e.name
	m := e.ToMap()
	id := bson.NewObjectId()
	//m["_id"] = md5sum()
	m["_id"] = id
	delete(m, "id")
	err := db.Store(m, collection)
	delete(m, "_id")

	if err == nil {
		return id.Hex(), nil
	}
	return "", err
}

type EntityList struct {
	name string
	baseEntity *Entity
	entities []Entity
}

func (el *EntityList) Restore(pid string) error {
	collection := "foraddon_"+pid+"_"+el.name
	ml := make([]map[string]interface{},0)
	err := db.RestoreList(&ml, bson.M{}, collection)
	if err != nil {
		return err
	}
	log.Debug("entities from database %v\n", ml)

	for _, entry := range(ml) {
		delete(entry, "id")
		ent := new(Entity)
		ent.attributes = make(map[string]attribute, 0)
		for ak, _ := range(el.baseEntity.attributes) {
			a := new(attribute)
			ent.attributes[ak] = *a
		}
		ent.name = el.name
		//ent.Function = make(map[string] func(*EntityContainer, Data) interface{})

		for k, v := range(entry) {
			if k == "_id" {
				id := v.(bson.ObjectId)
				ent.SetAttribute("id", id.Hex())
				continue
			}
			ent.SetAttribute(k, v.(string))
		}
		el.entities = append(el.entities, *ent)
	}

	log.Debug("entity list = %v\n", el)

	return nil
}

func (el *EntityList) ToMap() []map[string]interface{} {
	list := make([]map[string]interface{},0)
	for _, e := range(el.entities) {
		m := make(map[string]interface{},0)
		for name, attr := range(e.attributes) {
			m[name] = attr.value
		}
		list = append(list, m)
	}

	log.Debug("entity list as map %v\n", list)

	return list
}

type EntityContainer map[string]*Entity

func NewEntitiesList() *EntityContainer {
	el := make(EntityContainer, 0)

	return &el
}

func (es *EntityContainer) New(name string) *Entity {
	e :=new(Entity)
	e.name = name
	e.Function = make(map[string] func(*EntityContainer, Data) interface{})
	e.attributes = make(map[string]attribute)
	(*es)[name] = e
	return e
}

func (es *EntityContainer) GetEntity(name string) Entity {
	e := (*es)[name]
	return *e
}

func (es *EntityContainer) GetEntityList(name string) EntityList {
	e := (*es)[name]
	eList := new(EntityList)
	eList.name = name
	eList.baseEntity = e
	eList.entities = make([]Entity, 0)
	return *eList
}

func createPattern(method, path string) string {
    pattern := "(?i)^"

    if method != "" {
        pattern = pattern + method + ":/"
    } else {
        pattern = pattern + "(GET|POST)" + ":/"
    }

    if len(path) > 1 {
        pathSlice := strings.Split(path[1:], "/")
        for _, v := range(pathSlice) {
            if v[0] == '{' {
                pattern = pattern + "[0-9a-z_\\ \\.\\+\\-]*/"
            } else {
                pattern = pattern + v + "/"
            }
        }
    }
    pattern = pattern + "$"
    return pattern
}
