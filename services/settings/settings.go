package settings

import "github.com/bryopsida/gofiber-pug-starter/interfaces"

type settingsService struct {
	repo interfaces.ISettingsRepository
}

// NewSettingsService creates a new settingsService instance
func NewSettingsService(repo interfaces.ISettingsRepository) interfaces.ISettingsService {
	return &settingsService{repo: repo}
}

func (s *settingsService) GetString(key string) (string, error) {
	return s.repo.GetString(key)
}

func (s *settingsService) GetInt(key string) (int, error) {
	return s.repo.GetInt(key)
}

func (s *settingsService) GetBool(key string) (bool, error) {
	return s.repo.GetBool(key)
}

func (s *settingsService) SetString(key string, value string) error {
	return s.repo.Set(key, value)
}

func (s *settingsService) SetInt(key string, value int) error {
	return s.repo.Set(key, value)
}

func (s *settingsService) SetBool(key string, value bool) error {
	return s.repo.Set(key, value)
}
