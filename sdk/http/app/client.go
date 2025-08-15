package app

import "fmt"

var fqdn string

func Init(_fqdn string) {
	fqdn = fmt.Sprintf("%s/api", _fqdn)
}
