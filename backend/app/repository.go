package app

type appRepository struct {
	client *appClient
}

func newAppRepository(client *appClient) *appRepository {
	return &appRepository{client: client}
}
