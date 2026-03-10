package controller

// utilisationBuffer is a fixed-size circular buffer that stores recent utilisation
// readings and returns their average. Used for smoothing scaling signals.
type utilisationBuffer struct {
	values []float64
	size   int
	index  int
	count  int
}

// newUtilisationBuffer creates a buffer that holds up to `size` values.
// A size of 0 or 1 means no smoothing (average returns the latest value).
func newUtilisationBuffer(size int) *utilisationBuffer {
	if size < 1 {
		size = 1
	}
	return &utilisationBuffer{
		values: make([]float64, size),
		size:   size,
	}
}

// add inserts a new value into the buffer, overwriting the oldest if full.
func (b *utilisationBuffer) add(value float64) {
	b.values[b.index] = value
	b.index = (b.index + 1) % b.size
	if b.count < b.size {
		b.count++
	}
}

// average returns the mean of all stored values.
// Returns 0 if no values have been added.
func (b *utilisationBuffer) average() float64 {
	if b.count == 0 {
		return 0
	}
	sum := 0.0
	for i := 0; i < b.count; i++ {
		sum += b.values[i]
	}
	return sum / float64(b.count)
}
