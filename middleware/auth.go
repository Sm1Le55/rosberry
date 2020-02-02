package middleware

import (
	"rosberry/database"
	"fmt"
	"net/http"
	"strconv"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth middleware:", r.URL.Path)

		accessKey := r.Header.Get("AccessKey")		
		if accessKey == "" {
			fmt.Println("Bad access key:", r.URL.Path)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := strconv.Atoi(r.Header.Get("UserID"))
		if err != nil {
			fmt.Println("Bad userId:", r.Header.Get("UserID"), r.URL.Path)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !authVerification(userId, accessKey) {
			fmt.Printf("Incorrect combination userId/accessKey by userId=%v. Path:%v\n", userId, r.URL.Path)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authVerification(userId int, accessKey string) bool {
	result := database.AuthQuery(userId, accessKey)
	if result == database.AuthResultSuccess {
		return true
	}
	return false
}
