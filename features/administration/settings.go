package administration

import "errors"

var ErrSettingsNotFound = errors.New("settings not found")

type Settings struct {
	ID                  string
	RegistrationEnabled bool
}

func defaultSettings() Settings {
	return Settings{
		ID:                  "settings",
		RegistrationEnabled: false,
	}
}

type SettingsRepository interface {
	Save(settings Settings) error
	Get() (Settings, error)
}

type settingsManager struct {
	repository SettingsRepository
}

func (manager *settingsManager) Get() (Settings, error) {
	settings, err := manager.repository.Get()
	if errors.Is(err, ErrSettingsNotFound) {
		settings = defaultSettings()
		err = manager.repository.Save(settings)
		if err != nil {
			return Settings{}, err
		}
	}
	return settings, err
}

func newSettingsManager(repository SettingsRepository) *settingsManager {
	return &settingsManager{
		repository: repository,
	}
}
