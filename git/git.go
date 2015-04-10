package git

import (
	"encoding/base64"
	"os"
	"path/filepath"

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

func Push(content *string, path, message string) error {
	var err error

	ts := &tokenSource{
		&oauth2.Token{AccessToken: os.Getenv("GITHUBTOKEN")},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	g := github.NewClient(tc)
	ref, _, _ := g.Git.GetRef(owner, repo, head)
	cSHA := *(ref.Object.SHA)

	com, _, _ := g.Git.GetCommit(owner, repo, cSHA)
	tSHA := *(com.Tree.SHA)

	var enc string
	ext := filepath.Ext(path)
	if ext == ".png" || ext == ".mp3" || ext == ".ogg" {
		enc = "base64"
		b64 := base64.StdEncoding.EncodeToString([]byte(*content))
		content = &b64
	} else {
		enc = "utf-8"
	}
	blob, _, err := g.Git.CreateBlob(owner, repo, &github.Blob{
		Content:  content,
		Encoding: &enc,
	})

	tree, _, err := g.Git.CreateTree(owner, repo, tSHA, []github.TreeEntry{github.TreeEntry{
		Path: &path,
		Mode: &mode,
		SHA:  blob.SHA,
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
		Ref:    ref.Ref,
		URL:    ref.URL,
		Object: &github.GitObject{Type: ref.Object.Type, SHA: cSHAn, URL: ref.Object.URL},
	}, false)
	if err != nil {
		return err
	}

	return nil
}
