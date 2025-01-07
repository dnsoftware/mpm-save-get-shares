package grpc

type JWTProcessor interface {
	GetActualToken() (string, error)
}
