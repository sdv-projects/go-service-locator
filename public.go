package sl

import (
	"context"
	"reflect"
)

// GetService is a generic function for getting Service by type.
func GetService[T any](ctx context.Context, sl ServiceLocator) (T, error) {
	st := new(T)
	service := sl.GetService(ctx, reflect.TypeOf(st)).(T)

	return service, nil
}
