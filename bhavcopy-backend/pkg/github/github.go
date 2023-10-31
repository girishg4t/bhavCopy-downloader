package github

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/girishg4t/bhavcopy-backend/pkg/config"
	"github.com/google/go-github/v27/github"
	"golang.org/x/oauth2"
)

type GitConnection struct {
	gitAccessToken string
	email          string
	user           string
	nseDrive       string
	bseDrive       string
	client         *github.Client
	repo           string
	serverFilepath string
	ctx            context.Context
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func ConnectToGit(obj config.Symboles) *GitConnection {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	filename := obj.Date + ".csv"
	filePath := config.NSEDrive + "/" + filename
	if obj.Fund == "OPTIONS" {
		filename = obj.Date + ".json"
		filePath = config.OPTIONSDRIVE + "/" + strings.ToUpper(config.NSEDrive) + "/" + filename
	}
	if obj.Exchange == "BSE" {
		filePath = config.BSEDrive + "/" + filename
	}
	conn := GitConnection{
		client:         github.NewClient(tc),
		email:          config.Email,
		user:           config.GitUser,
		repo:           config.GitRepo,
		gitAccessToken: config.GitAccessToken,
		nseDrive:       config.NSEDrive,
		bseDrive:       config.BSEDrive,
		serverFilepath: filePath,
		ctx:            ctx,
	}
	return &conn
}

func (conn *GitConnection) UpdateToGithub(obj config.Symboles) {
	content, _ := ioutil.ReadFile(config.LocalFilePath)

	message := "Added data for " + conn.serverFilepath
	sha := config.GetSha()
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: content,
		SHA:     &sha,
		Committer: &github.CommitAuthor{Name: github.String(conn.user),
			Email: github.String(conn.email)},
	}
	resp, _, err := conn.client.Repositories.UpdateFile(conn.ctx, conn.user, conn.repo,
		conn.serverFilepath, repositoryContentsOptions)
	if err != nil {
		exitErrorf("Unable to upload %q to %q, %v", obj.Date, conn.serverFilepath, err.Error())
	}

	fmt.Printf(*resp.Message)
}

func (conn *GitConnection) UpdateToGithubOptions(content []byte, obj config.Symboles) {
	message := "Added data for " + conn.serverFilepath
	sha := config.GetSha()
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: content,
		SHA:     &sha,
		Committer: &github.CommitAuthor{Name: github.String(conn.user),
			Email: github.String(conn.email)},
	}
	resp, _, err := conn.client.Repositories.UpdateFile(conn.ctx, conn.user, conn.repo,
		conn.serverFilepath, repositoryContentsOptions)
	if err != nil {
		exitErrorf("Unable to upload %q to %q, %v", obj.Date, conn.serverFilepath, err.Error())
	}

	fmt.Printf(*resp.Message)
}

func (conn *GitConnection) ReadIfFileExistsFromGit(obj config.Symboles) [][]string {
	repos, _, _, err := conn.client.Repositories.GetContents(conn.ctx, conn.user, conn.repo,
		conn.serverFilepath, &github.RepositoryContentGetOptions{})
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	c, err := base64.StdEncoding.DecodeString(*repos.Content)
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}
	reader := csv.NewReader(strings.NewReader(string(c)))
	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}
	return csvData
}

func (conn *GitConnection) ReadIfFileExistsFromGitOptions(obj config.Symboles) []byte {
	repos, _, _, err := conn.client.Repositories.GetContents(conn.ctx, conn.user, conn.repo,
		conn.serverFilepath, &github.RepositoryContentGetOptions{})
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	c, err := base64.StdEncoding.DecodeString(*repos.Content)
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	return c
}
