package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ClaimsSymmetric Набор утвержданий для симметричных межмикросервисных токенов
type ClaimsSymmetric struct {
	jwt.RegisteredClaims
	ServiceName string `json:"servicename"` // имя сервиса, который делает запрос
}

// ServiceSymmetric -симметричное шифрование (у клиента и сервера один секретный ключ,
// взаимодействуют сервисы из определенного списка)
type ServiceSymmetric struct {
	serviceName       string
	validServicesList []string      // список названий сервисов, от которых принимаем запросы
	secret            string        // секрет для проверки и генерации подписи
	validityPeriod    time.Duration // Период действия в минутах
}

func NewServiceSymmetric(serviceName string, validServicesList []string, secret string) *ServiceSymmetric {
	s := &ServiceSymmetric{
		serviceName:       serviceName,
		validServicesList: validServicesList,
		secret:            secret,
		validityPeriod:    60, // 60 минут (TODO вынести в конфиг)
	}

	return s
}

// ServiceName Получить имя JWT сервиса
func (s *ServiceSymmetric) ServiceName() string {
	return s.serviceName
}

// GenerateJWT создает JWT токен
func (s *ServiceSymmetric) GenerateJWT() (string, error) {
	claims := ClaimsSymmetric{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.validityPeriod * time.Minute)),
		},
		ServiceName: s.serviceName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}

// GetClaims получение утверждений из токена
func (s *ServiceSymmetric) GetClaims(tokenStr string) (*ClaimsSymmetric, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsSymmetric{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClaimsSymmetric)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// IsServiceValid проверка валидности сервиса от которого идет запрос на выполнение удаленной процедуры
func (s *ServiceSymmetric) IsServiceValid(tokenStr string) (bool, error) {
	claims, err := s.GetClaims(tokenStr)
	if err != nil {
		return false, err
	}

	return s.validateServiceName(claims.ServiceName), nil
}

// validateServiceNAme Проверить есть ли имя сервиса в списке валидных (могут ли запускаться процедуры от имени этого сервиса)
func (s *ServiceSymmetric) validateServiceName(serviceName string) bool {
	for _, v := range s.validServicesList {
		if v == serviceName {
			return true
		}
	}
	return false
}
