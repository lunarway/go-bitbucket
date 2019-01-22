package bitbucket

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
)

type Link struct {
	Href string
}

type PullRequestRepository struct {
	Name     string
	Uuid     string
	FullName string
	Links    map[string]Link
}

type Commit struct {
	Links map[string]Link
	Hash  string
}

type Person struct {
	DisplayName string
	AccountId   string
	Links       map[string]Link
	Nickname    string
	Uuid        string
	Username    string
}

type PullRequest struct {
	Id                int
	TaskCount         int
	Author            Person
	CloseSourceBranch bool
	Source            struct {
		Commit     Commit
		Repository PullRequestRepository
		Branch     struct {
			Name string
		}
	}
	Destination struct {
		Commit     Commit
		Repository PullRequestRepository
		Branch     struct {
			Name string
		}
	}
	CommentCount int
	State        string
	Links        map[string]Link
	Title        string
	CreatedOn    string
	//Summary
	ClosedBy    Person
	Description string
	Reason      string
	MergeCommit Commit
}

type PullRequests struct {
	c *Client
}

func (p *PullRequests) CreateUntyped(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := p.c.requestUrl("/repositories/%s/%s/pullrequests/", po.Owner, po.RepoSlug)
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) Create(po *PullRequestsOptions) (*PullRequest, error) {
	response, err := p.CreateUntyped(po)
	if err != nil {
		return nil, err
	}
	return decodePullRequest(response)
}

func (p *PullRequests) UpdateUntyped(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	return p.c.execute("PUT", urlStr, data)
}

func (p *PullRequests) GetsUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/"
	response, err := p.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}
	return decodePullRequests(response)
}

func (p *PullRequests) GetWithQueryUntyped(po *PullRequestsOptions, query string) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/?q=" + url.QueryEscape(query)
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) GetWithQuery(po *PullRequestsOptions, query string) (*[]PullRequest, error) {
	response, err := p.GetWithQueryUntyped(po, query)
	if err != nil {
		return nil, err
	}
	return decodePullRequests(response)
}

func (p *PullRequests) GetUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) ActivitiesUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/activity"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) ActivityUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/activity"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) CommitsUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/commits"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) PatchUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/patch"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) DiffUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/diff"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) MergeUntyped(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/merge"
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) DeclineUntyped(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/decline"
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) GetCommentsUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) GetCommentUntyped(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/" + po.CommentID
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) buildPullRequestBody(po *PullRequestsOptions) string {

	body := map[string]interface{}{}
	body["source"] = map[string]interface{}{}
	body["destination"] = map[string]interface{}{}
	body["reviewers"] = []map[string]string{}
	body["title"] = ""
	body["description"] = ""
	body["message"] = ""
	body["close_source_branch"] = false

	if n := len(po.Reviewers); n > 0 {
		body["reviewers"] = make([]map[string]string, n)
		for i, user := range po.Reviewers {
			body["reviewers"].([]map[string]string)[i] = map[string]string{"username": user}
		}
	}

	if po.SourceBranch != "" {
		body["source"].(map[string]interface{})["branch"] = map[string]string{"name": po.SourceBranch}
	}

	if po.SourceRepository != "" {
		body["source"].(map[string]interface{})["repository"] = map[string]interface{}{"full_name": po.SourceRepository}
	}

	if po.DestinationBranch != "" {
		body["destination"].(map[string]interface{})["branch"] = map[string]interface{}{"name": po.DestinationBranch}
	}

	if po.DestinationCommit != "" {
		body["destination"].(map[string]interface{})["commit"] = map[string]interface{}{"hash": po.DestinationCommit}
	}

	if po.Title != "" {
		body["title"] = po.Title
	}

	if po.Description != "" {
		body["description"] = po.Description
	}

	if po.Message != "" {
		body["message"] = po.Message
	}

	if po.CloseSourceBranch == true || po.CloseSourceBranch == false {
		body["close_source_branch"] = po.CloseSourceBranch
	}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}

func decodePullRequests(response interface{}) (*[]PullRequest, error) {
	responseMap := response.(map[string]interface{})

	if responseMap["type"] == "error" {
		return nil, DecodeError(responseMap)
	}

	var pullRequests []PullRequest
	values := responseMap["values"].([]interface{})

	for _, value := range values {
		pullRequest := &PullRequest{}
		err := mapstructure.Decode(value, pullRequest)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, *pullRequest)
	}

	return &pullRequests, nil
}

func decodePullRequest(response interface{}) (*PullRequest, error) {
	responseMap := response.(map[string]interface{})

	if responseMap["type"] == "error" {
		return nil, DecodeError(responseMap)
	}


	pullRequest := &PullRequest{}
	err := mapstructure.Decode(response, pullRequest)
	if err != nil {
		return nil, err
	}

	return pullRequest, nil
}

//func decodePullRequest(response interface{}) (*PullRequest, error) {
//	repoMap := response.(map[string]interface{})
//
//	if repoMap["type"] == "error" {
//		return nil, DecodeError(repoMap)
//	}
//
//	pullRequest := PullRequest{}
//	err := mapstructure.Decode(repoMap, repository)
//	if err != nil {
//		return nil, err
//	}
//
//	return &pullRequest, nil
//}
