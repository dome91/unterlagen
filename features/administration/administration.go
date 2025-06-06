package administration

type Administration struct {
	*settingsManager
	*users
}

func (*Administration) UserRoles() []UserRole {
	return []UserRole{
		UserRoleAdmin,
		UserRoleUser,
	}
}

func New(settingsRepository SettingsRepository, userRepository UserRepository, userMessages UserMessages) *Administration {
	return &Administration{
		settingsManager: newSettingsManager(settingsRepository),
		users:           newUsers(userRepository, userMessages),
	}
}
