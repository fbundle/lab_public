package monad

import "github.com/fbundle/lab_public/lab/go_util/pkg/option"

type Monad[S any, T any] struct {
	Init S
	Next func(S) (S, option.Option[T])
}
