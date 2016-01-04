package github_service

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	util "github.com/rjz/hubbub/common"
)

func findByPath(entries []github.TreeEntry, path string) *github.TreeEntry {
	for _, entry := range entries {
		if *entry.Path == path {
			return &entry
		}
	}
	return nil
}

type FileService struct {
	Client    *github.Client
	RepoOwner string
	RepoName  string
	RefTrees  map[sha]*github.Tree
}

func NewFileService(client *github.Client, owner, name string) *FileService {
	fs := FileService{client, owner, name, make(map[sha]*github.Tree)}
	return &fs
}

// RefFacts fetches the current state (SHA, tree) of the reference
func (fs *FileService) TreeFacts(SHA sha) error {
	isRecursive := false // TODO: support recursive paths
	tree, _, treeErr := fs.Client.Git.GetTree(fs.RepoOwner, fs.RepoName, string(SHA), isRecursive)
	if treeErr != nil {
		return treeErr
	}

	fs.RefTrees[SHA] = tree
	return nil
}

// CommitTree commits the provided tree
func (fs *FileService) CommitTree(tree *github.Tree, refName, parentSHA, msg string) error {
	commit, _, cErr := fs.Client.Git.CreateCommit(fs.RepoOwner, fs.RepoName, &github.Commit{
		Message: &msg,
		Tree:    &github.Tree{SHA: tree.SHA},
		Parents: []github.Commit{{SHA: &parentSHA}},
	})
	if cErr != nil {
		return cErr
	}

	newRef := github.Reference{
		Ref:    &refName,
		Object: &github.GitObject{SHA: commit.SHA},
	}
	_, _, err := fs.Client.Git.UpdateRef(fs.RepoOwner, fs.RepoName, &newRef, false)
	return err
}

// CreateOrUpdate updates an existing file or creates it if it does not exist.
// The new file conforms to the specified params.
func (fs *FileService) CreateOrUpdate(parentSHA sha, params fileParams) error {
	if fs.RefTrees[parentSHA] == nil {
		return errors.New(fmt.Sprintf("No tree available for SHA '%s'", parentSHA))
	}

	existingTree := fs.RefTrees[parentSHA]
	existingTreeEntries := existingTree.Entries
	filepath := *params.Name
	newEntries := append(existingTreeEntries, github.TreeEntry{
		Path:    &filepath,
		Mode:    util.String("100644"),
		Type:    util.String("blob"),
		Content: params.Content,
	})

	sha := existingTree.SHA // this might be wrong..
	// Create a new tree including the updated file to obtain a SHA
	newTree, _, tErr := fs.Client.Git.CreateTree(fs.RepoOwner, fs.RepoName, *sha, newEntries)
	if tErr != nil {
		return tErr
	}

	// Compare old and new SHAs to decide whether to update
	oldEntry := findByPath(existingTreeEntries, filepath)
	newEntry := findByPath(newTree.Entries, filepath)
	if oldEntry != nil {
		if *oldEntry.SHA == *newEntry.SHA {
			// nothing updated / nothing to do.
			return nil
		}
	}

	return fs.CommitTree(newTree, *params.Ref, string(parentSHA), fmt.Sprintf("Adding '%s'", filepath))
}

// Remove attempts to delete a file from the parent SHA
func (fs *FileService) Remove(parentSHA sha, params fileParams) error {
	if fs.RefTrees[parentSHA] == nil {
		return errors.New(fmt.Sprintf("No tree available for sha '%s'", parentSHA))
	}

	existingTree := fs.RefTrees[parentSHA]
	existingTreeEntries := existingTree.Entries
	filepath := *params.Name
	oldEntry := findByPath(existingTreeEntries, filepath)
	if oldEntry == nil {
		return nil
	}

	newEntries := []github.TreeEntry{}
	for _, e := range existingTreeEntries {
		if *e.Path != filepath {
			newEntries = append(newEntries, e)
		}
	}

	// Create a new tree including the updated file to obtain a SHA
	sha := existingTree.SHA // this might be wrong..
	newTree, _, tErr := fs.Client.Git.CreateTree(fs.RepoOwner, fs.RepoName, *sha, newEntries)
	if tErr != nil {
		return tErr
	}

	return fs.CommitTree(newTree, *params.Ref, string(parentSHA), fmt.Sprintf("Adding '%s'", filepath))
}
