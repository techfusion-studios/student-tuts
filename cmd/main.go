package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/techfusion/school/student/pb/github.com/techfusion/student/v1"
	"github.com/techfusion/school/student/pkg/auth"
	"github.com/techfusion/school/student/pkg/config"
	"github.com/techfusion/school/student/pkg/database"
	"github.com/techfusion/school/student/pkg/datasource/repositories"
	"github.com/techfusion/school/student/pkg/datasource/services"
	"github.com/techfusion/school/student/pkg/middlewares"
	"github.com/techfusion/school/student/pkg/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	cfg = config.New()
	db  = database.New()
)

func init() {
	cfg.LoadEnv()
	db.Connect()
}

func main() {
	// Set up Keycloak OIDC configuration
	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		log.Fatalf("Failed to create OIDC authenticator: %v", err)
	}

	authMiddleware := middlewares.NewCustomGrpcAuthMiddleware(authenticator)

	// Create repository and service
	repo := repositories.NewStudentRepository(db.GetEngine())
	studentService := services.NewStudentService(repo)
	studentServer := servers.NewStudentServer(studentService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start REST server
	mux := runtime.NewServeMux()

	go func() {
		serverAddr := fmt.Sprintf(":%s", os.Getenv("APP.PORT"))
		if err := pb.RegisterStudentServiceHandlerServer(ctx, mux, studentServer); err != nil {
			log.Fatalf("Failed to register student service handler: %v", err)
		}

		log.Printf("REST controllers started on addr: %s\n", serverAddr)
		if err := http.ListenAndServe(serverAddr, mux); err != nil {
			log.Fatalln("failed to started REST controllers:", err)
		}

		log.Printf("REST endpoint started on addr: %s\n", serverAddr)
		if err := http.ListenAndServe(serverAddr, mux); err != nil {
			log.Fatalln("failed to started REST controllers:", err)
		}
	}()

	go func() {
		// create a standard HTTP router
		muxy := http.NewServeMux()

		// mount the gRPC HTTP gateway to the root
		muxy.Handle("/", mux)

		// mount a path to expose the generated OpenAPI specification on disk
		muxy.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./docs/swagger/protos/vendor.swagger.json")
		})

		// mount the Swagger UI that uses the OpenAPI specification path above
		muxy.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./swagger-ui"))))
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("APP.GRPC_PORT")))
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authMiddleware.CustomAuthInterceptor()))
	pb.RegisterStudentServiceServer(grpcServer, studentServer)
	reflection.Register(grpcServer)

	fmt.Printf("gRPC server is running on port :%s\n", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
