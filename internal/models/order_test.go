package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetStartTime(t *testing.T) {
	date := time.Date(2024, time.Month(1), 1, 0, 0, 0, 0, time.Local)

	o := Order{
		DepartureDay: date,
		StartTime:    "08:05",
	}

	startTime := o.GetStartTime()
	assert.Equal(t, "800A", startTime)

	o = Order{
		DepartureDay: date,
		StartTime:    "21:35",
	}

	startTime = o.GetStartTime()
	assert.Equal(t, "930P", startTime)
}
