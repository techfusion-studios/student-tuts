package repositories

import (
	"github.com/techfusion/school/student/pkg/data/models"
	"gorm.io/gorm"
)

// StudentRepository defines the methods for CRUD operations on the Student model.
type StudentRepository interface {
	CreateStudent(student *models.Student) error
	GetStudent(id string) (*models.Student, error)
	UpdateStudent(student *models.Student) error
	DeleteStudent(*models.Student) error
	GetAll() ([]*models.Student, error)
}

type studentRepository struct {
	store *gorm.DB
}

func (repository *studentRepository) GetAll() ([]*models.Student, error) {
	var students []*models.Student
	err := repository.store.Find(&students).Error
	return students, err
}

// NewStudentRepository initializes a new StudentRepository.
func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{store: db}
}

func (repository *studentRepository) CreateStudent(student *models.Student) error {
	if err := repository.store.Create(&student).Error; err != nil {
		return err
	}

	return nil
}

func (repository *studentRepository) GetStudent(id string) (*models.Student, error) {
	var student models.Student
	if err := repository.store.First(&student, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (repository *studentRepository) UpdateStudent(student *models.Student) error {
	if err := repository.store.Save(&student).Error; err != nil {
		return err
	}

	return nil
}

func (repository *studentRepository) DeleteStudent(s *models.Student) error {
	if err := repository.store.Delete(s).Error; err != nil {
		return err
	}
	return nil
}
