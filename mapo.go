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
	"github.com/maponet/utils/log"
	"github.com/maponet/utils/conf"

	"flag"
)

func main() {
	var err error

	// set config defaults
	conf.Set("logLevel", "INFO")

	// parse flags
	var flagLogLevel, flagConfPath string
	flag.StringVar(&flagConfPath, "conf", "/etc/mapo/mapo.conf", "set path to configuration file")
	flag.StringVar(&flagLogLevel, "log", "", "set loglevel [ERROR|INFO|DEBUG]")
	flag.Parse()

	// load config from file
	confErr := conf.ParseFile(flagConfPath)

	// override config settings with command line flags
	if flagLogLevel != "" {
		conf.Set("logLevel", flagLogLevel)
	}

	logLevel, _ := conf.GetString("logLevel")
	err = log.SetLevel(logLevel)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Starting application")
	if confErr == nil {
		log.Info("Loaded configuration file: %s", flagConfPath)
	} else {
		log.Fatal("Can't load config file at: %s (%s)", flagConfPath, confErr.Error())
	}
	log.Info("Setting log level to: %s", logLevel)

	// setup application

	// register with supervisor
	log.Info("Joining supervisor")

	// init db
	log.Info("Initializing db")

	// load addons
	log.Info("Loading addons")

	// register handlers
	log.Info("Registering handlers")

	// start server
	log.Info("Accepting requests")

	// inform supervisor that we are up

	// for each request
	// 	check authentication/authorization

	// 	extract request operation

	// 	extract request arguments

	// 	pass operation and arguments to api.router

	// 	find function mapped to operation

	// 	call function with arguments

	// return result to user

	// close on signal
	log.Info("Closing application")
}
