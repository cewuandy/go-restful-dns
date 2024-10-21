package domain

import "context"

type InitHandler interface {
	Initialize(ctx context.Context) error
}

type InitUseCase interface {
	ClearRedisData(ctx context.Context) error

	RecoverRecords(ctx context.Context) error
}
