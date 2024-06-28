package service

type Repository interface {
}

type ServiceSettings struct {
}

//type ServicePools struct {
//	MainPool,
//	AwaitPool worker.Pool
//}

type Service struct {
	Repo     Repository
	settings *ServiceSettings
	//pools    *ServicePools
}

func NewService(repo Repository, settings *ServiceSettings) *Service {
	return &Service{
		Repo:     repo,
		settings: settings,
	}
}
