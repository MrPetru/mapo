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

/*
Package addons contains a framework to create Mapo addons and the addons
shipped with the official distribution.
*/
package api

import (
	"mapo/addons"
	"github.com/maponet/utils/log"
	"fmt"
	"errors"
)

const (
	String = 0
)

/*
dati inviati al addon sono caratterizzati di questa interfaccia.
*/
type Data interface {
	GetValue(string) string
}

type addonContainer map[string]*addon

/*
lista globale dei addon disponibili. Identificati da un id unico.
*/
var Addons addonContainer

func newAddonContainer() {
	Addons = make(addonContainer)
}

///*
//Usato al avvio quando i addon vengono registrati nella lista globale dei
//Addons. Pero, viene chiamato dal addon stesso.
//*/
//Func (ac *addonContainer) NewAddon(id string) addons.Addon {
//	ad := new(addon)
//	ad.dependByAddons = make(map[string]*addon)
//
//	if _, ok := Addons[id]; !ok {
//		ad.id = id
//		Addons[id] = ad
//		return ad
//	}
//
//	return nil
//}

/*
definizione de un singolo addon
*/
type addon struct {
	id string
	name string
	author string
	version int

	Constructors []func(addons.EntityContainer)
	dependByAddons map[string]*addon
}

/*
ogni addon ha una funzione che costruisce le entità necessari ad un
funzionamento corretto. Al momento della registrazione del addon,
SetConstructor collega il costruttore al  addon.
*/
func (a *addon) SetConstructor(c func(addons.EntityContainer)) {
	a.Constructors = append(a.Constructors, c)
}

func (a *addon) AddDependency(dependencyId string) {
	dep, ok := Addons[dependencyId]
	if ok {
		a.dependByAddons[dependencyId] = dep
		return
	}
	log.Error("cant find addon with ID=%s", dependencyId)
}

func (a *addon) SetAuthor(author string) {
	a.author = author
}

func (a *addon) SetName(name string) {
	a.name = name
}

func (a *addon) SetVersion(v int) {
	a.version = v
}

func (a *addon) CreateId() error {
	if (len(a.author)>0 && len(a.name)>0 && a.version>0) {
		a.id = fmt.Sprintf("%s:%s:%04d", a.name, a.author, a.version)
		return nil
	}
	return errors.New("cant create addon ID")
}

func orderByDependency(addonsId []string, Addons addonContainer) []string{
	var hasDep []string = make([]string, 0)
	var isDep []string = make([]string, 0)

	for _, addId := range(addonsId) {
		a, ok := Addons[addId]
		if ok {
			cycle(a, &hasDep, &isDep)
		} else {
			log.Error("name %s isn't in addon list", addId)
		}
	}

	result := make([]string, 0)
	// 0 -> n
	for _, addId := range(isDep) {
		if !inList(result, addId) {
			result = append(result, addId)
		}
	}

	// n -> 0
	for i:=len(hasDep); i>0 ; i-- {
		if !inList(result, hasDep[i-1]) {
			result = append(result, hasDep[i-1])
		}
	}

	log.Debug("aranged by dependency addons = %v", result)

	return result
}

func cycle(a *addon, has, is *[]string) {
	if len(a.dependByAddons) < 1 {
		*is = append(*is, a.id)
		return
	}

	*has = append(*has, a.id)

	for _, d := range(a.dependByAddons) {
		cycle(d, has, is)
	}
}

func inList(list []string, elem string) bool {

	for _, s := range(list) {
		if s == elem {
			return true
		}
	}
	return false
}
