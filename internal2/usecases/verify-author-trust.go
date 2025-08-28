package usecases

import "github.com/inodaf/neoman/internal2/ports"

type VerifyAuthorTrustUseCase struct {
	AuthorService ports.AuthorService
	UserPrompter  ports.UserPrompter
}

func (u *VerifyAuthorTrustUseCase) Execute(authorName string) (bool, error) {
	if u.AuthorService.IsTrusted(authorName) {
		return true, nil
	}

	confirmed, err := u.UserPrompter.ConfirmTrust(authorName)
	if err != nil {
		return false, err
	}

	if confirmed {
		err := u.AuthorService.Trust(authorName)
		if err != nil {
			return false, err
		}
	}

	return confirmed, nil
}
