package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	MemeberId  uint
	MemberName string
	IsSeller   bool
	jwt.StandardClaims
}
