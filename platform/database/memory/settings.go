package memory

import (
	"unterlagen/features/administration"
)

var _ administration.SettingsRepository = &SettingsRepository{}

type SettingsRepository struct {
	settings administration.Settings
}

// Get implements administration.SettingsRepository.
func (s *SettingsRepository) Get() (administration.Settings, error) {
	return s.settings, nil
}

// Save implements administration.SettingsRepository.
func (s *SettingsRepository) Save(settings administration.Settings) error {
	s.settings = settings
	return nil
}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{
		settings: administration.Settings{},
	}
}
