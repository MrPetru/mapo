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
    "github.com/maponet/utils/log"
    "mapo/objectspace"

    "net/http"
)

// GetUser restituisce un utente che è gia salvato nella database
// func GetUser(inValues values) interface{} {
func GetUser(out http.ResponseWriter, in *http.Request) {

    log.Info("executing GetUser function")

    errors := NewCoreErr()

    // cearmo un nuovo ogetto/contenitore per il utente richiesto
    user := objectspace.NewUser()

    // aggiorniamo il valore del id del utente, che servirà per ricavare l'utente
    // dal database
    id := in.FormValue("uid")
    err := user.SetId(id)
    if err != nil {
        errors.append("id", err)
    }

    // fermiamo l'esecuzione se fino a questo momento abbiamo incontrato qualche errore
    if len(errors) > 0{
        WriteJsonResult(out, errors, "error")
        return
    }

    // ricavare i dati del utente dalla database
    err = user.Restore()
    if err != nil {
        errors.append("on restore", err)
        WriteJsonResult(out, errors, "error")
        return
    }

    log.Debug("%s", user.GetId())

    WriteJsonResult(out, user, "ok")
}
