package git

import (
	"encoding/base64"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	owner    = `gamedevja`
	repo     = `gamedevja.github.io`
	head     = `heads/master`
	treetype = `blob`
	mode     = `100644`
)

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

var (
	BlobBase64Encode = "base64"
	BlobUtf8Encode   = "utf-8"
)

type Blob struct {
	Blob github.Blob
	Path string
}

func Push(blobs []Blob, message *string) error {
	var err error

	ts := &tokenSource{
		&oauth2.Token{AccessToken: os.Getenv("GITHUBTOKEN")},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	g := github.NewClient(tc)
	ref, _, _ := g.Git.GetRef(owner, repo, head)
	cSHA := *(ref.Object.SHA)

	com, _, _ := g.Git.GetCommit(owner, repo, cSHA)

	tree := com.Tree
	for _, b := range blobs {
		if *b.Blob.Encoding == BlobBase64Encode {
			b64 := base64.StdEncoding.EncodeToString([]byte(*b.Blob.Content))
			b.Blob.Content = &b64
		}
		if blob, _, err := g.Git.CreateBlob(owner, repo, &b.Blob); err != nil {
			return err
		} else {
			if t, _, err := g.Git.CreateTree(owner, repo, *tree.SHA, []github.TreeEntry{github.TreeEntry{
				Path: &b.Path,
				Mode: &mode,
				Type: &treetype,
				SHA:  blob.SHA,
			}}); err != nil {
				return err
			} else {
				tree = t
			}
		}
	}

	comn, _, _ := g.Git.CreateCommit(owner, repo, &github.Commit{
		Message: message,
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
