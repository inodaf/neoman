package ports

type UserPrompter interface {
	ConfirmTrust(authorName string) (bool, error)
}
