package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/dgrijalva/jwt-go" 
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET")) 

func JWTAuth(next http.Handler) http.Handler{

return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
tokenHeader := r.Header.Get("Authorization") 

if tokenHeader == ""{
 logrus.Warn("Missing auth token")
 w.WriteHeader(http.StatusUnauthorized)
 return 
}

tokenParts := strings.Split(tokenHeader, " ")
if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
logrus.Warn("Invalid token format")
w.WriteHeader(http.StatusUnauthorized)
return 
}
token, err := jwt.Parse(tokenParts[1],func(token *jwt.Token) (interface{},error){
return JWTSecret, nil 
})
if err != nil || !token.Valid {
logrus.Warn("Invalid token:",err)
w.WriteHeader(http.StatusUnauthorized) 
return 
}
next.ServeHTTP(w,r)
})
}
jw