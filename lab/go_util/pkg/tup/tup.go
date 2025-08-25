package tup

type Tup[T1 any, T2 any] struct {
	fst T1
	snd T2
}

func MakeTup[T1 any, T2 any](v1 T1, v2 T2) Tup[T1, T2] {
	return Tup[T1, T2]{
		fst: v1,
		snd: v2,
	}
}
