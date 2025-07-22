package main

import (
	"encoding/json"
	"fmt"
	"go_util/pkg/fib"
	"go_util/pkg/int_ntt"
	"go_util/pkg/integer"
	"go_util/pkg/line_slice"
	"go_util/pkg/monad"
	"go_util/pkg/padic"
	"go_util/pkg/persistent/ordered_map"
	"go_util/pkg/persistent/vector"
	"go_util/pkg/ring"
	"go_util/pkg/tup"
	"go_util/pkg/vec"
	"math/rand"
	"os"
	"strconv"
)

func testPAdic() {
	const N = 10
	x := padic.NewPAdicFromInt(3, 23)
	y := padic.NewPAdicFromInt(3, 27)
	z := padic.NewPAdicFromInt(3, 92)
	fmt.Println(x.Approx(N))
	fmt.Println(y.Approx(N))
	fmt.Println(y.Mul(x).Approx(N)) // print 27 x 23 = 621
	fmt.Println(y.Sub(x).Approx(N)) // print 27 - 23 = 4
	fmt.Println(z.Div(x).Approx(N)) // print 92 / 23 = 4
}

func testFib() {
	if len(os.Args) < 2 {
		panic("Usage: go run main.go <integer>")
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Error: argument must be an integer")
		return
	}

	fmt.Println(fib.Fib(int_ntt.Nat{}, uint64(n)))
}

func testEA() {
	a, b := ring.EuclideanAlgorithm(integer.FromInt64(15), integer.FromInt64(46))
	fmt.Println(a, b)
}

func testIntNTT() {
	x := int_ntt.FromString("0x318346193417412890342342")
	z := int_ntt.FromString("0x539543980a084524")
	fmt.Println(z)
	fmt.Println(z.Add(x).Mod(x))
}

func testTup() {
	t := tup.MakeTup(2, 2.5)
	fmt.Println(t)
}

func testVecFunctor() {
	add1 := func(x int) int {
		return x + 1
	}

	v1 := vec.MakeVecFromSlice([]int{1, 2, 3, 4, 5})

	add1vec := vec.Wrap(add1)
	v2 := add1vec(v1)
	fmt.Println(v2)

}

func testMonad() {
	a := monad.None[uint]().Insert(1, 2, 3, 4)
	join := func(tr string, t uint) (string, bool) {
		if t >= 4 {
			return "", false
		}
		return fmt.Sprintf("%s%d,", tr, t), true
	}
	resultList := []interface{}{
		monad.Natural.TakeAtMost(10).Slice(),
		monad.Filter(monad.Natural, func(n uint) bool {
			return n%2 == 0
		}).TakeAtMost(10).Slice(),
		monad.Replicate(5).TakeAtMost(10).Slice(),
		a.TakeAtMost(0).Slice(),
		a.TakeAtMost(2).Slice(),
		a.TakeAtMost(5).Slice(),
		a.Insert(9, 8, 7).Slice(),
		a.DropAtMost(0).Slice(),
		a.DropAtMost(2).Slice(),
		a.DropAtMost(5).Slice(),
		monad.Map(a, func(x uint) uint {
			return x * 2
		}).Slice(),
		monad.Filter(a, func(n uint) bool {
			return n%2 == 0
		}).Slice(),
		monad.Reduce(a, join, ""),
		monad.Fold(a, join, "").Slice(),
		monad.Bind(a, func(ta uint) monad.Monad[uint] {
			return monad.Replicate(ta).TakeAtMost(int(ta))
		}).Slice(),
		monad.Fibonacci.TakeAtMost(10).Slice(),
		monad.Prime.TakeAtMost(10).Slice(),
	}
	for _, result := range resultList {
		fmt.Println(result)
	}
}

type Item struct {
	Id    uint   `json:"id"`
	Value string `json:"value"`
}

func (i Item) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func setupTmpDir(dir string) {
	if _, err := os.Stat(dir); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
	} else {
		panic(err)
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
}

func testLineSlice() {
	setupTmpDir("tmp")
	openLS := func() (line_slice.LineSlice[Item], error) {
		return line_slice.NewLineSlice[Item](
			"tmp/test.jsonl",
			func(b []byte) (Item, error) {
				item := Item{}
				err := json.Unmarshal(b, &item)
				if err != nil {
					return Item{}, err
				}
				return item, err
			},
			func(item Item) ([]byte, error) {
				b, err := json.Marshal(item)
				if err != nil {
					return nil, err
				}
				return b, nil
			},
			'\n',
		)
	}
	{
		ls, err := openLS()
		if err != nil {
			panic(err)
		}
		defer ls.Close()

		//
		err = ls.Push(Item{Id: 0, Value: "zero"})
		if err != nil {
			panic(err)
		}
		err = ls.Push(Item{Id: 1, Value: "one"})
		if err != nil {
			panic(err)
		}
		err = ls.Push(Item{Id: 2, Value: "two"})
		if err != nil {
			panic(err)
		}
		fmt.Println(ls.Get(0))
		fmt.Println(ls.Get(1))
		fmt.Println(ls.Get(2))
	}
	{
		ls, err := openLS()
		if err != nil {
			panic(err)
		}
		defer ls.Close()

		//
		err = ls.Push(Item{Id: 3, Value: "three"})
		if err != nil {
			panic(err)
		}
		err = ls.Push(Item{Id: 4, Value: "four"})
		if err != nil {
			panic(err)
		}
		err = ls.Push(Item{Id: 5, Value: "five"})
		if err != nil {
			panic(err)
		}
		fmt.Println(ls.Get(0))
		fmt.Println(ls.Get(1))
		fmt.Println(ls.Get(2))
		fmt.Println(ls.Get(3))
		fmt.Println(ls.Get(4))
		fmt.Println(ls.Get(5))
	}
}

func testPersistentOrderedMap() {
	w := ordered_map.NewOrderedMap[int, struct{}]()
	fmt.Println(w.Repr())
	w = w.
		Set(10, struct{}{}).
		Set(11, struct{}{}).
		Set(12, struct{}{}).
		Set(13, struct{}{}).
		Set(14, struct{}{}).
		Set(15, struct{}{}).
		Del(11)
	fmt.Println(w.Repr())
	l, r := w.Split(13)
	fmt.Println(l.Repr(), r.Repr())

	stressTest := true
	if !stressTest {
		return
	}

	// stress test
	type WH struct {
		Weight uint `json:"weight"`
		Height uint `json:"height"`
	}
	statistics := make([]WH, 0)
	n := 100000
	keys := make(map[int]struct{})
	for i := 0; i < n; i++ {
		x := rand.Int()
		w = w.Set(x, struct{}{})
		keys[x] = struct{}{}
		if rand.Float32() < 0.2 {
			// 20% remove one of the keys
			j := rand.Intn(len(keys))
			for k := range keys {
				j--
				if j == 0 {
					delete(keys, k)
					w = w.Del(k)
				}
			}
		}
		// write statistics
		statistics = append(statistics, WH{
			Weight: w.Weight(),
			Height: w.Height(),
		})
	}
	b, err := json.Marshal(statistics)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("tmp/statistics.json", b, 0644)
	if err != nil {
		panic(err)
	}
}

func testPersistentVector() {
	v := vector.NewVector[int]()
	fmt.Println(v.Repr(), v.Weight())
	v = v.Ins(v.Weight(), 0)
	v = v.Ins(v.Weight(), 1)
	v = v.Ins(v.Weight(), 2)
	v = v.Ins(v.Weight(), 3)
	v = v.Ins(v.Weight(), 4)
	v = v.Ins(v.Weight(), 5)
	fmt.Println(v.Repr(), v.Weight())

	v = v.Set(2, 22)
	fmt.Println(v.Repr(), v.Weight())

	v = v.Ins(3, 33)
	fmt.Println(v.Repr(), v.Weight())

	v = v.Del(4)
	fmt.Println(v.Repr(), v.Weight())

	v1, v2 := v.Split(3)
	fmt.Println(v1.Repr(), v1.Weight())
	fmt.Println(v2.Repr(), v2.Weight())

	v = v1.Concat(v2)
	fmt.Println(v.Repr(), v.Weight())

	stressTest := false
	if !stressTest {
		return
	}

	type WH struct {
		Weight uint `json:"weight"`
		Height uint `json:"height"`
	}
	statistics := make([]WH, 0)
	n := 100000

	for i := 0; i < n; i++ {
		x := rand.Int()
		v = v.Ins(v.Weight(), x)
		if rand.Float32() < 0.2 {
			// 20% remove one of the entry
			j := uint(rand.Intn(int(v.Weight())))
			v = v.Del(j)
		}
		// write statistics
		statistics = append(statistics, WH{
			Weight: v.Weight(),
			Height: v.Height(),
		})
	}
	b, err := json.Marshal(statistics)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("tmp/statistics.json", b, 0644)
	if err != nil {
		panic(err)

	}
}

func main() {
	testPersistentVector()
}
