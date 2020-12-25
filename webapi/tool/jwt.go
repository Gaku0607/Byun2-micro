package tool

import (
	"strconv"
	"time"

	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/dgrijalva/jwt-go"
)

var Key = []byte("Byun2-micro")

const (
	issuer string        = "Byun2-micro@v1BY_Gaku"
	exp    time.Duration = time.Duration(3) * time.Hour
	before time.Duration = time.Duration(10) * time.Second
)

func ResendJwt(member *models.Claims) (token string, err error) {
	now := time.Now()

	jwtId := member.MemberName + strconv.FormatInt(now.Unix(), 10)

	claims := &models.Claims{
		MemeberId:  member.MemeberId,
		MemberName: member.MemberName,
		IsSeller:   member.IsSeller,
		StandardClaims: jwt.StandardClaims{
			Id:        jwtId,                  //jwt唯一標示符
			Issuer:    issuer,                 //簽發人
			IssuedAt:  now.Unix(),             //當前時間
			ExpiresAt: now.Add(exp).Unix(),    //過期時間
			NotBefore: now.Add(before).Unix(), //指定時間後生效
		},
	}

	tokenclaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if token, err = tokenclaims.SignedString(Key); err != nil {
		return "", models.NewErr(models.ERROR_SERVER_FAILD, err)
	}
	return
}
