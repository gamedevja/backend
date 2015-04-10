package git

import (
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	owner = `gamedevja`
	repo  = `gamedevja.github.io`
	mode  = `100644`
	head  = `heads/master`
)

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func Push(path, content, message string) error {
	var err error

	ts := &tokenSource{
		&oauth2.Token{AccessToken: os.Getenv("GITHUBTOKEN")},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// Github API
	g := github.NewClient(tc)
	ref, _, _ := g.Git.GetRef(owner, repo, head)
	cSHA := *(ref.Object.SHA)

	com, _, _ := g.Git.GetCommit(owner, repo, cSHA)
	tSHA := *(com.Tree.SHA)
	tree, _, err := g.Git.CreateTree(owner, repo, tSHA, []github.TreeEntry{github.TreeEntry{
		Path:    &path,
		Mode:    &mode,
		Content: &content,
	}})
	if err != nil {
		return err
	}

	comn, _, _ := g.Git.CreateCommit(owner, repo, &github.Commit{
		Message: &message,
		Tree:    tree,
		Parents: []github.Commit{*com},
	})
	cSHAn := comn.SHA

	_, _, err = g.Git.UpdateRef(owner, repo, &github.Reference{
		Ref:    &head, // Due to go-github bug
		URL:    ref.URL,
		Object: &github.GitObject{Type: ref.Object.Type, SHA: cSHAn, URL: ref.Object.URL},
	}, false)
	if err != nil {
		return err
	}

	return nil
}
