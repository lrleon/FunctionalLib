package FunctionalLib

import (
	Seq "github.com/geniussportsgroup/Slist"
	Set "github.com/geniussportsgroup/treaps"
	"github.com/stretchr/testify/assert"
	"testing"
)

func cmpInt(i1, i2 interface{}) bool {
	item1, ok := i1.(int)
	if !ok {
		panic("First parameter is not int")
	}
	item2, ok := i2.(int)
	if !ok {
		panic("Second parameter is not int")
	}
	return item1 < item2
}

const N = 100

func createSet() *Set.Treap {

	tree := Set.New(3, cmpInt)
	for i := 0; i < N; i++ {
		tree.Insert(i)
	}

	return tree
}

func TestForEach(t *testing.T) {

	tree := createSet()

	i := 0
	ForEach(tree, func(k interface{}) {
		assert.Equal(t, k.(int), i)
		i++
	})
}

func TestAll(t *testing.T) {

	tree := createSet()

	assert.True(t, All(tree, func(i interface{}) bool {
		return i.(int) < N
	}))

	assert.False(t, All(tree, func(i interface{}) bool {
		return i.(int) >= N
	}))
}

func TestExist(t *testing.T) {

	tree := createSet()

	assert.True(t, Exist(tree, func(i interface{}) bool {
		return i.(int) <= N-1
	}))

	assert.False(t, Exist(tree, func(i interface{}) bool {
		return i.(int) >= N
	}))
}

func TestMap(t *testing.T) {

	tree := createSet()

	m := Map(tree, func(i interface{}) interface{} {
		return 2 * i.(int)
	})

	assert.True(t, All(Zip(tree, m), func(i interface{}) bool {
		p := i.(Pair)
		return 2*p.Item1.(int) == p.Item2.(int)
	}))
}

func TestZipUnzip(t *testing.T) {

	l1 := createSet()
	l2 := createSet()

	lzip := Zip(l1, l2)
	assert.Equal(t, lzip.Size(), l1.Size())
	assert.Equal(t, lzip.Size(), l2.Size())

	assert.True(t, All(lzip, func(i interface{}) bool {
		p := i.(Pair)
		return p.Item1 == p.Item2
	}))

	r1, r2 := Unzip(lzip)
	assert.Equal(t, r1.Size(), r2.Size())
	assert.Equal(t, r1.Size(), l1.Size())
	for it1, it2 := r1.CreateIterator().(*Seq.Iterator), r2.CreateIterator().(*Seq.Iterator); it1.HasCurr() && it2.HasCurr(); /* nothing */ {
		assert.Equal(t, it1.GetCurr().(int), it2.GetCurr().(int))

		it1.Next()
		it2.Next()
	}
}

func TestFilter(t *testing.T) {

	l := createSet()

	l50 := Filter(l, func(i interface{}) bool {
		return i.(int) < N/2
	})

	assert.True(t, All(l50, func(i interface{}) bool {
		return i.(int) < N/2
	}))

	assert.False(t, Exist(l50, func(i interface{}) bool {
		return i.(int) >= N/2
	}))
}

func TestMapIf(t *testing.T) {

	pred := func(item interface{}) bool {
		return item.(int) >= 20 && item.(int) <= 60
	}

	tree := createSet()

	lmap := MapIf(tree, func(i interface{}) interface{} {
		return 2 * i.(int)
	}, pred)

	assert.True(t, All(Zip(lmap, Filter(tree, pred)), func(i interface{}) bool {
		p := i.(Pair)
		return p.Item1.(int) == 2*p.Item2.(int)
	}))
}

func TestSplit(t *testing.T) {

	l := createSet()

	l1, l2 := Split(l, func(item interface{}) bool {
		return item.(int) < N/2
	})

	assert.True(t, All(l1, func(i interface{}) bool {
		return i.(int) < N/2
	}))

	assert.True(t, All(l2, func(i interface{}) bool {
		return i.(int) >= N/2
	}))
}

func TestFind(t *testing.T) {

	tree := createSet()
	iN2 := Find(tree, func(item interface{}) bool {
		return item.(int) == N/2
	})

	assert.Equal(t, iN2.(int), N/2)

	assert.Nil(t, Find(tree, func(item interface{}) bool {
		return item.(int) >= N
	}))
}

func TestTake(t *testing.T) {

	tree := createSet()
	l10 := Take(tree, 10)
	assert.Equal(t, l10.Size(), 10)
	assert.True(t, All(l10, func(i interface{}) bool {
		return i.(int) <= 9
	}))
}

func TestDrop(t *testing.T) {

	tree := createSet()
	l90 := Drop(tree, 10)
	assert.Equal(t, l90.Size(), N-10)
	assert.True(t, All(l90, func(i interface{}) bool {
		return i.(int) > 9
	}))
}

func TestFoldl(t *testing.T) {

	assert.Equal(t, Foldl(createSet(), 0, func(acu, item interface{}) interface{} {
		return acu.(int) + item.(int)
	}).(int), N*(N-1)/2)
}

func TestNth(t *testing.T) {

	tree := createSet()

	assert.Nil(t, Nth(tree, -1))
	assert.Nil(t, Nth(tree, N))
	for i := 0; i < tree.Size(); i++ {
		assert.Equal(t, Nth(tree, i), i)
	}
}

func TestPosition(t *testing.T) {

	tree := createSet()

	assert.Equal(t, Position(tree, func(item interface{}) bool {
		return item.(int) == -1
	}), -1)

	assert.Equal(t, Position(tree, func(item interface{}) bool {
		return item.(int) == N
	}), -1)

	for i := 0; i < tree.Size(); i++ {
		assert.Equal(t, Position(tree, func(item interface{}) bool {
			return item.(int) == i
		}), i)
	}
}

func TestNewTuple(t *testing.T) {

	tuple4 := NewTuple(1, 2, 3, 4)
	assert.Equal(t, tuple4.Size(), 4)
	assert.Equal(t, tuple4.Nth(0).(int), 1)
}

func TestTZip(t *testing.T) {

	l1 := Seq.New(1, 2, 3, 4, 5)
	l2 := Seq.New("A", "B", "C", "D", "E")
	l3 := Seq.New(-5, -4, -3, -2, -1)

	zl := TZip(l1, l2, l3)

	assert.Equal(t, zl.Size(), l1.Size())

	assert.True(t, All(Zip(zl, l1), func(item interface{}) bool {
		pair := item.(Pair)
		tuple := pair.Item1.(*Tuple)
		i := pair.Item2.(int)
		return tuple.Nth(0) == i
	}))

	assert.True(t, All(Zip(zl, l2), func(item interface{}) bool {
		pair := item.(Pair)
		tuple := pair.Item1.(*Tuple)
		str := pair.Item2.(string)
		return tuple.Nth(1) == str
	}))

	assert.True(t, All(Zip(zl, l3), func(item interface{}) bool {
		pair := item.(Pair)
		tuple := pair.Item1.(*Tuple)
		i := pair.Item2.(int)
		return tuple.Nth(2) == i
	}))
}

func TestTUnzip(t *testing.T) {

	l1 := Seq.New(1, 2, 3, 4, 5)
	l2 := Seq.New("A", "B", "C", "D", "E")
	l3 := Seq.New(-5, -4, -3, -2, -1)

	zl := TZip(l1, l2, l3)

	tuple := TUnzip(zl)

	assert.True(t, All(Zip(tuple.Nth(0).(*Seq.Slist), l1), func(p interface{}) bool {
		pair := p.(Pair)
		i1 := pair.Item1.(int)
		i2 := pair.Item2.(int)
		return i1 == i2
	}))

	assert.True(t, All(Zip(tuple.Nth(1).(*Seq.Slist), l2), func(p interface{}) bool {
		pair := p.(Pair)
		i1 := pair.Item1.(string)
		i2 := pair.Item2.(string)
		return i1 == i2
	}))

	assert.True(t, All(Zip(tuple.Nth(2).(*Seq.Slist), l3), func(p interface{}) bool {
		pair := p.(Pair)
		i1 := pair.Item1.(int)
		i2 := pair.Item2.(int)
		return i1 == i2
	}))
}
