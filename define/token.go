package define

import (
	"ddaom/domain"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//Creating Access Token
func CreateToken(userToken *domain.UserToken, secretKey string) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = userToken.Authorized
	atClaims["seq_member"] = userToken.SeqMember
	atClaims["email"] = userToken.Email
	atClaims["user_level"] = userToken.UserLevel
	atClaims["exp"] = time.Now().Add(time.Minute * 60 * 12).Unix() // 개발 편의상 12시간 설정
	atClaims["allocated"] = userToken.Allocated

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secretKey)) // 1 Hour
	if err != nil {
		return "", err
	}
	return token, nil
}

// Token validation check
func VerifyToken(tokenString string, secretKey string) (string, error) {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			return EXPIRED_TOKEN, err
		} else {
			return INVALID_TOKEN, err
		}
	}
	return SUCCESS, nil
}

func ExtractTokenMetadata(tokenString string, secretKey string) (*domain.UserToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		authorized, _ := claims["authorized"].(bool)
		seqMember, _ := claims["seq_member"].(float64)
		email, _ := claims["email"].(string)
		userLevel, _ := claims["user_level"].(float64)
		allocated, _ := claims["allocated"].(float64)

		return &domain.UserToken{
			Authorized: authorized,
			SeqMember:  int64(seqMember),
			Email:      email,
			UserLevel:  int(userLevel),
			Allocated:  int8(allocated),
		}, nil
	}
	return nil, err
}
