package objectspace


type studio struct {
    id string
    name string
    owners []string
    projects []string
}

func NewStudio() studio {
    s := new(studio)
    
    return *s
}
