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

package admin

import (
	"net/http"
)

type mux interface {
	HandleFunc(string, string, func(http.ResponseWriter, *http.Request))
}

func Activate(muxer mux) {

	muxer.HandleFunc("GET", "/admin/user/{uid}", Authenticate(httpGetUser))

	muxer.HandleFunc("POST", "/admin/studio", Authenticate(httpNewStudio))
	muxer.HandleFunc("GET", "/admin/studio", Authenticate(httpGetStudioAll))
	muxer.HandleFunc("GET", "/admin/studio/{sid}", Authenticate(httpGetStudio))
	muxer.HandleFunc("GET", "/admin/studio/{sid}/update", Authenticate(httpUpdateStudio))

	muxer.HandleFunc("POST", "/admin/project", Authenticate(httpNewProject))
	muxer.HandleFunc("GET", "/admin/project", Authenticate(httpGetProjectAll))
	muxer.HandleFunc("GET", "/admin/project/{pid}", Authenticate(httpGetProject))
	muxer.HandleFunc("GET", "/admin/project/{pid}/appendaddon", Authenticate(httpAppendAddon))

	muxer.HandleFunc("GET", "/login/{oauthprovider}", Login)
	//muxer.HandleFunc("GET", "/logout", admin.Logout)

	// OAuth
	// su questo url viene reinderizato il cliente dopo che la procedura di authenticazione
	// sul server del servizio aviene con successo o meno.
	muxer.HandleFunc("GET", "/oauth2callback", OAuthCallBack)

}
