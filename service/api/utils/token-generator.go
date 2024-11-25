package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var secretKey = []byte("Deadpool-Is-The-Marvel-Jesus")

func TokenGenerator(userId string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

// CheckToken verifica il token JWT e restituisce i dati contenuti o un errore
func GetTokenInfo(tokenString string) (map[string]interface{}, error) {
	// Parsing e validazione del token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Controllo del metodo di firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo di firma non valido: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("errore durante il parsing del token: %w", err)
	}

	// Controlla se il token è valido
	if !token.Valid {
		return nil, fmt.Errorf("token non valido")
	}

	// Recupera i claims (dati) dal token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("impossibile recuperare i claims dal token")
	}

	// Ritorna i claims come una mappa generica
	return claims, nil
}
