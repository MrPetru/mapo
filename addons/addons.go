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

/*
dati inviati al addon sono caratterizzati di questa interfaccia.
*/
type Data interface {
	GetValue(string) string
}

/*
lista globale dei addon disponibili. Identificati da un id unico.
*/
var Addons map[string]*addon = make(map[string]*addon)

/*
definizione de un singolo addon
*/
type addon struct {
	id string
	Constructors []func(*EntityContainer)
	dependByAddons interface{}
}

/*
usato al avvio quando i addon vengono registrati nella lista globale dei
addons. Pero, viene chiamato dal addon stesso.
*/
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

/*
ogni addon ha una funzione che costruisce le entità necessari ad un
funzionamento corretto. Al momento della registrazione del addon,
SetConstructor collega il costruttore al  addon.
*/
func (a *addon) SetConstructor(c func(*EntityContainer)) {
	a.Constructors = append(a.Constructors, c)
}
