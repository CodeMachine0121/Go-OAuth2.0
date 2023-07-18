package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	fmt.Println("Initializing Google OAuth2")
	fmt.Println("Client ID: " + os.Getenv("GOOGLE_CLIENT_ID"))
	fmt.Println("Client Secret: " + os.Getenv("GOOGLE_CLIENT_SECRET"))

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://6f72-223-136-43-237.ngrok-free.app/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<!DOCTYPE html><body>
		<a href="/login">Google Login</a>
	</body></html>`
	fmt.Fprintf(w, htmlIndex)
}

var (
	oauthStateString = "pseudo-state-string"
)

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {

	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if IsErrorOccur(err) {
		return ReturnError(err)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if IsErrorOccur(err) {
		return ReturnError(err)
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if IsErrorOccur(err) {
		return ReturnError(err)
	}

	return contents, nil

}

func IsErrorOccur(err error) bool {
	return err != nil
}
func ReturnError(err error) ([]byte, error) {
	return nil, fmt.Errorf(err.Error())
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	http.ListenAndServe(":8080", nil)
}
