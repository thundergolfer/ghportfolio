package main

import (
	"bufio"
	"bytes"
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

	"github.com/buger/goterm"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

const (
	ghGraphQLAPIRoot      string = "https://api.github.com/graphql"
	ghRestAPIRoot         string = "https://api.github.com/"
	appDataFolderLocation string = ".ghportfolio"
	ghTokenLocation       string = ".ghportfolio/token"
	ghUsernameLocation    string = ".ghportfolio/username"
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
	GhPageSize int
}

func getToken() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fullTokenPath := path.Join(usr.HomeDir, ghTokenLocation)
	dat, err := ioutil.ReadFile(fullTokenPath)
	if err != nil {
		return "", err
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
		return "", err
	}

	return strings.TrimSpace(string(dat)), nil
}

func fileExistsAndAccesible(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		// other error
		return false
	}
	return true
}

func validateSetup() (bool, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	if !fileExistsAndAccesible(path.Join(usr.HomeDir, ghTokenLocation)) {
		return false, err
	}

	username, err := getUsername()
	if err != nil {
		return false, err
	}
	if username == "" {
		fmt.Println("Github Username is empty. Please run `ghportfolio setup`")
		return false, nil
	}

	if !fileExistsAndAccesible(path.Join(usr.HomeDir, ghUsernameLocation)) {
		return false, err
	}
	token, err := getToken()
	if err != nil {
		return false, err
	}
	if token == "" {
		fmt.Println("Github Token is empty. Please run `ghportfolio setup`")
		return false, nil
	}

	return true, nil
}

func appSetup() error {
	fmt.Println("Running Setup")
	fmt.Println("--------------------------------------")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return err
	}

	pth := path.Join(usr.HomeDir, appDataFolderLocation)

	if _, err = os.Stat(pth); os.IsNotExist(err) {
		err = os.Mkdir(pth, 0777)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Github username: ")
	text, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(usr.HomeDir, ghUsernameLocation), []byte(text), 0777)
	if err != nil {
		return err
	}

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter a Github Access Token (goto 'Settings/Developer settings/Personal Access Tokens') with 'user', 'notifications', and 'push' permissions: ")
	text, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(usr.HomeDir, ghTokenLocation), []byte(text), 0777)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------------------")
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

func (app *App) getInterestStats(project string) (map[string]int, map[string]int, map[string]int, error) {
	query := fmt.Sprintf("repos/%s/%s/events", app.GhUsername, project)

	header := &requestHeader{key: "Authorization", value: "bearer " + app.GhToken}
	resp, err := doRequest("GET", ghRestAPIRoot+query, "", header)
	if err != nil {
		return nil, nil, nil, err
	}

	respBody, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return nil, nil, nil, err

	}

	res := []map[string]interface{}{}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, nil, nil, err

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
		return nil, nil, nil, err

	}

	respBody, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return nil, nil, nil, err

	}

	res2 := map[string]interface{}{}
	if err := json.Unmarshal(respBody, &res2); err != nil {
		return nil, nil, nil, err
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

	return stars, forks, clones, nil
}

func (app *App) getInterest(project string, sparklineDisplay bool) (string, error) {
	stars, forks, clones, err := app.getInterestStats(project)
	if err != nil {
		return "", err
	}

	interestStr := ""
	if sparklineDisplay {
		interestStr += fmt.Sprintf("       %s|\n", timelineHeader()) // 1 less space as dodgy alingment hack
		interestStr += fmt.Sprintf("Stars:  %s\n", timelineCountSparkline(stars, "stars"))
		interestStr += fmt.Sprintf("Forks:  %s\n", timelineCountSparkline(forks, "forks"))
		interestStr += fmt.Sprintf("Clones: %s\n", timelineCountSparkline(clones, "clones"))
	} else {
		interestStr += fmt.Sprintf("        %s|\n", timelineHeader())
		interestStr += fmt.Sprintf("Stars:  %s|\n", TimelineCount(stars))
		interestStr += fmt.Sprintf("Forks:  %s|\n", TimelineCount(forks))
		interestStr += fmt.Sprintf("Clones: %s|\n", TimelineCount(clones))
	}

	return interestStr, nil
}

func (app *App) getInterestTotals(project string, sparklineDisplay bool) (string, error) {
	repo := Repository{
		FullName: app.GhUsername + "/" + project,
	}
	stargazers, err := app.Stargazers(repo)
	if err != nil {
		return "", err
	}

	interest := ""

	min := time.Now()
	max := time.Now()
	// find range of chart
	for _, stargazer := range stargazers {
		var curr time.Time = stargazer.StarredAt
		if curr.Before(min) {
			min = stargazer.StarredAt
		}
	}

	stars := map[string]int{}

	// map all stargazers into day buckets
	for _, stargazer := range stargazers {
		dateString := timeToDateStr(stargazer.StarredAt)
		stars[dateString] += 1
	}

	dayTotals := []float64{}
	curr := max
	// sum
	for curr.After(min) {
		dateString := timeToDateStr(curr)
		dayTotals = append(dayTotals, float64(stars[dateString]))

		curr = curr.AddDate(0, 0, -1)
	}

	// reverse the slice day the most recent date is last
	for i, j := 0, len(dayTotals)-1; i < j; i, j = i+1, j-1 {
		dayTotals[i], dayTotals[j] = dayTotals[j], dayTotals[i]
	}

	// accumulate
	for i, _ := range dayTotals {
		if i == 0 {
			continue
		}
		dayTotals[i] += dayTotals[i-1]
	}

	cumulativeTotals := dayTotals
	overallTotal := cumulativeTotals[len(cumulativeTotals)-1]

	chart := goterm.NewLineChart(100, 20)
	data := new(goterm.DataTable)
	data.AddColumn("Time")
	data.AddColumn("Stars")

	for i := 0; i < len(cumulativeTotals); i++ {
		data.AddRow(float64(i), cumulativeTotals[i])
	}

	goterm.Println()

	interest += chart.Draw(data) + fmt.Sprintf(" current: %d", int(overallTotal))

	return interest, nil

}

func timeToDateStr(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d%02d", int(year), int(month), int(day))
}

type jsonResp struct {
	Data struct {
		User struct {
			Repositories struct {
				Edges []map[string]interface{}
			}
		}
	}
}

func (app *App) getProjectsDataDump() (*jsonResp, error) {
	query := `
    query {
      user(login: "%s") {
        repositories(first: 50 privacy:PUBLIC ) {
          edges {
            node {
             	name
              url
              forkCount
              stargazers {
                totalCount
              }
              issues(states: [OPEN] first: 10) {
                  nodes {
                    number
                  }
              }
              pullRequests(first: 10 states:OPEN) {
                nodes {
                  number
                }
              }
            }
          }
        }
      }
    }
  `
	payload := graphQLRequest{
		Query:     fmt.Sprintf(query, app.GhUsername),
		Variables: map[string]interface{}{},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	header := &requestHeader{key: "Authorization", value: "bearer " + app.GhToken}
	resp, err := doRequest("POST", ghGraphQLAPIRoot, string(body), header)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return nil, err
	}

	respBodyJSON := &jsonResp{}

	err = json.Unmarshal(respBody, respBodyJSON)
	if err != nil {
		return nil, err
	}

	return respBodyJSON, nil
}

func (app *App) getPortfolioStats() (string, error) {
	projectDataJSON, err := app.getProjectsDataDump()
	if err != nil {
		return "", err
	}

	var totalForks, totalStars int
	repos := projectDataJSON.Data.User.Repositories.Edges
	for _, r := range repos {
		node := r["node"].(map[string]interface{})
		totalForks += int(node["forkCount"].(float64))
		stargazers := node["stargazers"].(map[string]interface{})
		totalStars += int(stargazers["totalCount"].(float64))
	}

	return fmt.Sprintf("Portfolio Stars: %d  Portfolio Forks: %d", totalStars, totalForks), nil
}

func (app *App) getProjects(filter bool) (string, error) {
	respBodyJSON, err := app.getProjectsDataDump()
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Name", "Open Issues", "Open PRs"})
	repos := respBodyJSON.Data.User.Repositories.Edges
	for _, r := range repos {
		node := r["node"].(map[string]interface{})

		var openIssues, openPrs string
		if hasOpenIssues(node) {
			openIssues = " ! "
		} else {
			openIssues = ""
		}
		if hasOpenPRs(node) {
			openPrs = " ! "
		} else {
			openPrs = ""
		}

		if filter && (openPrs == "") && (openIssues == "") {
			continue
		}
		row := []string{node["name"].(string), openIssues, openPrs}
		table.Append(row)
	}

	table.Render()
	return buf.String(), nil
}

func hasOpenIssues(node map[string]interface{}) bool {
	issues := node["issues"].(map[string]interface{})
	issuesNodes := issues["nodes"].([]interface{})

	return len(issuesNodes) > 0
}

func hasOpenPRs(node map[string]interface{}) bool {
	issues := node["pullRequests"].(map[string]interface{})
	prNodes := issues["nodes"].([]interface{})

	return len(prNodes) > 0
}

func main() {
	valid, err := validateSetup()
	if err != nil {
		fmt.Println("Please run `ghportfolio setup`")
		fmt.Println(err.Error())

		err = appSetup()
		if err != nil {
			fmt.Printf("Error: failed to setup the CLI tool, sorry. %v", err.Error())
			os.Exit(1)
		}
	}
	if !valid {
		err = appSetup()
		if err != nil {
			fmt.Printf("Error: failed to setup the CLI tool, sorry. %v", err.Error())
			os.Exit(1)
		}
	}

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
		GhPageSize: 100,
	}

	app := cli.NewApp()
	app.Name = "ghportfolio"
	app.Usage = "for catching up on the activity and health of your public Github projects"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:    "totals",
			Aliases: []string{"t"},
			Usage:   "display overall interest in your profile/portfolio",
			Action: func(c *cli.Context) error {
				portfolioStats, err := driver.getPortfolioStats()
				if err != nil {
					fmt.Println("Failed to list projects")
					fmt.Println(err.Error())
				} else {
					fmt.Println(portfolioStats)
				}

				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "display all public repos under your profile",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "filter"},
			},
			Action: func(c *cli.Context) error {
				projectsDetails, err := driver.getProjects(c.Bool("filter"))
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
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "chart"},
				cli.BoolFlag{Name: "totals"},
			},
			Action: func(c *cli.Context) error {
				var interest string
				if c.Bool("totals") {
					interest, err = driver.getInterestTotals(c.Args().First(), c.Bool("chart"))
				} else {
					interest, err = driver.getInterest(c.Args().First(), c.Bool("chart"))
				}
				if err != nil {
					fmt.Println("Failed to get project interest info")
					fmt.Println(err.Error())
				}
				fmt.Println(interest)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
