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

package main

import (
    "mapo/database"
    "mapo/managers/addon"
    "mapo/log"
    
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "sync"
    "time"
)

func main() {

    log.SetLevel("INFO")
    
    // create a new connection to database
    database.NewConnection()
    log.Msg("created a new database connection")
    
    // register for all available addons
    addons := addon.GetAll()
    addons = addons
    log.Msg("load addons and generate a list")
    
    // create and start a server http
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT)//, syscall.SIGTERM)
    
    s := new(server)
    // TODO: register this node to load balancing service
    go s.signals(c)
    
    log.Info("start listening for requests")
    log.Msg("close server with message: %v", http.ListenAndServe(":8081", s))
}

type server struct {
    current_connections int
    lock sync.Mutex
    closing bool
}

func (s *server) RequestHandler(out http.ResponseWriter, in *http.Request) {

    log.Msg("executing RequestHandler function")
    // collect request data
    
    // authenticate
    
    // run router
    
    // send response to client
}

func (s *server) ServeHTTP(out http.ResponseWriter, in *http.Request) {
    
    if !s.closing {
        s.lock.Lock()
        s.current_connections++
        s.lock.Unlock()
        
        defer func() {
            s.lock.Lock()
            s.current_connections--
            s.lock.Unlock()
        }()
        
        s.RequestHandler(out, in)
    }
}

func (s *server) signals(c chan os.Signal) {

    _ = <-c
    log.Info("closing ...")
    s.closing = true
    
    // TODO: send notification to load balancing that this node is unavailable
    
    for {
        if s.current_connections == 0 {
            log.Info("bye ... :)")
            os.Exit(1)
        } else {
            log.Info("waiting for %d opened connections", s.current_connections)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
