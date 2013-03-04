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

package shotpatch

import (
	"mapo/addons"
)

func Register(addon addons.Addon) {
	//addon := addonContainer.NewAddon("sh_base_v01_patch")
	addon.SetName("shot_base_structure_patch")
	addon.SetAuthor("maponet")
	addon.SetVersion(1)
	addon.SetConstructor(constructor)
	addon.AddDependency("shot_base_structure:maponet:0001")
}

func constructor(entityContainer addons.EntityContainer) {
	// creare le entità qui
	shot := entityContainer.GetEntity("shot")
	shot.AddAttribute("parentname", addons.String)
}
