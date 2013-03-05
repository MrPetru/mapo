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
	"errors"
)

import (
	"github.com/maponet/utils/log"
	"labix.org/v2/mgo/bson"
	"mapo/db"

	"mapo/addons"
	"reflect"
	"regexp"
	"strings"
)

func newEntity(entity addons.CompEntity, data addons.RequestData) (addons.CompEntity, error) {
	var err error

	compE, ok := entity.(*composedEntity)
	if !ok {
		return nil, errors.New("not a known type")
	}
	localEntity := compE.s

	for key, _ := range localEntity.attributes {
		attrValue := data.GetValue(key)
		if len(attrValue) > 0 {
			entity.SetAttribute(key, attrValue)
			// TODO: gestire le validazioni
		}
	}

	err = entity.Store()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func getOne(entity addons.CompEntity, requestData addons.RequestData) (addons.CompEntity, error) {

	id := requestData.GetValue("id")

	err := entity.Restore(id)
	if err != nil {
		return nil, err
	}

	log.Debug("done with get entity\n")

	return entity, nil
}

func getAll(entity addons.CompEntity, requestData addons.RequestData) (addons.CompEntity, error) {

	err := entity.Restore("")
	if err != nil {
		return nil, err
	}

	log.Debug("done with get entity list\n")

	return entity, nil
}

// ===== end default functions ===== 

type composedEntity struct {
	isList bool
	l      []*Entity
	s      *Entity
}

func newCompEntity(entity interface{}) addons.CompEntity {
	entityType := reflect.TypeOf(entity).String()

	compE := new(composedEntity)

	switch entityType {
	case "*api.Entity":
		compE.isList = false
		e, _ := entity.(*Entity)
		compE.s = e
		compE.l = nil
		log.Debug("is a single entity")
		return compE
	case "[]*api.Entity":
		compE.isList = true
		el, _ := entity.([]*Entity)
		compE.l = el
		compE.s = nil
		return compE
	}

	return nil
}

func (compE *composedEntity) IsList() bool {
	return compE.isList
}

func (compE *composedEntity) SetAttribute(name, value string) {
	if compE.IsList() {
		return
	}

	compE.s.SetAttribute(name, value)
}

func (compE *composedEntity) GetAttribute(name string) string {
	if compE.IsList() {
		panic("cant get attribute from a list of elements")
	}

	return compE.s.GetAttribute(name)
}

func (compE *composedEntity) Store() error {
	if compE.IsList() {
		// store a list to database
		return errors.New("can't store a list of entities")
	}

	// store a single element to database
	return compE.s.Store()
}

func (compE *composedEntity) Restore(id string) error {

	if len(id) > 0 {
		err := compE.s.Restore(id)
		compE.isList = false
		return err
	}

	err := compE.restoreList()

	return err
}

func (compE *composedEntity) restoreList() error {

	collection := "projectentities"
	ml := make([]map[string]interface{}, 0)
	filter := bson.M{"projectId": compE.s.projectId, "entityName": compE.s.name}
	err := db.RestoreList(&ml, filter, collection)
	if err != nil {
		return err
	}
	log.Debug("restored list of entities from database %v\n", ml)

	for _, entry := range ml {
		delete(entry, "projectId")
		ent := new(Entity)
		ent.attributes = make(map[string]attribute, 0)
		for ak, _ := range compE.s.attributes {
			a := new(attribute)
			ent.attributes[ak] = *a
		}
		ent.name = compE.s.name

		for k, v := range entry {
			if k == "_id" {
				id := v.(bson.ObjectId)
				ent.SetAttribute("id", id.Hex())
				continue
			}
			ent.SetAttribute(k, v.(string))
		}
		compE.l = append(compE.l, ent)
	}

	compE.isList = true

	log.Debug("entity list = %v\n", compE.l)
	return nil
}

func (compE *composedEntity) ToMap() interface{} {
	if compE.IsList() {
		// restore a list to database
		log.Debug("need to process a list o entities")

		if len(compE.l) < 1 {
			return nil
		}

		list := make([]map[string]interface{}, 0)
		for _, e := range compE.l {
			m := make(map[string]interface{}, 0)
			for name, attr := range e.attributes {
				m[name] = attr.value
			}
			list = append(list, m)
		}

		log.Debug("entity list as map %v\n", list)
		return list

	} else {
		// restore a single element to database
		return compE.s.ToMap()
	}

	return nil
}

func (compE *composedEntity) List() []addons.CompEntity {
	if compE.IsList() {
		l := make([]addons.CompEntity, 0)
		for _, E := range compE.l {
			l = append(l, E)
		}
		return l
	}

	return []addons.CompEntity{compE.s}
}

type Entity struct {
	name      string
	projectId string
	//addonsId []string
	Function   map[string]addons.Method
	attributes map[string]attribute
}

type attribute struct {
	value string
	t     int
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

func (e *Entity) GetAttribute(name string) string {
	attr, ok := e.attributes[name]
	if ok {
		return attr.value
	}

	return ""
}

func (e *Entity) List() []addons.CompEntity {
	return nil
}

func (e *Entity) AddMethod(method, path string, f addons.Method) {
	if e.Function == nil {
		e.Function = make(map[string]addons.Method)
	}
	pattern := createPattern(method, path)
	e.Function[pattern] = f
}

func (e *Entity) ToMap() map[string]interface{} {
	m := make(map[string]interface{}, 0)
	for name, attr := range e.attributes {
		m[name] = attr.value
	}
	log.Debug("entity as map %v\n", m)

	return m
}

func (e *Entity) Restore(id string) error {
	//collection := "foraddon_"+e.projectId+"_"+e.name
	collection := "projectentities"
	m := make(map[string]interface{})
	m["_id"] = ""
	delete(m, "id")

	filter := bson.M{"_id": bson.ObjectIdHex(id), "projectId": e.projectId, "entityName": e.name}
	err := db.RestoreOne(m, filter, collection)
	if err != nil {
		return err
	}

	m["id"] = ""

	for k, v := range m {
		if k == "_id" || k == "id" {
			e.SetAttribute("id", id)
			continue
		}
		e.SetAttribute(k, v.(string))
	}
	delete(m, "_id")
	return nil
}

func (e *Entity) Store() error {
	//collection := "foraddon_"+pid+"_"+e.name
	collection := "projectentities"
	m := e.ToMap()
	id := bson.NewObjectId()
	m["_id"] = id
	m["projectId"] = e.projectId
	m["entityName"] = e.name
	delete(m, "id")
	err := db.Store(m, collection)
	delete(m, "_id")

	if err == nil {
		e.SetAttribute("id", id.Hex())
		return nil
	}
	return err
}

type EntityContainer map[string]*Entity

func NewEntitiesList() *EntityContainer {
	el := make(EntityContainer, 0)

	return &el
}

func (es EntityContainer) NewEntity(name string) addons.Entity {
	e := new(Entity)
	e.name = name
	e.Function = make(map[string]addons.Method)
	e.attributes = make(map[string]attribute)
	e.AddAttribute("id", addons.String)
	e.AddAttribute("entityName", addons.String)
	e.AddMethod("GET", "/{id}", getOne)
	e.AddMethod("GET", "/", getAll)
	e.AddMethod("POST", "/", newEntity)
	es[name] = e
	return e
}

func (es *EntityContainer) GetEntity(name string) addons.Entity {
	e := (*es)[name]
	return e
}

func (es EntityContainer) Run(entityName, method, path string, data Data) (addons.CompEntity, error) {
	entity, ok := es[entityName]
	if !ok {
		return nil, errors.New("cant't find entity")
	}

	// find method to be run
	// run method
	var f addons.Method
	for pattern, v := range entity.Function {
		log.Debug("pattern is = %v", pattern)
		matching, _ := regexp.MatchString(pattern, method+":"+path)
		if matching {
			log.Debug("matching")
			f = v
			break
		}
		log.Debug("not matching")

	}

	if f == nil {
		return nil, errors.New("entity don't handle this action")
	}

	compE := newCompEntity(entity)

	resultEntity, err := f(compE, data)
	if err == nil {
		return resultEntity, nil
	}

	return nil, err
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
		for _, v := range pathSlice {
			if v[0] == '{' {
				pattern = pattern + "[0-9a-z_\\ \\.\\+\\-]{24,}/"
			} else {
				pattern = pattern + v + "/"
			}
		}
	}
	pattern = pattern + "$"
	return pattern
}
