package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Branch struct {
	Name string
	//Commit map[string]interface{}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//get branches from git
	var branches []Branch

	branches = getBranches(1, branches)

	//git config merge.pin.driver true
	//echo sitemap.xml merge=pin >> .git/info/attributes

	for _, branch := range branches {

		//merge with API

		mergeABranch("Robot merging", os.Getenv("GIT_MAIN_BRANCH"), branch.Name)

		//merge with CMD

		//cd '/Applications/MAMP/htdocs/<repo folder>' &&
		/*cmdString := "git checkout '" + branch.Name + "' && git merge master"

		cmd := exec.Command(cmdString)

		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}*/
	}

	fmt.Printf("%v branches processed", len(branches))
}

/**
 * The Repo Merging API supports merging branches in a repository.
 */
func mergeABranch(commitMessage string, base string, head string) {

	mergeABranchEndpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/merges",
		os.Getenv("GIT_ORG"), os.Getenv("GIT_REPO"))

	mergeABranchParams := map[string]string{
		"commit_message": commitMessage,
		"base":           base,
		"head":           head,
	}

	postBody, _ := json.Marshal(mergeABranchParams)

	requestBody := bytes.NewBuffer(postBody)

	request, err := http.NewRequest(http.MethodPost, mergeABranchEndpoint, requestBody)

	request.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("GIT_TOKEN")))

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}

/**
 * get repo branches
 */
func getBranches(page int, branches []Branch) []Branch {

	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches",
		os.Getenv("GIT_ORG"), os.Getenv("GIT_REPO"))

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)

	request.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("GIT_TOKEN")))

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var response []Branch

	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalln(err)
	}

	for _, branch := range response {
		branches = append(branches, branch)
	}

	linkHeader := resp.Header.Get("link")

	if strings.Contains(linkHeader, "rel=\"next\"") {
		page++

		/*
			//if want to limit merges
			if page == 3 {
				return branches
			}*/

		return getBranches(page, branches)
	}

	return branches
}
