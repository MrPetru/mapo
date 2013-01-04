/*
objectspace contiene la definizione dei oggetti che poi verranno usati al interno
dell'applicazione.
*/
package objectspace

import (
    "fmt"
    "crypto/md5"
)

func Md5sum(value string) string {
    sum := md5.New()
    sum.Write([]byte(value))

    result := fmt.Sprintf("%x", sum.Sum(nil))

    //fmt.Printf("md5sum is %s\n", result)
    return result
}
