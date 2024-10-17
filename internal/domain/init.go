package domain

import "context"

type InitHandler interface {
	Initialize(ctx context.Context) error
}

type InitUseCase interface {
	RecoverRecords(ctx context.Context) error
}
