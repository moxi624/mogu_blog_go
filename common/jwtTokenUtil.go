package common

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/15 9:10 上午
 * @version 1.0
 */
/*type jwtTokenUtil struct {
	AdminUid   string
	Role       string
	CreateTime string
	jwt.StandardClaims
}

func (jwtTokenUtil) CreateJWT(username string, adminUid string, roleName string, audience string, issuer string, TTLMillis int64, base64Security string) string {
	jwtKey, _ := base64.StdEncoding.DecodeString(base64Security)
	expireTime := time.Now().Add(time.Duration(TTLMillis) * time.Millisecond).Unix()
	claims := &jwtTokenUtil{
		adminUid,
		roleName,
		time.Now().String(),
		jwt.StandardClaims{
			Audience:  audience,
			ExpiresAt: expireTime,
			Issuer:    issuer,
			NotBefore: time.Now().Unix(),
			Subject:   username,
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := tokenObj.SignedString(jwtKey)
	return tokenStr
}
func (jwtTokenUtil) ParseToken(tokenString string, base64Security string) (*jwt.Token, *jwtTokenUtil, error) {
	s := &jwtTokenUtil{}
	token, err := jwt.ParseWithClaims(tokenString, s, func(token *jwt.Token) (i interface{}, err error) {
		return base64Security, nil
	})

	return token, s, err
}

var JwtTokenUtil = &jwtTokenUtil{}*/
