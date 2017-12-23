package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/urfave/cli"
)

const (
	ghGraphQLAPIRoot string = "https://api.github.com/graphql"
	ghTokenLocation  string = ".ghportfolio/token"
)

type requestHeader struct {
	key   string
	value string
}

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func getToken() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fullTokenPath := path.Join(usr.HomeDir, ghTokenLocation)
	dat, err := ioutil.ReadFile(fullTokenPath)
	if err != nil {
		fmt.Printf("Error: failed to read Github Access Token at: %s\n", fullTokenPath)
		return "", err
	}

	if string(dat) == "" {
		fmt.Println("Github Access Token is empty. Please run `ghportfolio setup`")
	}

	return strings.TrimSpace(string(dat)), nil
}

func toolSetup() error {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	_ = usr

	return nil
}

func doRequest(method string, url string, body string, headers ...*requestHeader) (*http.Response, error) {
	reader := strings.NewReader(body)

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	ghToken, err := getToken()
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		req.Header.Set(header.key, header.value)
	}
	req.Header.Set("Authorization", "bearer "+ghToken) // always need to authenticate for the GH GraphQL API

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getProjects() (string, error) {
	query := `
    query {
      user(login: "thundergolfer") {
        repositories(first: 50) {
          edges {
            node {
             	name
              url
            }
          }
        }
      }
    }
  `
	payload := graphQLRequest{
		Query:     query,
		Variables: map[string]interface{}{},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := doRequest("POST", ghGraphQLAPIRoot, string(body))
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func main() {
	app := cli.NewApp()
	app.Name = "ghportfolio"

	app.Usage = "for catching up on the activity and health of your public Github projects"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Action = func(c *cli.Context) error {
		name := "person"
		if c.NArg() > 0 {
			name = c.Args().Get(0)
		}
		if c.String("lang") == "spanish" {
			fmt.Println("Hola", name)
		} else {
			fmt.Println("Hello", name)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "display all public repos under your profile, and others you've assigned to yourself",
			Action: func(c *cli.Context) error {
				fmt.Println("Repos: ", c.Args().First())
				projectsDetails, err := getProjects()
				if err != nil {
					fmt.Println("Failed to list projects")
					fmt.Println(err.Error())
				} else {
					fmt.Println(projectsDetails)
				}

				return nil
			},
		},
		{
			Name:    "interest",
			Aliases: []string{"i"},
			Usage:   "display historical stats on stars, forks, and clones of a project (90 days max)",
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "setup",
			Aliases: []string{"i"},
			Usage:   "setup configuration and local data files for this CLI tool",
			Action: func(c *cli.Context) error {
				fmt.Println("Setting up the ghportfolio tool...")
				err := toolSetup()
				if err != nil {
					fmt.Printf("Error: failed to setup the CLI tool, sorry. %v", err.Error())
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}
