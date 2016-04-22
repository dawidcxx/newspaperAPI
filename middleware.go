package main

import "github.com/gin-gonic/gin"
import jwt "github.com/dgrijalva/jwt-go"
import "net/http"
import "strings"
import "fmt"

const defaultAuthSchema = "Bearer"
// Auth aborts the requests if the user is not Authenticated / has invalid token
func Auth (c *gin.Context) {
  auth := c.Request.Header.Get("Authorization")

  if strings.HasPrefix(auth, defaultAuthSchema) {
    token, err := jwt.Parse(auth[len(defaultAuthSchema)+1:], func(token *jwt.Token) (interface{}, error) {
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
      }
      return []byte(SECRET), nil
    })
    
    if err == nil && token.Valid {
      
      userID := int(token.Claims["UserID"].(float64))
      
      c.Set("UserID", userID)
      c.Next()
      return
    }
    
  }

  c.Status(http.StatusUnauthorized)
  c.Abort()

}
