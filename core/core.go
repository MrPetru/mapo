/*

USERS
    post:/admin/user - create a new user
    post:/admin/user/id - update a specific user data
    get:/admin/user/id - get specific user data
    get:/admin/user - get data of all users

STUDIOS
    post:/admin/studio - create a new studio
    post:/admin/studio/id - update a studio data
    get:/admin/studio/id - get a specific studio data
    get:/admin/studio - get data for all studios

PROJECTS

*/
package core

import (
    "mapo/log"
    
    "errors"
    "fmt"
)

// tipo di dati che riceverà dai moduli superiori
// TODO: non sono convinto su questa
type values map[string][]string

// SetError e usate per aggiunger una chiave in più alla mappa di valori in
// entrata. In caso di errore, al cliente vera restituito lo stesso oggetto
// che lui ha inviato, aggiungendo soltanto il valore del errore.
func (v values) SetError(value error) {
    v["error"] = append(v["error"], value.Error())
}

// GetSingleValue restituisce un valore singolo di tipo string.
// Visto la forma dei dati in entrata, questa funzione ci aiuta a riprender
// soltanto un valore dalla lista intera che il cliente potrebbe inviare.
func (v values) GetSingleValue(name string) (string, error) {
    value, ok := v[name]
    if ok {
        if fmt.Sprintf("%T", value) == "[]string" {

            if len(value) == 1 {
                return value[0], nil
            }
            return *(new(string)), errors.New("requested value is not a single value")
        
        } else {
            return value[0], nil
        }
    }
    
    return *(new(string)), errors.New("not found")
}

// TODO: GetListValue che simile al GetSingleValue ci aiuta a prendere una lista
// di valori dal oggetto in entrata.


// Start, è un nome ambiguo, il scopo di questa funzione è di avviare il processo
// della richiesta verso il modulo core.
// ha, principalmente il ruolo di identificare e eseguire la funzione giusta
// per la risorsa richiesta. Questo è un approccio poco espandibile.
func Start(resourcePath []string, requestMethod string, formValues map[string][]string) interface{} {
    
    if len(resourcePath) < 1 {
        return nil
    }
    
    inValues := make(values,0)
    inValues = formValues
    log.Debug("%v", inValues)
    
    switch r := resourcePath[0]; r {
    
        case "user":
            // POST functions
            if requestMethod == "POST" {
                // check for function using GET method
                switch l := len(resourcePath); l {
                    case 1:
                        user := NewUser(inValues)
                        return user
                    case 2:
                        user := UpdateUser(inValues)
                        return user
                }
            }
            
            // GET functions
            if requestMethod == "GET" {
                // check for function using GET method
                switch l := len(resourcePath); l {
                    case 2:
                        inValues["id"] = []string{resourcePath[1]}
                        user := GetUser(inValues)
                        return user
                    case 1:
                        result := GetUserAll()
                        return result
                }
            }
            
        case "studio":
            // run studio functions
            if requestMethod == "POST" {
                switch l := len(resourcePath); l {
                    case 1:
                        // POST:/admin/studio   ->  create new studio
                        studio := NewStudio(inValues)
                        return studio
                    case 2:
                        // POST:/admin/studio/id    ->  update/edit studio
                        studio := UpdateStudio(inValues)
                        return studio
                }
            }
            
            if requestMethod == "GET" {
                switch l := len(resourcePath); l {
                    case 1:
                        // GET:/admin/studio    ->  get all studios
                        studios := GetStudioAll(inValues)
                        return studios
                    case 2:
                        // GET:/admin/studio/id ->  get studio by id
                        studio := GetStudio(inValues)
                        return studio
                }
            }
            
        case "project":
            // run project functions
    }
    
    return nil
}
