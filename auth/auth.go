package auth

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	sqlite "../db"
	"github.com/dgrijalva/jwt-go"
)

type userModel struct {
	Name     string
	Password string
}

var db *sql.DB
var err error

func InitBaseForUser() {
	db, err = sqlite.InitDB()
	if err != nil {
		log.Println("Error on DB init")
	}
	sqlite.CreateUsersTable(db)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := &userModel{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
	}
	token := "Invalid auth"
	hasherName := md5.New()
	hasherName.Write([]byte(user.Name))

	hasherPassword := md5.New()
	hasherPassword.Write([]byte(user.Password))

	if hex.EncodeToString(hasherName.Sum(nil)) == "72c5c1a762214c72f3a5cfe043c03454" && hex.EncodeToString(hasherPassword.Sum(nil)) == "8017d0408f41b75489701e3fb1c3e773" {
		token = createToken()
		sqlite.InstertJwt(db, token)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	resp := map[string]interface{}{"auth": token}

	json.NewEncoder(w).Encode(resp)
}

func createToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-app",
		"sub": "authorised",
		"aud": "any",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	tokenString, error := token.SignedString([]byte("anycameracontroller"))
	if error != nil {
		fmt.Println(error)
	}
	return tokenString
}

func ValidateUser(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	log.Println(token)
	rows, err := db.Query("SELECT * FROM sessions WHERE sessionKey=\"" + token + "\"")
	defer rows.Close()
	resp := map[string]interface{}{"pass": "OK"}
	if err != nil {
		panic(err)
		resp = map[string]interface{}{"pass": "Error on query"}
	}
	if rows.Next() {
		resp = map[string]interface{}{"pass": "OK"}
	} else {
		resp = map[string]interface{}{"pass": "Error on validate"}
	}

	json.NewEncoder(w).Encode(resp)
}
