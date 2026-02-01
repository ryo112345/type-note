package port

import "context"

// TxManager はトランザクション境界を制御するインターフェース。
type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
