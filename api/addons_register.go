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
	"mapo/addons/repo/go/scene"
	"mapo/addons/repo/go/shot"
	"mapo/addons/repo/go/shotpatch"

	"mapo/addons"
)

func RegisterAddons() {
	// addons
	registers := []func(addons.Addon){
		scene.Register, shot.Register, shotpatch.Register,
		}

	newAddonContainer()
	for _, r := range(registers) {
		add := new(addon)
		add.dependByAddons = make(map[string]*addon)
		r(add)
		err := add.CreateId()
		if err == nil {
			Addons[add.id] = add
		}
	}
	//scene.Register(&Addons)
	//shot.Register(&Addons)
	//shotpatch.Register(&Addons)
}
