package common

import (
	"encoding/base64"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"strings"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/30 1:03 下午
 * @version 1.0
 */

type jwx struct{}

func (jwx) CreateJWT(username string, adminUid string, roleName string, audience string, issuer string, TTLMillis int64, base64Security string) string {
	jwtKey, _ := base64.StdEncoding.DecodeString(base64Security)
	expireTime := time.Now().Add(time.Duration(TTLMillis) * time.Millisecond).Unix()
	token := jwt.New()
	token.Set(jwt.SubjectKey, username)
	token.Set(jwt.AudienceKey, audience)
	token.Set(jwt.IssuerKey, issuer)
	token.Set(jwt.ExpirationKey, expireTime)
	token.Set(jwt.NotBeforeKey, time.Now().Unix())
	token.Set("adminUid", adminUid)
	token.Set("role", roleName)
	token.Set("createTime", time.Now().String())
	payload, err := jwt.Sign(token, jwa.HS256, jwtKey)
	if err != nil {
		fmt.Printf("failed to generate signed payload: %s\n", err)
	}
	return string(payload)
}

func (jwx) ParseToken(tokenString string) jwt.Token {
	token, err := jwt.Parse(strings.NewReader(tokenString))
	if err != nil {
		fmt.Printf("failed to parse JWT token: %s\n", err)
	}
	return token
}

var Jwx = &jwx{}
