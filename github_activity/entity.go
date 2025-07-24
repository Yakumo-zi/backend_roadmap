package main

type Repo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Actor struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarId   string `json:"gravatar_id"`
	Url          string `json:"url"`
	AvatarUrl    string `json:"avatar_url"`
}
type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Commit struct {
	Sha      string `json:"sha"`
	Author   Author `json:"author"`
	Message  string `json:"message"`
	Distinct bool   `json:"distinct"`
	Url      string `json:"url"`
}

type Payload struct {
	Action       string      `json:"action,omitempty"`
	RepositoryId int         `json:"repository_id,omitempty"`
	PushId       int64       `json:"push_id,omitempty"`
	Size         int         `json:"size,omitempty"`
	DistinctSize int         `json:"distinct_size,omitempty"`
	Ref          *string     `json:"ref,omitempty"`
	Head         string      `json:"head,omitempty"`
	Before       string      `json:"before,omitempty"`
	Commits      []Commit    `json:"commits,omitempty"`
	RefType      string      `json:"ref_type,omitempty"`
	MasterBranch string      `json:"master_branch,omitempty"`
	Description  interface{} `json:"description"`
	PusherType   string      `json:"pusher_type,omitempty"`
}

type Org struct {
	Id         int    `json:"id"`
	Login      string `json:"login"`
	GravatarId string `json:"gravatar_id"`
	Url        string `json:"url"`
	AvatarUrl  string `json:"avatar_url"`
}
