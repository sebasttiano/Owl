package service

type Repository interface {
}

type AuthService struct {
	Repo *Repository
}

func NewAuthService(repo *Repository) *AuthService {
	return &AuthService{Repo: repo}
}

type BinaryService struct {
	Repo *Repository
}

func NewBinaryService(repo *Repository) *BinaryService {
	return &BinaryService{Repo: repo}
}

type TextService struct {
	Repo *Repository
}

func NewTextService(repo *Repository) *TextService {
	return &TextService{Repo: repo}
}
