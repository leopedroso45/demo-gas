package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	model "demogas.com/m/model"
	"demogas.com/m/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/basic"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"github.com/shaj13/go-guardian/store"
)

var authenticator auth.Authenticator
var cache store.Cache

func main() {
	r := mux.NewRouter()
	//port := os.Getenv("PORT")
	port := "8000"
	setupGoGuardian()
	r.HandleFunc("/v1/auth/token", middleware(http.HandlerFunc(createToken))).Methods("GET")
	r.HandleFunc("/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/editAccount", editAccount).Methods("PUT")
	r.HandleFunc("/removeAccount", removeAccount).Methods("DELETE")
	// r.HandleFunc("/removeAccount", middleware(http.HandlerFunc(removeAccount))).Methods("DELETE")
	http.ListenAndServe("127.0.0.1:"+port, r)
}

func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	if userName == "medium" && password == "medium" {
		return auth.NewDefaultUser("medium", "1", nil, nil), nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

func verifyToken(ctx context.Context, r *http.Request, tokenString string) (auth.Info, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := auth.NewDefaultUser(claims["medium"].(string), "", nil, nil)
		return user, nil
	}
	return nil, fmt.Errorf("invaled token")
}

func setupGoGuardian() {
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*5)
	basicStrategy := basic.New(validateUser, cache)
	tokenStrategy := bearer.New(verifyToken, cache)
	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
	authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}

func middleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Auth Middleware")
		user, err := authenticator.Authenticate(r)
		if err != nil {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("User %s Authenticated\n", user.UserName())
		next.ServeHTTP(w, r)
	})
}

func createToken(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-app",
		"sub": "medium",
		"aud": "any",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	jwtToken, _ := token.SignedString([]byte("secret"))
	w.Write([]byte(jwtToken))
}
func createAccount(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var newUser model.User
	err := decoder.Decode(&newUser)

	if err != nil {
		code := http.StatusBadRequest
		http.Error(w, err.Error(), code)
		return
	}

	newUser.Id = uuid.NewString()
	result, err := mongo.CreateUser(newUser)
	if err != nil {
		code := http.StatusNotAcceptable
		http.Error(w, err.Error(), code)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result.(string)))

}
func editAccount(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newUser model.User
	err := decoder.Decode(&newUser)

	if err != nil {
		code := http.StatusBadRequest
		http.Error(w, err.Error(), code)
		return
	}

	updatedUserBson, err := mongo.EditUser(newUser)
	if err != nil {
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}

	var updatedUser model.User

	bson.Unmarshal(updatedUserBson, &updatedUser)

	updatedUser.ClearUserDetails()

	w.WriteHeader(http.StatusAccepted)
	w.Write(updatedUser.ToJSON())
}

func removeAccount(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newUser model.User
	err := decoder.Decode(&newUser)

	if err != nil {
		code := http.StatusBadRequest
		http.Error(w, err.Error(), code)
		return
	}

	_, err = mongo.DeleteUser(newUser)
	if err != nil {
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}

	newUser.ClearUserDetails()

	w.WriteHeader(http.StatusAccepted)
	w.Write(newUser.ToJSON())
}
