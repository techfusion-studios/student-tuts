package services

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	pb "github.com/techfusion/school/student/pb/github.com/techfusion/student/v1"
	"github.com/techfusion/school/student/pkg/auth"
	"github.com/techfusion/school/student/pkg/data/models"
	"github.com/techfusion/school/student/pkg/datasource/repositories"
	"log"
)

type StudentService interface {
	CreateStudent(context.Context, *pb.CreateStudentRequest) (*pb.CreateStudentResponse, error)
	DeleteStudent(context.Context, *pb.DeleteStudentRequest) (*pb.DeleteStudentResponse, error)
	GetStudent(context.Context, *pb.GetStudentRequest) (*pb.GetStudentResponse, error)
	UpdateStudent(context.Context, *pb.UpdateStudentRequest) (*pb.UpdateStudentResponse, error)
	GetAllStudents(context.Context, *pb.ListStudentRequest) (*pb.ListStudentResponse, error)
}

// studentService implements the generated gRPC service
type studentService struct {
	repository repositories.StudentRepository
}

func (service *studentService) GetAllStudents(ctx context.Context, request *pb.ListStudentRequest) (*pb.ListStudentResponse, error) {
	students, err := service.repository.GetAll()
	if err != nil {
		return nil, err
	}

	studentProtos := make([]*pb.Student, 0, len(students))
	for _, student := range students {
		studentProtos = append(studentProtos, student.ToProto())
	}

	return &pb.ListStudentResponse{Students: studentProtos}, nil
}

// NewStudentService initializes a new studentService
func NewStudentService(studentRepository repositories.StudentRepository) StudentService {
	return &studentService{repository: studentRepository}
}

func (service *studentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentResponse, error) {
	// Retrieve claims from context to ensure the user is authenticated
	userClaims, ok := ctx.Value(auth.ContextKeyUser).(map[string]interface{})
	if !ok {
		return nil, errors.New("unauthenticated")
	}

	log.Println("UserClaims ::: |", userClaims)

	studentRequest := req.GetStudent()
	student := models.NewStudent(studentRequest.GetName(), int(studentRequest.GetAge()), studentRequest.GetEmail())
	if err := service.repository.CreateStudent(student); err != nil {
		return nil, err
	}

	return &pb.CreateStudentResponse{Student: student.ToProto()}, nil
}

func (service *studentService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentResponse, error) {
	student, err := service.repository.GetStudent(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetStudentResponse{Student: student.ToProto()}, nil
}

func (service *studentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentResponse, error) {
	student, err := service.repository.GetStudent(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(student, req.GetStudent()); err != nil {
		return nil, err
	}

	if err := service.repository.UpdateStudent(student); err != nil {
		return nil, err
	}

	return &pb.UpdateStudentResponse{Student: student.ToProto()}, nil
}

func (service *studentService) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentResponse, error) {
	student, err := service.repository.GetStudent(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := service.repository.DeleteStudent(student); err != nil {
		return nil, err
	}

	return &pb.DeleteStudentResponse{Message: "Student deleted successfully"}, nil
}
