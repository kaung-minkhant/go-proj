package main

import (
	"fmt"
	"net/http"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandlerFunc(handler apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{
				Error: err.Error(),
			})
		}
	}
}

func withJWTAuth(handler http.HandlerFunc, store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("With jwt")
		auth := r.Header.Get("authorization")
		token, err := getAndValidateJWTToken(auth)
		if err != nil {
			ErrorLog("JWT validation error %v", err.Error())
			PermissionDenied(w)
			return
		}
		accountId, err := GetAccountIDFromUrl(r)
		if err != nil {
			ErrorLog("Could not get account id %v", err)
			PermissionDenied(w)
			return
		}
		account, err := store.GetAccountById(accountId)
		if err != nil {
			ErrorLog("Could not find account with id %v with error %v", accountId, err)
			PermissionDenied(w)
			return
		}
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			ErrorLog("Claims cannot be found in JWT %v", token)
			PermissionDenied(w)
			return
		}
		claimId := claims.AccountNumber
		if claimId != account.Number.String() {
			ErrorLog("mismatch credentials")
			PermissionDenied(w)
			return
		}
		handler(w, r)
	}
}

// Embedding testings
// type interfaceA interface {
// 	SayHello()
// }
// type structB struct {
// }
// func (str structB) SayHello() {
// 	fmt.Println("hello")
// }
// type strctA struct {
// 	*structB
// }
// func main2() {
// 	var strA strctA
// 	strA.SayHello()
// }
