package administration

import "unterlagen/features/common"

type Administration struct {
	*settingsManager
	*users
	taskRepository common.TaskRepository
}

func (*Administration) UserRoles() []UserRole {
	return []UserRole{
		UserRoleAdmin,
		UserRoleUser,
	}
}

func (a *Administration) GetAllTasks() ([]common.Task, error) {
	return a.taskRepository.FindAll()
}

func (a *Administration) GetTasksPaginated(page int) ([]common.Task, int, int, error) {
	limit := 10
	offset := (page - 1) * limit
	tasks, total, err := a.taskRepository.FindPaginated(limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	
	totalPages := (total + limit - 1) / limit
	return tasks, total, totalPages, nil
}

func New(settingsRepository SettingsRepository, userRepository UserRepository, userMessages UserMessages, taskRepository common.TaskRepository) *Administration {
	return &Administration{
		settingsManager: newSettingsManager(settingsRepository),
		users:           newUsers(userRepository, userMessages),
		taskRepository:  taskRepository,
	}
}
