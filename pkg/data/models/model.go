package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	pb "github.com/techfusion/school/student/pb/github.com/techfusion/student/v1"
	"gorm.io/gorm"
	"time"
)

type Student struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	Name         string
	Age          int
	Email        string `gorm:"unique"`
	RegisteredBy string
	RegisteredAt *time.Time
}

func (s *Student) ToProto() *pb.Student {
	p := new(pb.Student)
	_ = copier.Copy(p, s)
	return p
}

func NewStudent(name string, age int, email string) *Student {
	return &Student{
		ID:    uuid.NewString(),
		Name:  name,
		Age:   age,
		Email: email,
	}
}
