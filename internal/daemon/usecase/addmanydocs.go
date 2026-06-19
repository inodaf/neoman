package usecase

import (
	"fmt"
	"sync"
	"time"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type AddManyDocsInput struct {
	Limit  uint
	Author string
}

func (u *UseCase) AddManyRemoteDocs(input AddManyDocsInput) error {
	if input.Author == "" {
		return ErrAddRemoteDocsAuthorRequired
	}

	repos, err := u.gitHost.ListActiveRepos(input.Author, input.Limit)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		fmt.Printf("Checking %s/%s...\n", repo.Name, repo.UpdatedAt.Format(time.DateTime))
	}
	return nil
	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Go(func () {		
			exists, err := u.docsRepository.Exists(input.Author, repo.Name)
			if err != nil {
				return
			}
			
			if exists {
				return
			}
	
			availability, err := u.gitHost.HasDocs(input.Author, []string{repo.Name})
			if err != nil || availability[repo.Name] == false {
				return
			}
	
			docs := domain.NewRemoteDocs(
				input.Author,
				repo.Name,
				domain.RemoteSource(u.gitHost.Name()),
			)
	
			err = u.sourceRegistry.Download(*docs)
			if err != nil {
				return
			}
	
			err = u.docsRepository.Save(*docs)
			if err != nil {
				return
			}
		})
	}

	wg.Wait()
	
	return nil
}
