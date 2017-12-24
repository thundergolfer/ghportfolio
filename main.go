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
	"time"

	"github.com/urfave/cli"
)

const (
	ghGraphQLAPIRoot   string = "https://api.github.com/graphql"
	ghRestAPIRoot      string = "https://api.github.com/"
	ghTokenLocation    string = ".ghportfolio/token"
	ghUsernameLocation string = ".ghportfolio/username"
)

type requestHeader struct {
	key   string
	value string
}

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// App is the main application object
type App struct {
	GhToken    string
	GhUsername string
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

func getUsername() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fullTokenPath := path.Join(usr.HomeDir, ghUsernameLocation)
	dat, err := ioutil.ReadFile(fullTokenPath)
	if err != nil {
		fmt.Printf("Error: failed to read Github Username at: %s\n", fullTokenPath)
		return "", err
	}

	if string(dat) == "" {
		fmt.Println("Github Username is empty. Please run `ghportfolio setup`")
	}

	return strings.TrimSpace(string(dat)), nil
}

func appSetup() error {
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

	for _, header := range headers {
		req.Header.Set(header.key, header.value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (app *App) getInterest(project string) (string, error) {
	interestStr := ""
	query := fmt.Sprintf("repos/%s/%s/events", app.GhUsername, project)

	header := &requestHeader{key: "Authorization", value: "bearer " + app.GhToken}
	resp, err := doRequest("GET", ghRestAPIRoot+query, "", header)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return "", err
	}

	res := []map[string]interface{}{}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return "", err
	}

	stars := map[string]int{}
	forks := map[string]int{}
	clones := map[string]int{}

	for _, event := range res {
		eType := event["type"].(string)
		eTime, err := time.Parse(time.RFC3339, event["created_at"].(string))
		if err != nil {
			panic(err) // Github API should never really return invalid input
		}

		switch eType {
		case "ForkEvent":
			eDateStr := timeToDateStr(eTime)
			if val, ok := forks[eDateStr]; ok {
				forks[eDateStr] = val + 1
			} else {
				forks[eDateStr] = 1
			}
		case "WatchEvent":
			eDateStr := timeToDateStr(eTime)
			if val, ok := stars[eDateStr]; ok {
				stars[eDateStr] = val + 1
			} else {
				stars[eDateStr] = 1
			}
		default:
		}
	}

	// https: //api.github.com/repos/thundergolfer/interview-with-python/traffic/clones
	query = fmt.Sprintf("repos/%s/%s/traffic/clones", app.GhUsername, project)

	header = &requestHeader{key: "Authorization", value: "bearer " + app.GhToken}
	resp, err = doRequest("GET", ghRestAPIRoot+query, "", header)
	if err != nil {
		return "", err
	}

	respBody, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return "", err
	}

	res2 := map[string]interface{}{}
	if err := json.Unmarshal(respBody, &res2); err != nil {
		return "", err
	}

	clonesList := res2["clones"].([]interface{})

	for _, cloneInterface := range clonesList {
		clone := cloneInterface.(map[string]interface{})
		cloneTime, err := time.Parse(time.RFC3339, clone["timestamp"].(string))
		if err != nil {
			panic(err)
		}
		eDateStr := timeToDateStr(cloneTime)
		if val, ok := clones[eDateStr]; ok {
			clones[eDateStr] = val + int(clone["uniques"].(float64))
		} else {
			clones[eDateStr] = int(clone["uniques"].(float64))
		}
	}

	interestStr += fmt.Sprintf("        %s|\n", timelineHeader())
	interestStr += fmt.Sprintf("Stars:  %s|\n", TimelineCount(stars))
	interestStr += fmt.Sprintf("Forks:  %s|\n", TimelineCount(forks))
	interestStr += fmt.Sprintf("Clones: %s|\n", TimelineCount(clones))

	return interestStr, nil
}

func timeToDateStr(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d%02d", int(year), int(month), int(day))
}

func (app *App) getProjects() (string, error) {
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

	header := &requestHeader{key: "Authorization", value: "bearer " + app.GhToken}
	resp, err := doRequest("POST", ghGraphQLAPIRoot, string(body), header)
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
	token, err := getToken()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	username, err := getUsername()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	driver := App{
		GhToken:    token,
		GhUsername: username,
	}

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
			Usage:   "display all public repos under your profile",
			Action: func(c *cli.Context) error {
				fmt.Println("Repos: ", c.Args().First())
				projectsDetails, err := driver.getProjects()
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
				interest, err := driver.getInterest(c.Args().First())
				if err != nil {
					fmt.Println("Failed to get project interest info")
					fmt.Println(err.Error())
				}
				fmt.Println(interest)
				return nil
			},
		},
		{
			Name:    "setup",
			Aliases: []string{"i"},
			Usage:   "setup configuration and local data files for this CLI tool",
			Action: func(c *cli.Context) error {
				fmt.Println("Setting up the ghportfolio tool...")
				err := appSetup()
				if err != nil {
					fmt.Printf("Error: failed to setup the CLI tool, sorry. %v", err.Error())
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}
