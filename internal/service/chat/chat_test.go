package chat

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/AlGrushino/chat/internal/repository"
	"github.com/AlGrushino/chat/internal/repository/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) Create(ctx context.Context, chat *models.Chat) error {
	args := m.Called(ctx, chat)

	if chat != nil && args.Error(0) == nil {
		if chat.ID == 0 {
			chat.ID = 1
		}
	}

	return args.Error(0)
}

func (m *MockChatRepository) GetByID(ctx context.Context, id int) (*models.Chat, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Chat), args.Error(1)
}

func (m *MockChatRepository) ChatExists(ctx context.Context, title string) (bool, error) {
	args := m.Called(ctx, title)
	return args.Bool(0), args.Error(1)
}

func (m *MockChatRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Chat, error) {
	args := m.Called(ctx, limit, offset)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Chat), args.Error(1)
}

func (m *MockChatRepository) Update(ctx context.Context, chat *models.Chat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

func (m *MockChatRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockLogger struct {
	messages []string
	errors   []string
	warns    []string
	infos    []string
}

func (m *MockLogger) WithError(err error) *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}

func (m *MockLogger) WithField(key string, value any) *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}

func (m *MockLogger) Warnf(format string, args ...any) {
	m.warns = append(m.warns, format)
}

func (m *MockLogger) Error(msg ...any) {
	m.errors = append(m.errors, fmt.Sprint(msg...))
}

func (m *MockLogger) Infof(format string, args ...any) {
	m.infos = append(m.infos, fmt.Sprintf(format, args...))
}

func (m *MockLogger) ContainsWarn(msg string) bool {
	for _, w := range m.warns {
		if strings.Contains(w, msg) {
			return true
		}
	}
	return false
}

type ChatServiceTestSuite struct {
	suite.Suite
	ctx        context.Context
	mockRepo   *MockChatRepository
	mockLogger *logrus.Logger
	service    *ChatService
}

func (suite *ChatServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockRepo = new(MockChatRepository)
	suite.mockLogger = logrus.New()
	suite.service = NewChatService(suite.mockLogger, &repository.Repository{
		Chat: suite.mockRepo,
	})
}

func (suite *ChatServiceTestSuite) TestCreateChat_Success() {
	expectedTitle := "Успешное создание"
	suite.mockRepo.On("ChatExists", suite.ctx, expectedTitle).Return(false, nil).Once()
	suite.mockRepo.On("Create", suite.ctx, mock.MatchedBy(func(chat *models.Chat) bool {
		return chat.Title == expectedTitle
	})).Return(nil).Once()

	result, err := suite.service.CreateChat(suite.ctx, expectedTitle)

	suite.NoError(err)
	suite.Equal(expectedTitle, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ChatServiceTestSuite) TestCreateChat_Duplicate() {
	duplicateTitle := "Дубликат"
	suite.mockRepo.On("ChatExists", suite.ctx, duplicateTitle).Return(true, nil).Once()

	result, err := suite.service.CreateChat(suite.ctx, duplicateTitle)

	suite.Error(err)
	suite.Contains(err.Error(), "already exists")
	suite.Empty(result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ChatServiceTestSuite) TestCreateChat_DatabaseError() {
	title := "Проблемный"
	suite.mockRepo.On("ChatExists", suite.ctx, title).Return(false, nil).Once()
	suite.mockRepo.On("Create", suite.ctx, mock.Anything).
		Return(errors.New("database error")).
		Once()

	result, err := suite.service.CreateChat(suite.ctx, title)

	suite.Error(err)
	suite.Contains(err.Error(), "failed to create chat")
	suite.Empty(result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestChatServiceSuite(t *testing.T) {
	suite.Run(t, new(ChatServiceTestSuite))
}
