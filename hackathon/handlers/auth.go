import (
	"net/http"
	"os"
	"time"
)

func isProduction() bool {
	return os.Getenv("ENV") == "production"
}

func SessionLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	idToken := r.FormValue("idToken")
	if idToken == "" {
		http.Error(w, "ID token is required", http.StatusBadRequest)
		return
	}

	// (Firebaseトークン検証は省略)

	cookie := &http.Cookie{
		Name:     "session",
		Value:    idToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   isProduction(), // 本番では true
	}

	if isProduction() {
		cookie.SameSite = http.SameSiteNoneMode // クロスサイトに必要
	} else {
		cookie.SameSite = http.SameSiteLaxMode // 開発中はこれで十分
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func SessionLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   isProduction(),
	}
	if isProduction() {
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
