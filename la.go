package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	apiURL     = "https://learn-anything.xyz/api/maps"
	apiQuery   = "q"
	client *http.Client
)

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		},
	}
}

// Repo is a GitHub repo
type Repo struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Owner       *repoOwner `json:"owner"`
	URL         string     `json:"html_url"`
	Stars       int64      `json:"stargazers_count"`
	Topics      []string   `json:"topics"`
	Lang        string     `json:"language"`
}

// FullName returns standard "owner/repo" format
func (r *Repo) FullName() string {
	return fmt.Sprintf("%s/%s", r.Username(), r.Name)
}

// Username is GitHub user login
func (r *Repo) Username() string { return r.Owner.Login }

type repoOwner struct {
	Login string `json:"login"`
}

type apiResponse struct {
	Repos []*Repo `json:"items"`
	Total int     `json:"total_count"`
}

// fetchRepos fetches all repos with topic "alfred-workflow" from GitHub.
func fetchRepos() ([]*Repo, error) {
	repos := []*Repo{}
	var (
		pageCount int
		pageNum   = 1
	)

	for {
		if pageCount != 0 && pageNum > pageCount {
			break
		}
		log.Printf("fetching page %d of %d ...", pageNum, pageCount)

		URL, _ := url.Parse(apiURL)
		q := URL.Query()
		q.Set("page", fmt.Sprintf("%d", pageNum))
		q.Set("q", apiQuery)
		URL.RawQuery = q.Encode()

		log.Printf("fetching %s ...", URL)
		req, err := http.NewRequest("GET", URL.String(), nil)
		if err != nil {
			return nil, err
		}
		for k, v := range apiHeaders {
			req.Header.Add(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		log.Printf("[%d] %s", resp.StatusCode, URL)
		if resp.StatusCode > 299 {
			return nil, errors.New(resp.Status)
		}

		data, err := ioutil.ReadAll(resp.Body)
		r := apiResponse{}
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, err
		}
		repos = append(repos, r.Repos...)

		// Populate pageCount if unset
		if pageCount == 0 {
			pageCount = r.Total / 100
			if math.Mod(float64(r.Total), 100.0) > 0.0 {
				pageCount++
			}
		}
		pageNum++

	}
	return repos, nil
}
