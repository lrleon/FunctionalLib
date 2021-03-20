package FunctionalLib

import Seq "github.com/geniussportsgroup/Slist"

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

type Pair struct {
	item1, item2 interface{}
}

// Zip two lists into one list of pair. The result is truncated to the shortest list
func Zip(s1, s2 Sequence) *Seq.Slist {

	ret := Seq.New()

	it1, it2 := s1.CreateIterator().(SequentialIterator), s2.CreateIterator().(SequentialIterator)
	for it1.HasCurr() && it2.HasCurr() {

		ret.Append(Pair{
			item1: it1.GetCurr(),
			item2: it2.GetCurr(),
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
		l1.Append(curr.item1)
		l2.Append(curr.item2)
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
