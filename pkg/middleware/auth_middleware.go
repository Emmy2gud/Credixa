package middleware
import (
	"context"

	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"payme/pkg/utils"
)

// AuthMiddleware validates user authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//strings.Replace(original, old, new, count)
		//It means:
		//“Find old inside original,
		//replace it with new,
		//only count times.”
		tokenString := strings.Replace(auth, "Bearer ", "", 1)
		// Validate token
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//request.
		//context.WithValue → creates a new copy of the context with a key-value pair.
		//"user_id" → the key
		//claims["user_id"] → the value
		claims := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
//Adminmiddleware 
func Adminmiddleware(next http.Handler)http.Handler{
return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(auth, "Bearer ", "", 1)
		// Validate token
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["role"] != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}



func GetUserID(r *http.Request) (uint, bool) {
	idValue := r.Context().Value("user_id")
	
	if idValue == nil {
		return 0, false
	}

	// JWT stores numbers as float64
	idFloat, ok := idValue.(float64)
	if !ok {
		return 0, false
	}
// fmt.Printf("value=%v, type=%T\n", idValue, idValue)
	return uint(idFloat), true
}

