package servers

import (
	"context"
	pb "github.com/techfusion/school/student/pb/github.com/techfusion/student/v1"
	"github.com/techfusion/school/student/pkg/datasource/services"
)

type studentServer struct {
	pb.UnimplementedStudentServiceServer
	studentService services.StudentService
}

func (server *studentServer) ListStudent(ctx context.Context, req *pb.ListStudentRequest) (*pb.ListStudentResponse, error) {
	studentResponse, err := server.studentService.GetAllStudents(ctx, req)
	if err != nil {
		return nil, err
	}

	return studentResponse, nil
}

func (server *studentServer) CreateStudent(ctx context.Context, request *pb.CreateStudentRequest) (*pb.CreateStudentResponse, error) {
	studentResponse, err := server.studentService.CreateStudent(ctx, request)
	if err != nil {
		return nil, err
	}

	return studentResponse, nil
}

func (server *studentServer) GetStudent(ctx context.Context, request *pb.GetStudentRequest) (*pb.GetStudentResponse, error) {
	studentResponse, err := server.studentService.GetStudent(ctx, request)
	if err != nil {
		return nil, err
	}

	return studentResponse, nil
}

func (server *studentServer) UpdateStudent(ctx context.Context, request *pb.UpdateStudentRequest) (*pb.UpdateStudentResponse, error) {
	studentResponse, err := server.studentService.UpdateStudent(ctx, request)
	if err != nil {
		return nil, err
	}

	return studentResponse, nil
}

func (server *studentServer) DeleteStudent(ctx context.Context, request *pb.DeleteStudentRequest) (*pb.DeleteStudentResponse, error) {
	studentResponse, err := server.studentService.DeleteStudent(ctx, request)
	if err != nil {
		return nil, err
	}

	return studentResponse, nil
}

func NewStudentServer(studentService services.StudentService) pb.StudentServiceServer {
	return &studentServer{
		studentService: studentService,
	}
}
