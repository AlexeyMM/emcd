package log

import "context"

func CreateClientNamedContext(ctx context.Context, name string) context.Context {

	return context.WithValue(ctx, serviceNameStruct{}, name)
}
