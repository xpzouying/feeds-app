package middleware

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/sony/gobreaker"
)

func WithGoBreaker(ep endpoint.Endpoint, breakerName string) endpoint.Endpoint {

	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name: breakerName,
	})

	return circuitbreaker.Gobreaker(breaker)(ep)
}
