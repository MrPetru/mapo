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

package addons

const (
	String = 0
)

type Addons interface {
	NewAddon(string) Addon
}

type Addon interface {
	SetConstructor(func(EntityContainer))
	AddDependency(string)
	SetName(string)
	SetAuthor(string)
	SetVersion(int)
}

type EntityContainer interface {
	NewEntity(string) Entity
	GetEntity(string) Entity
}

type CompEntity interface {
	SetAttribute(string, string)
	GetAttribute(string) string
	Restore(string) error
	Store() error
	List() []CompEntity
}

type Entity interface {
	AddAttribute(string, int)
	SetAttribute(string, string)
	AddMethod(string, string, Method)
	Restore(string) error
	Store() error
}

type EntityList interface {
	Restore() error
}

type Method func(CompEntity, RequestData) (CompEntity, error)

type RequestData interface {
	GetValue(string) string
}
