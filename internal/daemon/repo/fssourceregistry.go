package repo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/inodaf/neoman/internal/daemon/domain"
	"github.com/inodaf/neoman/pkg/config"
	"github.com/inodaf/neoman/pkg/git"
)

const MAX_FILE_SIZE = 1024 * 1024 * 10 // 10MB

// isMarkdownFile checks if the filename has a markdown extension (.md or .mdx)
// using case-insensitive matching.
func isMarkdownFile(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".mdx")
}

// isHiddenFile checks if the filename starts with a dot.
func isHiddenFile(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

func NewFsSourceRegistry(remote git.GitRemote) SourceRegistry {
	return &FsSourceRegistry{GitRemoteProvider: remote}
}

type FsSourceRegistry struct {
	GitRemoteProvider git.GitRemote
}

// Download optimimally clones the remote repository into into the registry directory.
func (r *FsSourceRegistry) Download(docs domain.RemoteDocs) error {
	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		return err
	}

	authorDir := path.Join(registryDir, "remote", docs.Author)
	if _, err := os.Stat(authorDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(authorDir, os.ModePerm); err != nil {
			return err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory")
	}

	err = os.Chdir(authorDir)
	if err != nil {
		return fmt.Errorf("could not change directory")
	}

	defer os.Chdir(wd)

	err = git.Clone(docs.Author, docs.Repository, r.GitRemoteProvider)
	if err != nil {
		return fmt.Errorf("could not clone repository from remote")
	}

	err = os.Chdir(path.Join(authorDir, docs.Repository))
	if err != nil {
		return fmt.Errorf("could not access cloned repository directory")
	}

	return git.SparseCheckout()
}

// GetAllPaths returns all Markdown (.md, .mdx) file paths from the remote documentation.
func (r *FsSourceRegistry) GetAllPaths(docs domain.RemoteDocs) ([]string, error) {
	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		return nil, err
	}

	docsDir := path.Join(registryDir, "remote", docs.Author, docs.Repository, config.PrimaryDocsDirName)
	if _, err := os.Stat(docsDir); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("docs does not exist")
	}

	var matches []string
	err = filepath.WalkDir(docsDir, func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || isHiddenFile(entry.Name()) {
			return nil
		}
		if isMarkdownFile(entry.Name()) {
			matches = append(matches, filePath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not locate documentation files: %w", err)
	}

	return matches, nil
}

// GetAllContents reads the content of all Markdown (.md, .mdx) files in the remote
// documentation and returns them as a slice of RegistryContent.
func (r *FsSourceRegistry) GetAllContents(docs domain.RemoteDocs) ([]RegistryContent, error) {
	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		return nil, err
	}

	docsDir := path.Join(registryDir, "remote", docs.Author, docs.Repository, config.PrimaryDocsDirName)

	matches, err := r.GetAllPaths(docs)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(matches))

	contents := make([]RegistryContent, 0, len(matches))
	for _, docPath := range matches {
		go func(docPath string) {
			defer wg.Done()

			file, err := os.Open(docPath)
			if err != nil {
				slog.Error("unable to open doc file", "path", docPath)
				return
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				slog.Error("unable to get doc file info", "path", docPath)
				return
			}

			if fileInfo.Size() >= MAX_FILE_SIZE {
				slog.Error("doc file is too large", "path", docPath)
				return
			}

			reader := bufio.NewReader(file)
			content, err := io.ReadAll(reader)
			if err != nil {
				slog.Error("unable to read doc file", "path", docPath)
				return
			}

			relPath, err := filepath.Rel(docsDir, docPath)
			if err != nil {
				slog.Error("unable to get doc file relative path", "path", docPath)
				return
			}

			contents = append(contents, RegistryContent{
				Text:           string(content),
				RelPath:        relPath,
				LastModifiedAt: fileInfo.ModTime(),
			})
		}(docPath)
	}

	wg.Wait()

	return contents, nil
}
