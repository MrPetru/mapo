package core

func ExtractSingleValue(data map[string][]string, name string) string {
    v, ok := data[name]
    if !ok {
        return ""
    }
    
    if len(v) < 1 {
        return ""
    }
    
    if len(v) > 1 {
        return ""
    }
    
    return v[0]
}
