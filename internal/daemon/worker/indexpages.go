package worker

import (
	"log/slog"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type IndexPagesInput struct {
	Author     string
	Repository string
	Source     domain.RemoteSource
}

func (w *Worker) IndexPages(input IndexPagesInput) error {
	docs, err := w.docsRepository.GetOne(input.Author, input.Repository, input.Source)
	if err != nil {
		return err
	}

	files, err := w.sourceRegistry.GetAll(docs)
	if err != nil {
		return err
	}

	docs.StartIndexing()
	err = w.docsRepository.Save(docs)
	if err != nil {
		return err
	}

	defer func() {
		docs.StopIndexing()
		err := w.docsRepository.Save(docs)
		if err != nil {
			slog.Error("unable to save docs after indexing", "error", err)
		}
	}()

	for _, content := range files {
		docsPage, err := domain.NewDocsPage(docs.Author, docs.Repository, content.Text, content.RelPath)
		if err != nil {
			slog.Error("unable to create docs page", "error", err)
			continue
		}

		err = docsPage.SetLastModifiedAt(content.LastModifiedAt)
		if err != nil {
			slog.Error("unable to set last modified at", "error", err)
			continue
		}

		err = w.docsPageRepository.Save(*docsPage)
		if err != nil {
			slog.Error("unable to save docs page", "error", err)
			continue
		}
	}

	return nil
}
