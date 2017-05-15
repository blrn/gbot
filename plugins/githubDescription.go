package plugins

import (
	"context"
	"fmt"
	"github.com/blrn/gbot/bot"
	"github.com/google/go-github/github"
	"github.com/mattermost/platform/model"
	"regexp"
)

const github_url_regex = `(?:(?:https?:\/\/)?github\.com\/)(.+?)\/(.+?)(?:/|\s|$)`

func init() {
	bot.Register("githubDescription", func(gbot *bot.Bot) {
		gbot.MatchRegex(onMessage, github_url_regex)
	})

}

func onMessage(incommingMessage *model.Post, gbot *bot.Bot) {
	_, repositories := IsGithubRepo(incommingMessage.Message)
	for _, repo := range repositories {
		msg := createMessage(*repo)
		resp := new(model.Post)
		resp.ChannelId = incommingMessage.ChannelId
		resp.Message = msg
		resp.RootId = incommingMessage.Id
		gbot.SendMessage(resp)
	}
}

var compiled_regex *regexp.Regexp

type myRepo struct {
	Owner       string
	Name        string
	Description string
	url         string
}

func createMessage(repo github.Repository) string {
	/*
		## [Name](Repo URL)
		**Last commit maybe**
		_Description_
	*/
	nameLine := fmt.Sprintf("## [%s](%s)\n", repo.GetFullName(), repo.GetURL())
	descriptionLine := fmt.Sprintf("_%s_\n\n", repo.GetDescription())
	statLine := fmt.Sprintf("Stars: %d | Forks: %d | Open Issues: %d\n\n", repo.GetStargazersCount(), repo.GetForksCount(), repo.GetOpenIssuesCount())
	return nameLine + descriptionLine + statLine
}

func IsGithubRepo(url string) (bool, []*github.Repository) {
	if compiled_regex == nil {
		compiled_regex = regexp.MustCompile(github_url_regex)
	}
	res := compiled_regex.FindAllStringSubmatch(url, -1)
	if res == nil || len(res) == 0 {
		return false, nil
	} else {
		repos := make([]*github.Repository, 0, len(res))
		for _, repoRes := range res {
			if len(repoRes) == 3 {

				tempRepo := myRepo{Owner: repoRes[1], Name: repoRes[2]}
				ghRepo := tempRepo.fetchRepository()
				if ghRepo != nil {
					repos = append(repos, ghRepo)
				}
			}
		}
		return true, repos
	}

}

func (r *myRepo) fetchRepository() *github.Repository {
	client := github.NewClient(nil)
	repo, resp, err := client.Repositories.Get(context.Background(), r.Owner, r.Name)
	if err != nil {
		fmt.Printf("ERROR: could not get myRepo %s/%s\n", r.Owner, r.Name)
		fmt.Println(err)
		return nil
	} else {
		fmt.Printf("myRepo: %+v\n", repo)
		fmt.Printf("resp: %+v\n", resp)
		return repo
	}
}
