package sl

import (
	"context"
	"reflect"
)

// ProvideServiceFn is a delegat for initializing service.
type ProvideServiceFn func(ctx context.Context, c ServiceLocator) any

type ServiceLocator interface {
	OpenScope(ctx context.Context) context.Context
	CloseScope(ctx context.Context)
	AddSingleton(t reflect.Type, f ProvideServiceFn)
	AddTransient(t reflect.Type, f ProvideServiceFn)
	AddScoped(t reflect.Type, f ProvideServiceFn)
	GetService(ctx context.Context, t reflect.Type) any
}
