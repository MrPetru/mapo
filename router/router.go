package router

import (
    "mapo/log"
    
    "errors"
    "strings"
)

// New create a routerData container for router.
func New(method, resourcePath string) (r routerData, err error) {
    
    log.Debug(resourcePath)
    if !(method == "GET" || method == "POST") {
        err = errors.New("invalid method, ex: GET, POST, ...")
        return
    }
    
    if len(resourcePath) < 1 {
        err = errors.New("invalid resource path, need at least / as path")
        return
    }
    
    r = *(new(routerData))
    r.method = method
    r.resourcePath = strings.Split(resourcePath[1:], "/")
    log.Debug("%v", r.resourcePath)
    
    return
}

// This data is pushed in from a upper module that will use some interfaces
// defined here.
type routerData struct {

    // user container
    user user
    
    // studio container
    studio studio
    
    // project container
    project project
    
    // user can use a special toke to authenticate
    token string
    
    // path to requested resource
    resourcePath []string
    
    // method of request
    method string
    
    // remaning values from request
    otherValues interface{}
    
    // if request contain files, this will be accessible from here
    files map[string][]byte
    
}

// Run will pass controll to router, real processing of request begin from here
func (rd *routerData) Run() (result interface{}, err error) {
    result = "result from router - request processor"
    return
}

// SetUserLogin store user login in de router data container
func (rd *routerData) SetUserLogin(value string) {
    log.Debug(value)
    rd.user.login = value
}

// SetUserPassword store user password in de router data container
func (rd *routerData) SetUserPassword(value string) {
    log.Debug(value)
    rd.user.password = value
}

// SetUserToken store user token in de router data container
func (rd *routerData) SetUserToken(value string) {
    log.Debug(value)
    rd.token = value
}

// SetOtherValues store remeaning values in de router data container
func (rd *routerData) SetOtherValues(value interface{}) {
    log.Debug("%v", value)
    rd.otherValues = value
}

// Authenticate will try to authenticate the user based on info founded in 
// router data container
func (rd *routerData) Authenticate() (ok bool) {
    log.Msg("executing authentication function")
    ok = true
    
    return
}

type user struct {
    id string
    login string
    name string
    password string
    contacts []interface{}
    description string
    rating float32
    studios map[string]string
}

type studio struct {
    id string
    name string
    owners []string
    administrators []string
    collaborators []string
    projects []string
}

type project struct {
    id string
    name string
    studio string
    administrators []string
    supervisors []string
    artists []string
    guests []string
}
