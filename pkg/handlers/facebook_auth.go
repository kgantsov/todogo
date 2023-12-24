package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"gopkg.in/gin-gonic/gin.v1"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/login_callback/",
		ClientID:     os.Getenv("FACEBOOK_APP_ID"),
		ClientSecret: os.Getenv("FACEBOOK_APP_SECRET"),
		Scopes: []string{
			"public_profile",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "random"
)

func FacebookLogin(c *gin.Context) {
	fmt.Println("!!!!!", c.Request.URL.Query())
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func FacebookLoginCallback(c *gin.Context) {
	fmt.Println("!!!!!", c.Request.URL.Query())

	state := c.Request.URL.Query().Get("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	code := c.Request.URL.Query().Get("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	fmt.Println(
		"???????",
		"https://graph.facebook.com/me?fields=id,name,email,gender&access_token="+token.AccessToken,
	)
	response, err := http.Get(
		"https://graph.facebook.com/me?fields=id,name,email,gender&access_token=" + token.AccessToken,
	)

	defer response.Body.Close()
	fmt.Println(response.Body)
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Println(">>>>>________", contents)

	var res interface{}

	err = json.Unmarshal(contents, &res)

	if err != nil {
		fmt.Println(err.Error(), response.Status)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	fmt.Println("######", res)
	c.JSON(200, gin.H{"result": res})
}
