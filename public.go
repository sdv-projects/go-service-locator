package sl

import (
	"context"
	"reflect"
)

// GetService is a generic function for getting Service by type.
func GetService[T any](ctx context.Context, sl ServiceLocator) (T, error) {
	var res T
	res = sl.GetService(ctx, reflect.TypeOf(res)).(T)

	return res, nil
}
