package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type AddDocsInput struct {
	Author     string
	Repository string
}

func (u *UseCase) AddRemoteDocs(input AddDocsInput) error {
	if input.Author == "" {
		return ErrAddRemoteDocsAuthorRequired
	}

	if input.Repository == "" {
		return ErrAddRemoteDocsRepositoryRequired
	}

	exists, err := u.docsRepository.Exists(input.Author, input.Repository)
	if err != nil {
		return err
	}

	if exists {
		return ErrAddRemoteDocsAlreadyExist
	}

	err = u.gitRemoteClient.IsDocsDirPresent(input.Author, input.Repository)
	if err != nil {
		return ErrAddRemoteDocsDirNotFound
	}

	docs := domain.NewRemoteDocs(
		input.Author,
		input.Repository,
		domain.RemoteSource(u.gitRemoteClient.Name()),
	)

	err = u.sourceRegistry.Download(*docs)
	if err != nil {
		return err
	}

	err = u.docsRepository.Save(*docs)
	if err != nil {
		return err
	}

	jobData, err := json.Marshal(map[string]string{
		"author":     input.Author,
		"repository": input.Repository,
		"source":     string(docs.Source),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal job payload: %w", err)
	}

	job, err := domain.NewJob(domain.JobTypeIndexPages, string(jobData))
	if err != nil {
		return fmt.Errorf("failed to create indexing job: %w", err)
	}

	err = u.jobRepository.Save(job)
	if err != nil {
		return fmt.Errorf("failed to save indexing job: %w", err)
	}

	return nil
}

var (
	ErrAddRemoteDocsAuthorRequired        = fmt.Errorf("author is required")
	ErrAddRemoteDocsRepositoryRequired    = fmt.Errorf("repository is required")
	ErrAddRemoteDocsAlreadyExist          = fmt.Errorf("docs already exist")
	ErrAddRemoteDocsDirNotFound           = fmt.Errorf("docs dir not found in repo")
	ErrAddRemoteDocsIndexingJobFailed     = fmt.Errorf("failed to create indexing job")
	ErrAddRemoteDocsIndexingJobSaveFailed = fmt.Errorf("failed to save indexing job")
)
