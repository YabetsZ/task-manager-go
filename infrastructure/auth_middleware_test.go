package infrastructure_test

import (
	"net/http"
	"net/http/httptest"
	"task-manager/domain"
	"task-manager/errs"
	"task-manager/infrastructure"
	"task-manager/usecases/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	mockUserUsecase *mocks.UserUsecase
	jwtService      infrastructure.JWTServiceV5
	router          *gin.Engine
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockUserUsecase = new(mocks.UserUsecase)
	s.jwtService = *infrastructure.NewJWTServiceV5()
	s.router = gin.Default()
}

func TestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (s *AuthMiddlewareTestSuite) performRequestWithAuth(header string, requiredRole string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	if header != "" {
		req.Header.Set("Authorization", header)
	}

	s.router.GET("/protected", infrastructure.AuthMiddleware(s.mockUserUsecase, requiredRole), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s.router.ServeHTTP(w, req)
	return w
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_NoHeader() {
	w := s.performRequestWithAuth("", domain.RoleUser)
	s.Assert().Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_MalformedHeader() {
	w := s.performRequestWithAuth("InvalidToken", domain.RoleUser)
	s.Assert().Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	w := s.performRequestWithAuth("Bearer invalid-token-string", domain.RoleUser)
	s.Assert().Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_UserNotFound() {
	user := &domain.User{ID: "123", Username: "test", Role: domain.RoleUser}
	token, _ := s.jwtService.GenerateJWT(user)

	s.mockUserUsecase.On("GetUserByID", user.ID).Return(nil, errs.ErrUserNotFound).Once()

	w := s.performRequestWithAuth("Bearer "+token, domain.RoleUser)

	s.Assert().Equal(http.StatusUnauthorized, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_InsufficientPermissions() {
	user := &domain.User{ID: "123", Username: "test", Role: domain.RoleUser}
	token, _ := s.jwtService.GenerateJWT(user)

	s.mockUserUsecase.On("GetUserByID", user.ID).Return(user, nil).Once()

	w := s.performRequestWithAuth("Bearer "+token, domain.RoleAdmin)

	s.Assert().Equal(http.StatusForbidden, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_AdminAccessingUserRoute_Success() {
	admin := &domain.User{ID: "456", Username: "admin", Role: domain.RoleAdmin}
	token, _ := s.jwtService.GenerateJWT(admin)

	s.mockUserUsecase.On("GetUserByID", admin.ID).Return(admin, nil).Once()

	w := s.performRequestWithAuth("Bearer "+token, domain.RoleUser)

	s.Assert().Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_Success() {
	user := &domain.User{ID: "123", Username: "test", Role: domain.RoleUser}
	token, _ := s.jwtService.GenerateJWT(user)

	s.mockUserUsecase.On("GetUserByID", user.ID).Return(user, nil).Once()

	w := s.performRequestWithAuth("Bearer "+token, domain.RoleUser)

	s.Assert().Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}
