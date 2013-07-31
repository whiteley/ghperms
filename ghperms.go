package main

import "flag"
import "fmt"
import "log"
import "os"
import "os/user"
import "path"

import "code.google.com/p/goauth2/oauth"
import "github.com/google/go-github/github"

var (
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	code         = flag.String("code", "", "Authorization Code")
)

func auth(config *oauth.Config) *github.Client {
	transport := &oauth.Transport{Config: config}

	token, err := config.TokenCache.Token()
	if err != nil {
		if *clientId == "" || *clientSecret == "" {
			flag.Usage()
			os.Exit(2)
		}
		if *code == "" {
			url := config.AuthCodeURL("")
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			os.Exit(0)
		}
		token, err = transport.Exchange(*code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}

	transport.Token = token

	client := github.NewClient(transport.Client())
	return client
}

func main() {
	flag.Parse()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cachefile := path.Join(usr.HomeDir, ".ghperms.json")
	authURL := "https://github.com/login/oauth/authorize"
	tokenURL := "https://github.com/login/oauth/access_token"

	var config = &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		TokenCache:   oauth.CacheFile(cachefile),
		AuthURL:      authURL,
		TokenURL:     tokenURL,
	}

	client := auth(config)

	repos, err := client.Repositories.List("", nil)
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	} else {
		fmt.Printf("%#v\n\n", repos)
	}
}
