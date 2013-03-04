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
Package api implements the API framework used by addons.
*/
package api

import (
	"mapo/admin"
	"net/http"
)

type mux interface {
	HandleFunc(string, string, func(http.ResponseWriter, *http.Request))
}

func Activate(muxer mux) {

    muxer.HandleFunc("GET", "/api/{pid}/.*", admin.Authenticate(HttpWrapper))
    muxer.HandleFunc("POST", "/api/{pid}/.*", admin.Authenticate(HttpWrapper))
    muxer.HandleFunc("GET", "/api/{pid}", admin.Authenticate(HttpWrapper))

}
