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
package addons

// GetAll restituisce la lista dei addon disponibili
func GetAll() []string {
    // crea una lista di tutti i addon
    return nil
}

const (
	String = 0
)

type Data interface {
	GetValue(string) string
}

var Addons map[string]*addon = make(map[string]*addon)

type addon struct {
	id string
	Constructors []func(*Entities)
	dependByAddons interface{}
}

func NewAddon(id string) *addon {
	ad := new(addon)
	//ad.entity = entity
	ad.id = id

	if _, ok := Addons[id]; !ok {
		Addons[id] = ad
		return ad
	}

	return nil
}

func (a *addon) SetConstructor(c func(*Entities)) {
	a.Constructors = append(a.Constructors, c)
}
