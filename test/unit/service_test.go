package unit

import (
	"testing"

	"github.com/Sphirium/learning-projects/wb-tech-demo-lo/internal/cache"
	"github.com/Sphirium/learning-projects/wb-tech-demo-lo/internal/models"
	"github.com/Sphirium/learning-projects/wb-tech-demo-lo/internal/repository"
	"github.com/Sphirium/learning-projects/wb-tech-demo-lo/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *MockRepo) FindByOrderUID(uid string) (*models.Order, error) {
	args := m.Called(uid)
	if result := args.Get(0); result != nil {
		return result.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRepo) GetAllOrderUIDs() ([]string, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockCache struct{ mock.Mock }

func (m *MockCache) Set(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *MockCache) Get(uid string) (*models.Order, error) {
	args := m.Called(uid)
	if result := args.Get(0); result != nil {
		return result.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Проверка соответствия интерфейсам
var _ repository.OrderRepositoryInterface = (*MockRepo)(nil)
var _ cache.OrderCacheInterface = (*MockCache)(nil)

func TestOrderService_GetOrderByUID_CacheHit(t *testing.T) {
	mockCache := new(MockCache)
	mockRepo := new(MockRepo)
	serv := service.NewOrderService(mockRepo, mockCache)

	expected := &models.Order{OrderUID: "test123"}
	mockCache.On("Get", "test123").Return(expected, nil)
	mockRepo.On("GetAllOrderUIDs").Maybe()

	order, err := serv.GetOrderByUID("test123")

	assert.NoError(t, err)
	assert.Equal(t, expected, order)
	mockRepo.AssertNotCalled(t, "FindByOrderUID")
	mockCache.AssertExpectations(t)
}
