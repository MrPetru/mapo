/*
objectspace contiene la definizione dei oggetti come User, Studio, Project
*/
package objectspace

import (
    "fmt"
    "crypto/md5"
)

// crea la soma md5 di una stringa
func Md5sum(value string) string {
    sum := md5.New()
    sum.Write([]byte(value))

    result := fmt.Sprintf("%x", sum.Sum(nil))

    return result
}
