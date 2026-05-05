package services

import (
	"planx/internal/repository"
	"planx/internal/models"
)

func ListEngineers(limit, offset int, specialization, city string) ([]models.Engineer, error) {
	return repository.GetEngineers(limit, offset, specialization, city)
}

func CreateEngineer(name, phone, email, specialization, city string) error {
	return repository.InsertEngineer(name, phone, email, specialization, city)
}