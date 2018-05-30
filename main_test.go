package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFlagCondition(t *testing.T) {
	cases := []struct {
		from, to uint
		hasErr   bool
	}{
		{from: 2, to: 1, hasErr: true},
		{from: 1, to: 1, hasErr: true},
		{to: 5, hasErr: true},
		{from: 1, to: 2},
		{from: 5},
	}

	for _, c := range cases {
		*from = 0
		*to = 0

		*from = c.from
		*to = c.to

		err := checkFlagCondition()
		if c.hasErr {
			assert.Error(t, err, "from: %d, to: %d", c.from, c.to)
		} else {
			assert.NoError(t, err, "from: %d, to: %d", c.from, c.to)
		}
	}
}
