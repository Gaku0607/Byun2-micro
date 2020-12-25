package models

import "github.com/dgrijalva/jwt-go"

//Jwt參數
type Claims struct {
	MemeberId  uint
	MemberName string
	IsSeller   bool
	jwt.StandardClaims
}
