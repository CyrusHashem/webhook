package main

import (
	"crypto/subtle"
)

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (a *Auth) Check(username, password string) bool {
	u := subtle.ConstantTimeCompare([]byte(a.Username), []byte(username))
	p := subtle.ConstantTimeCompare([]byte(a.Password), []byte(password))
	return u&p == 1
}
