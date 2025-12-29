//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/deque.go
//

package flagparser

// deque implements a generic deque.
type deque[T any] struct {
	values []T
}

// Empty returns true if the deque is empty.
func (d *deque[T]) Empty() bool {
	return len(d.values) <= 0
}

// Front returns the element at the front.
func (d *deque[T]) Front() (value T, ok bool) {
	if !d.Empty() {
		value, ok = d.values[0], true
	}
	return
}

// PopFront removes the first element if possible.
func (d *deque[T]) PopFront() {
	if !d.Empty() {
		d.values = d.values[1:]
	}
}

// PushBack appends an element to the back.
func (d *deque[T]) PushBack(val T) {
	d.values = append(d.values, val)
}
