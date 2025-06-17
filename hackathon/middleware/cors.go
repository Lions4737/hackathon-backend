package middleware

import "net/http"

func WithCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
        // 許可するオリジンを動的に検証して追加
        allowedOrigins := []string{
            "http://localhost:3000",
            "https://hackathon-frontend-wheat-tau.vercel.app", // ← Vercel の本番ドメイン
        }

        for _, o := range allowedOrigins {
            if o == origin {
                w.Header().Set("Access-Control-Allow-Origin", o)
                break
            }
        }
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
