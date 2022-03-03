package login

import (
	"blion-auth/internal/env"
	"blion-auth/internal/logger"
	"blion-auth/internal/models"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

var (
	signKey    *rsa.PrivateKey
	privateKey string
)

// init lee los archivos de firma y validación RSA
func init() {
	c := env.NewConfiguration()
	privateKey = c.App.RSAPrivateKey
	signBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		logger.Error.Printf("leyendo el archivo privado de firma: %s", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		logger.Error.Printf("realizando el parse en auth RSA private: %s", err)
	}
}

// GenerateJWT Genera el token
func GenerateJWT(u *models.User) (string, int, error) {
	tk := jwt.New(jwt.SigningMethodRS256)
	claims := tk.Claims.(jwt.MapClaims)
	claims["user"] = u
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token, err := tk.SignedString(signKey)
	if err != nil {
		logger.Error.Printf("firmando el token: %v", err)
		return "", 70, err
	}
	return token, 29, nil
}
