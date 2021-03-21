package FunctionalLib

import (
	Seq "github.com/geniussportsgroup/Slist"
)

type SequentialIterator interface {
	ResetFirst() interface{}
	HasCurr() bool
	GetCurr() interface{}
	Next() interface{}
}

type Sequence interface {
	Traverse(func(interface{}) bool) bool
	Append(item interface{}, items ...interface{}) interface{}
	Size() int
	Swap(other interface{}) interface{}
	IsEmpty() bool
	CreateIterator() interface{}
}

type Pair struct {
	Item1, Item2 interface{}
}

type Tuple struct {
	l *[]interface{}
}

func NewTuple(items ...interface{}) *Tuple {

	s := make([]interface{}, 0, len(items))
	tuple := Tuple{l: &s}
	for _, i := range items {
		*tuple.l = append(*tuple.l, i)
	}
	return &tuple
}

func BuildTuple(n int) *Tuple {
	s := make([]interface{}, n, n)
	return &Tuple{l: &s}
}

func (tuple *Tuple) Set(i int, item interface{}) {
	(*tuple.l)[i] = item
}

func (tuple *Tuple) Traverse(f func(interface{}) bool) bool {
	ptr := tuple.l
	for i := 0; i < len(*tuple.l); i++ {
		if !f((*ptr)[i]) {
			return false
		}
	}
	return true
}

func (tuple *Tuple) Append(item interface{}, items ...interface{}) interface{} {
	*tuple.l = append(*tuple.l, items)
	for _, i := range items {
		*tuple.l = append(*tuple.l, i)
	}
	return tuple
}

func (tuple *Tuple) Size() int {
	return len(*tuple.l)
}

func (tuple *Tuple) Swap(other interface{}) interface{} {
	otherTuple := other.(*Tuple)
	tuple.l, otherTuple.l = otherTuple.l, tuple.l
	return tuple
}

func (tuple *Tuple) IsEmpty() bool {
	return tuple.Size() == 0
}

type TupleIterator struct {
	tuple *Tuple
	pos   int
}

func (tuple *Tuple) CreateIterator() interface{} {
	return &TupleIterator{
		tuple: tuple,
		pos:   0,
	}
}

func NewTupleIterator(tuple Tuple) *TupleIterator {
	return tuple.CreateIterator().(*TupleIterator)
}

func (it *TupleIterator) HasCurr() bool {
	return it.pos < len(*it.tuple.l)
}

func (it *TupleIterator) GetCurr() interface{} {
	return (*it.tuple.l)[it.pos]
}

func (tuple *Tuple) Nth(i int) interface{} {
	return (*tuple.l)[i]
}

func (it *TupleIterator) Next() interface{} {
	it.pos++
	return it
}

func (it *TupleIterator) ResetFirst() interface{} {
	it.pos = 0
	return it
}

// Execute operation receiving every item of the sequence. Return seq
func ForEach(seq Sequence, operation func(interface{})) interface{} {

	seq.Traverse(func(i interface{}) bool {
		operation(i)
		return true
	})

	return seq
}

// Return true if all the elements of the sequence meets predicate
func All(seq Sequence, predicate func(interface{}) bool) bool {

	return seq.Traverse(predicate)
}

func Exist(seq Sequence, predicate func(interface{}) bool) bool {

	return !All(seq, func(i interface{}) bool {
		return !predicate(i)
	})
}

// Return a new seq with the items of sequence transformed with transformation operation
func Map(seq Sequence, transformation func(interface{}) interface{}) *Seq.Slist {

	ret := Seq.New()
	ForEach(seq, func(item interface{}) {
		ret.Append(transformation(item))
	})
	return ret
}

// Return a new seq with the items of sequence transformed with transformation operation for the items
// satisfying the predicate
func MapIf(seq Sequence,
	transformation func(interface{}) interface{},
	predicate func(interface{}) bool) *Seq.Slist {

	ret := Seq.New()
	ForEach(seq, func(item interface{}) {
		if predicate(item) {
			ret.Append(transformation(item))
		}
	})
	return ret
}

// Return a list containing the items satisfying predicate
func Filter(seq Sequence, predicate func(interface{}) bool) *Seq.Slist {

	ret := Seq.New()
	ForEach(seq, func(item interface{}) {
		if predicate(item) {
			ret.Append(item)
		}
	})
	return ret
}

// Zip two lists into one list of pair. The result is truncated to the shortest list
func Zip(s1, s2 Sequence) *Seq.Slist {

	ret := Seq.New()

	it1, it2 := s1.CreateIterator().(SequentialIterator), s2.CreateIterator().(SequentialIterator)
	for it1.HasCurr() && it2.HasCurr() {

		ret.Append(Pair{
			Item1: it1.GetCurr(),
			Item2: it2.GetCurr(),
		})

		it1.Next()
		it2.Next()
	}

	return ret
}

// Unzip a list of pairs into two separated lists
func Unzip(seq *Seq.Slist) (*Seq.Slist, *Seq.Slist) {

	l1 := Seq.New()
	l2 := Seq.New()

	for it := seq.CreateIterator().(*Seq.Iterator); it.HasCurr(); it.Next() {

		curr := it.GetCurr().(Pair)
		l1.Append(curr.Item1)
		l2.Append(curr.Item2)
	}

	return l1, l2
}

// Split seq into two lists. First contains items satisfying predicate and the second the complement
func Split(seq Sequence, predicate func(item interface{}) bool) (*Seq.Slist, *Seq.Slist) {

	l1 := Seq.New()
	l2 := Seq.New()

	ForEach(seq, func(i interface{}) {
		if predicate(i) {
			l1.Append(i)
		} else {
			l2.Append(i)
		}
	})

	return l1, l2
}

// Return the first item in seq satisfying predicate. If not item is found, the it returns nil
func Find(seq Sequence, predicate func(item interface{}) bool) interface{} {

	for it := seq.CreateIterator().(SequentialIterator); it.HasCurr(); it.Next() {
		item := it.GetCurr()
		if predicate(item) {
			return item
		}
	}
	return nil
}

// Return a sequence containing the first n items from the sequence
func Take(seq Sequence, n int) *Seq.Slist {

	ret := Seq.New()
	for it := seq.CreateIterator().(SequentialIterator); it.HasCurr() && n > 0; it.Next() {
		ret.Append(it.GetCurr())
		n--
	}

	return ret
}

// Return a sequence containing the items after the first n from the sequence
func Drop(seq Sequence, n int) *Seq.Slist {

	ret := Seq.New()
	for i, it := 0, seq.CreateIterator().(SequentialIterator); it.HasCurr(); it.Next() {
		if i < n {
			i++
			continue
		}
		ret.Append(it.GetCurr())
	}

	return ret
}

// Return f(in, ..., f(i2, f(i1, initVal) ... ))
func Foldl(seq Sequence, initVal interface{},
	f func(acu, item interface{}) interface{}) interface{} {

	retVal := initVal
	ForEach(seq, func(i interface{}) {
		retVal = f(retVal, i)
	})
	return retVal
}

// Return the n-th item in the sequence. Return nil if n is negative o greater than seq.Size()
func Nth(seq Sequence, n int) interface{} {

	if n < 0 || n >= seq.Size() {
		return nil
	}

	for it := seq.CreateIterator().(SequentialIterator); it.HasCurr(); it.Next() {
		if n == 0 {
			return it.GetCurr()
		}
		n--
	}

	return nil // it should not be reached!
}

// Return the position in the sequence of the first element satisfying predicate. If no element satisfies
// predicate then it returns -1
func Position(seq Sequence, predicate func(item interface{}) bool) int {

	for pos, it := 0, seq.CreateIterator().(SequentialIterator); it.HasCurr(); it.Next() {
		if predicate(it.GetCurr()) {
			return pos
		}
		pos++
	}

	return -1
}

// Zip all the lists into a list of tuples
func TZip(list *Seq.Slist, lists ...*Seq.Slist) *Seq.Slist {

	sz := len(lists) + 1
	ret := Seq.New()

	for it := Seq.NewIterator(list); it.HasCurr(); it.Next() {
		tuple := BuildTuple(sz)
		tuple.Set(0, it.GetCurr())
		ret.Append(tuple)
	}

	for i, l := range lists {

		zipL := Zip(ret, l)
		for it := Seq.NewIterator(zipL); it.HasCurr(); it.Next() { // traverse i-th list
			pair := it.GetCurr().(Pair)
			tuple := pair.Item1.(*Tuple)
			item := pair.Item2
			tuple.Set(i+1, item)
		}
	}

	return ret
}

// Unzip a list of tuples into a tuple of lists
func TUnzip(tupleList *Seq.Slist) *Tuple {

	tupleSize := tupleList.First().(*Tuple).Size()
	result := BuildTuple(tupleSize)
	for i := 0; i < tupleSize; i++ {
		result.Set(i, Seq.New())
	}

	for it := Seq.NewIterator(tupleList); it.HasCurr(); it.Next() {
		tuple := it.GetCurr().(*Tuple)
		for i := 0; i < tuple.Size(); i++ {
			result.Nth(i).(*Seq.Slist).Append(tuple.Nth(i))
		}
	}

	return result
}
