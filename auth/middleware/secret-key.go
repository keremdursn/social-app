package middleware

import "os"

var mySigningKey = []byte(os.Getenv("JWT_KEY"))
