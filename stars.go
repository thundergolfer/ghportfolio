// Credit to: https://github.com/caarlos0/starcharts
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var errNoMorePages = errors.New("no more pages to get")

// Stargazer is a star at a given time
type Stargazer struct {
	StarredAt time.Time `json:"starred_at"`
}

type Repository struct {
	FullName        string `json:"full_name"`
	StargazersCount int    `json:"stargazers_count"`
	CreatedAt       string `json:"created_at"`
}

// Stargazers returns all the stargazers of a given repo
func (app *App) Stargazers(repo Repository) (stars []Stargazer, err error) {
	sem := make(chan bool, 10)
	var g errgroup.Group
	var lock sync.Mutex
	for page := 1; page <= app.lastPage(repo); page++ {
		sem <- true
		page := page
		g.Go(func() error {
			defer func() { <-sem }()
			result, err := app.getStargazersPage(repo, page)
			if err == errNoMorePages {
				return nil
			}
			if err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			stars = append(stars, result...)
			return nil
		})
	}
	err = g.Wait()
	sort.Slice(stars, func(i, j int) bool {
		return stars[i].StarredAt.Before(stars[j].StarredAt)
	})
	return
}

func (app *App) getStargazersPage(repo Repository, page int) (stars []Stargazer, err error) {
	// err = app.cache.Get(fmt.Sprintf("%s_%d", repo.FullName, page), &stars)
	// if err == nil {
	// 	ctx.Info("got from cache")
	// 	return
	// }
	var url = fmt.Sprintf(
		"https://api.github.com/repos/%s/stargazers?page=%d&per_page=%d",
		repo.FullName,
		page,
		app.GhPageSize,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return stars, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3.star+json")
	if app.GhToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", app.GhToken))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return stars, err
	}
	defer resp.Body.Close()

	// rate limit
	if resp.StatusCode == http.StatusForbidden {
		time.Sleep(10 * time.Second)
		return app.getStargazersPage(repo, page)
	}
	if resp.StatusCode != http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return stars, err
		}
		return stars, fmt.Errorf("failed to get stargazers from github api: %v", string(bts))
	}
	err = json.NewDecoder(resp.Body).Decode(&stars)
	if len(stars) == 0 {
		return stars, errNoMorePages
	}
	// if err := app.cache.Put(
	// 	fmt.Sprintf("%s_%d", repo.FullName, page),
	// 	stars,
	// 	expire,
	// ); err != nil {
	// 	ctx.WithError(err).Warn("failed to cache")
	// }
	return
}

func (app *App) lastPage(repo Repository) int {
	return (repo.StargazersCount / app.GhPageSize) + 1
}
