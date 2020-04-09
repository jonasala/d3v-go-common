package middleware

import (
	"context"
	"net/http"
)

//CarregarEmpresa do header IDEmpresa para o contexto. Utilize a chave IDEmpresa
func CarregarEmpresa(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), IDEmpresa, r.Header.Get("IDEmpresa"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
