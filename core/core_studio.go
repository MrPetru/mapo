package core

import (
    "mapo/log"
    "mapo/objectspace"
)

func NewStudio(inValues values) interface{} {
    // create new studio
    log.Msg("executing NewStudio function")
    
    // creamo un nuovo contenitore di tipo studio
    studio := objectspace.NewStudio()
    
    return studio
}

func GetStudio(inValues values) interface{} {
    // create new studio
    
    return nil
}

func GetStudioAll(inValues values) interface{} {
    // create new studio
    
    return nil
}

func UpdateStudio(inValues values) interface{} {
    // create new studio
    
    return nil
}
