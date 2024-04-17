package crawler

import (
	"testing"
	"time"
	"tix-worker/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestFilterValidTrain(t *testing.T) {
	order := models.Order{
		DepartureDay: time.Date(2024, time.Month(1), 1, 16, 0, 0, 0, time.UTC),
		StartTime:    "18:30",
		EndTime:      "19:00",
	}

	trainDatas := []TrainData{
		{
			TrainCode:     "1",
			Value:         "radio31",
			Date:          "01/01",
			DepartureTime: "18:20",
			ArrivalTime:   "20:20",
		},
		{
			TrainCode:     "2",
			Value:         "radio31",
			Date:          "01/01",
			DepartureTime: "18:30",
			ArrivalTime:   "20:30",
		},
		{
			TrainCode:     "3",
			Value:         "radio31",
			Date:          "01/01",
			DepartureTime: "18:40",
			ArrivalTime:   "20:40",
		},
		{
			TrainCode:     "4",
			Value:         "radio31",
			Date:          "01/01",
			DepartureTime: "19:00",
			ArrivalTime:   "21:00",
		},
		{
			TrainCode:     "5",
			Value:         "radio31",
			Date:          "01/01",
			DepartureTime: "19:05",
			ArrivalTime:   "21:05",
		},
	}

	c := Crawler{}
	validTrans := c.filterValidTrain(order, trainDatas)
	assert.Equal(t, 3, len(validTrans))
}
