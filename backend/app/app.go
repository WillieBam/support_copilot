package app

type App struct {
	Client     *appClient
	Repository *appRepository
	Service    *appService
}

func NewApp() *App {
	appClient := newAppClient()
	appRepository := newAppRepository(appClient)
	appService := newAppService(appClient, appRepository)

	return &App{
		Client:     appClient,
		Repository: appRepository,
		Service:    appService,
	}
}
