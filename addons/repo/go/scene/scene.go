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

package scene

import (
	"mapo/addons"
	"github.com/maponet/utils/log"
)

func Register() {
	addon := addons.NewAddon("sc_base_v01")
	addon.SetConstructor(constructor)
	//addon.AddDependency = ""
}

func constructor(entities *addons.Entities) {
	// creare le entità qui
	scene := entities.New("scene")
	scene.AddAttribute("id", addons.String)
	scene.AddAttribute("name", addons.String)
	scene.AddAttribute("description", addons.String)

	scene.AddFunction("GET", "/scene/{id}", getOne)
	scene.AddFunction("GET", "/scene", getAll)
	scene.AddFunction("POST", "/scene", newScene)
}

func getOne(entities *addons.Entities, requestData addons.Data) interface{} {//*addons.Entity {
	project := requestData.GetValue("pid")
	//user := requestData.GetValue("currentuid")
	id := requestData.GetValue("id")

	scene := entities.GetEntity("scene")
	scene.Restore(project, id)

	log.Debug("done with get scene\n")

	return &scene
}

func getAll(entities *addons.Entities, requestData addons.Data) interface{} {//[]addons.Entity {
	project := requestData.GetValue("pid")
	//user := requestData.GetValue("currentuid")
	//id := requestData.GetValue("id")

	sceneList := entities.GetEntityList("scene")
	sceneList.Restore(project)

	log.Debug("done with get scene list\n")

	return sceneList
}


func newScene(entities *addons.Entities, requestData addons.Data) interface{} {//*addons.Entity {
	project := requestData.GetValue("pid")
	//user := requestData.GetValue("currentuid")

	name := requestData.GetValue("name")
	description := requestData.GetValue("description")

	scene := entities.GetEntity("scene")
	scene.SetAttribute("name", name)
	scene.SetAttribute("description", description)

	id, err := scene.Store(project)
	if err != nil {
		panic(err)
	}
	scene.SetAttribute("id", id)

	return &scene
}
