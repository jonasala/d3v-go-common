package middleware

import (
	"context"
	"net/http"
)

//CarregarUsuario do header IDUsuario para o contexto. Utilize a chave IDUsuario
func CarregarUsuario(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), IDUsuario, r.Header.Get("IDUsuario"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
