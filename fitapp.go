package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"

	"io"
	"log"
	"net/http"
	//"time"

	"html/template"
)

var homeTmpl = template.Must(template.New("home").ParseFiles("templates/epoch.html"))
var homeLoggedOutTmpl = template.Must(template.New("loggout").ParseFiles("templates/loggedout.html"))

func fitbitUserServeAuthorize(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" === FITBIT USER SERVE AUTHORIZE ===")
	ctx := context.Background()
	/*
		conf := &oauth2.Config{
			ClientID:     "22D6FQ",
			ClientSecret: "be9c1fb74ca0d6b8c93deb35ba305093",
			Scopes:       []string{"SLEEP"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.fitbit.com/oauth2/authorize",
				TokenURL: "https://api.fitbit.com/oauth2/token",
			},
		}
	*/
	conf := &oauth2.Config{
		ClientID:     "22DD2F",
		ClientSecret: "a62ee79d8e9ab5b3f6e99c6a775a16b5",
		//Scopes:       []string{"SLEEP", "activity"},
		Scopes: []string{"activity"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.fitbit.com/oauth2/authorize",
			TokenURL: "https://api.fitbit.com/oauth2/token",
		},
	}
	/*
				22DD2F
				a62ee79d8e9ab5b3f6e99c6a775a16b5
				http://127.0.0.1:8080
				https://www.fitbit.com/oauth2/authorize
				https://api.fitbit.com/oauth2/token

		https://www.fitbit.com/oauth2/authorize?response_type=token&client_id=22DD2F&redirect_uri=http%3A%2F%2F127.0.0.1%3A8080&scope=activity%20heartrate%20location%20nutrition%20profile%20settings%20sleep%20social%20weight&expires_in=604800

	*/

	xcallback := oauth2.SetAuthURLParam("redirect_uri", "http://127.0.0.1:8080/callback")
	//xcallback := oauth2.SetAuthURLParam("redirect_uri", "https://app-settings.fitbitdevelopercontent.com/simple-redirect.html")
	xtimeout := oauth2.SetAuthURLParam("expires_in", "325800")
	//xresponse := oauth2.SetAuthURLParam("response_type", "token")
	//fmt.Println()
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	//url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	url := conf.AuthCodeURL("state", xcallback, xtimeout)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	fmt.Println()
	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	fmt.Println("STEP 1")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	fmt.Println("STEP 2")
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("STEP 3")
	client := conf.Client(ctx, tok)
	//client.Get("...")
	fmt.Println(client)

	//redir := oauthClient.AuthorizationURL(tempCred, nil)
	fmt.Println("=== TOKEN 1 ===")
	//fmt.Println(redir)
	fmt.Println("=== TOKEN 2 ===")

	//http.Redirect(w, r, redir, 302)

	//fitbit-app-224008
}

// authHandler reads the auth cookie and invokes a handler with the result.
type authHandler struct {
	//handler  func(w http.ResponseWriter, r *http.Request, c *oauth.Credentials)
	handler  func(w http.ResponseWriter, r *http.Request)
	optional bool
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" *****  SERVE HTTP  *****")
	//cred, _ := session.Get(r)[tokenCredKey].(*oauth.Credentials)
	//if cred == nil && !h.optional {
	//	http.Error(w, "Not logged in.", 403)
	//}

	h.handler(w, r)
}

/*
// response responds to a request by executing the html remplate t with data.
func respond(w http.ResponseWriter, t *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println(" *****  RESPOND  *****")
	if err := t.Execute(w, data); err != nil {
		log.Print(err)
	}
}
*/

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" *****  SERVE HOME  *****")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	//fmt.Println(cred)
	fmt.Println("HOME")

	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	if err := homeLoggedOutTmpl.ExecuteTemplate(w, "loggedout.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//if cred == nil {
	//respond(w, homeLoggedOutTmpl, nil)
	//respond(w, homePage, nil)
	// for TEST
	//respond(w, homeTmpl, nil)
	//} else {
	//respond(w, homeTmpl, nil)
	//}
}

func fitbitConfig() {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "22D6FQ",
		ClientSecret: "be9c1fb74ca0d6b8c93deb35ba305093",
		Scopes:       []string{"SCOPE1", "SCOPE2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.fitbit.com/oauth2/authorize",
			TokenURL: "https://api.fitbit.com/oauth2/token",
		},
	}
	fmt.Println(ctx)
	fmt.Println(conf)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!, HOW DO YOU DO ??\n")
}

func main() {
	//fitbitConfig()
	//ctx := context.Background()
	data, err := ioutil.ReadFile("settings/index.json")
	fmt.Println("what ??")
	fmt.Println(data)
	fmt.Println("YES")
	//fmt.Println(ctx)
	if err != nil {
		fmt.Println(err)
	}

	//
	//helloHandler := func(w http.ResponseWriter, req *http.Request) {
	//	io.WriteString(w, "Hello, world! what happen ??\n")
	//}

	http.Handle("/", &authHandler{handler: serveHome, optional: true})

	http.HandleFunc("/authorize", fitbitUserServeAuthorize)
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
	//http.HandleFunc("/authorize", serveAuthorize)
	//log.Fatal(http.ListenAndServe(":8010", nil))
	/*
		s := &http.Server{
			Addr:           ":8080",
			Handler:        helloHandler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Fatal(s.ListenAndServe())
	*/
}
