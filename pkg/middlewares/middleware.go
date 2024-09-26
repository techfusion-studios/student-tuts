package middlewares

import (
	"context"
	"github.com/techfusion/school/student/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomGrpcAuthMiddleware interface {
	CustomAuthInterceptor() grpc.UnaryServerInterceptor
}

type customGrpcAuthMiddleware struct {
	authenticator auth.AuthenticatorService
}

// CustomAuthInterceptor selectively applies authentication only to private methods
func (middleware *customGrpcAuthMiddleware) CustomAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Define public methods that don't require authentication
		publicMethods := map[string]bool{
			"/protos.techfusion.student.v1.StudentService/ListStudent": true,
			"/protos.techfusion.student.v1.StudentService/GetStudent":  true,
		}

		// If the method is public, skip authentication
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Otherwise, apply the authentication middleware
		token, err := middleware.authenticator.ExtractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// Validate the token
		idToken, err := middleware.authenticator.VerifyToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// Add user info to context
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Inject claims into context
		ctx = context.WithValue(ctx, auth.ContextKeyUser, claims)

		// Proceed with the handler
		return handler(ctx, req)
	}
}

func NewCustomGrpcAuthMiddleware(authenticator auth.AuthenticatorService) CustomGrpcAuthMiddleware {
	return &customGrpcAuthMiddleware{authenticator: authenticator}
}
