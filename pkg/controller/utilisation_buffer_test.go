package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUtilisationBuffer(t *testing.T) {
	t.Run("size 0 creates buffer of size 1", func(t *testing.T) {
		buf := newUtilisationBuffer(0)
		assert.Equal(t, 1, buf.size)
	})

	t.Run("negative size creates buffer of size 1", func(t *testing.T) {
		buf := newUtilisationBuffer(-5)
		assert.Equal(t, 1, buf.size)
	})

	t.Run("size 1 creates buffer of size 1", func(t *testing.T) {
		buf := newUtilisationBuffer(1)
		assert.Equal(t, 1, buf.size)
	})

	t.Run("size 5 creates buffer of size 5", func(t *testing.T) {
		buf := newUtilisationBuffer(5)
		assert.Equal(t, 5, buf.size)
	})
}

func TestUtilisationBufferAverage(t *testing.T) {
	t.Run("empty buffer returns 0", func(t *testing.T) {
		buf := newUtilisationBuffer(5)
		assert.Equal(t, 0.0, buf.average())
	})

	t.Run("single value returns that value", func(t *testing.T) {
		buf := newUtilisationBuffer(1)
		buf.add(45.5)
		assert.Equal(t, 45.5, buf.average())
	})

	t.Run("buffer with 3 values out of 5 returns average of 3", func(t *testing.T) {
		buf := newUtilisationBuffer(5)
		buf.add(10.0)
		buf.add(20.0)
		buf.add(30.0)
		assert.Equal(t, 20.0, buf.average())
	})

	t.Run("full buffer returns average of all values", func(t *testing.T) {
		buf := newUtilisationBuffer(5)
		buf.add(10.0)
		buf.add(20.0)
		buf.add(30.0)
		buf.add(40.0)
		buf.add(50.0)
		assert.Equal(t, 30.0, buf.average())
	})

	t.Run("circular overwrite returns average of last N values", func(t *testing.T) {
		buf := newUtilisationBuffer(3)
		buf.add(10.0)
		buf.add(20.0)
		buf.add(30.0)
		buf.add(40.0)
		buf.add(50.0)
		assert.Equal(t, 40.0, buf.average())
	})

	t.Run("size 1 buffer with multiple adds returns last value", func(t *testing.T) {
		buf := newUtilisationBuffer(1)
		buf.add(10.0)
		buf.add(20.0)
		buf.add(30.0)
		assert.Equal(t, 30.0, buf.average())
	})
}

func TestUtilisationBufferAdd(t *testing.T) {
	t.Run("add increments count until full", func(t *testing.T) {
		buf := newUtilisationBuffer(3)
		assert.Equal(t, 0, buf.count)

		buf.add(10.0)
		assert.Equal(t, 1, buf.count)

		buf.add(20.0)
		assert.Equal(t, 2, buf.count)

		buf.add(30.0)
		assert.Equal(t, 3, buf.count)

		buf.add(40.0)
		assert.Equal(t, 3, buf.count)

		buf.add(50.0)
		assert.Equal(t, 3, buf.count)
	})

	t.Run("index wraps around correctly", func(t *testing.T) {
		buf := newUtilisationBuffer(3)

		buf.add(10.0)
		assert.Equal(t, 1, buf.index)

		buf.add(20.0)
		assert.Equal(t, 2, buf.index)

		buf.add(30.0)
		assert.Equal(t, 0, buf.index)

		buf.add(40.0)
		assert.Equal(t, 1, buf.index)
	})
}

func TestUtilisationBufferSmoothingScenarios(t *testing.T) {
	t.Run("smoothing prevents action on transient dip", func(t *testing.T) {
		buf := newUtilisationBuffer(5)
		buf.add(55.0)
		buf.add(55.0)
		buf.add(55.0)
		buf.add(55.0)
		buf.add(45.0)
		assert.Equal(t, 53.0, buf.average())
	})

	t.Run("smoothing detects sustained low utilisation", func(t *testing.T) {
		buf := newUtilisationBuffer(5)
		buf.add(45.0)
		buf.add(45.0)
		buf.add(45.0)
		buf.add(45.0)
		buf.add(45.0)
		assert.Equal(t, 45.0, buf.average())
	})

	t.Run("smoothing with partial buffer", func(t *testing.T) {
		buf := newUtilisationBuffer(5)
		buf.add(50.0)
		buf.add(50.0)
		buf.add(95.0)
		assert.Equal(t, 65.0, buf.average())
	})
}
