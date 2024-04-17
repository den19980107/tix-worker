package crawler

import (
	"fmt"
	"strconv"
	"strings"
)

type ConfirmTrain struct {
	BookingS2FormHf0                  string `json:"BookingS2Form:hf:0"`
	TrainQueryDataViewPanelTrainGroup string `json:"TrainQueryDataViewPanel:TrainGroup"`
}

type ConfirmTicket struct {
	BookingS3FormSpHf0                                                          string `json:"bookingS3FormSP:hf:0"`
	DiffOver                                                                    string `json:"diffOver"`
	IsSPromotion                                                                string `json:"isSPromotion"`
	PassengerCount                                                              string `json:"passengerCount"`
	IsGoBackM                                                                   string `json:"isGoBackM"`
	BackHome                                                                    string `json:"backHome"`
	TgoError                                                                    string `json:"tgoError"`
	IdInputRadio                                                                string `json:"idInputRadio"`
	DummyId                                                                     string `json:"dummyId"`
	DummyPhone                                                                  string `json:"dummyPhone"`
	Email                                                                       string `json:"email"`
	Agree                                                                       string `json:"agree"`
	TicketMemberSystemInputPanelTakerMemberSystemDataViewMemberSystemRadioGroup string `json:"TicketMemberSystemInputPanel:TakerMemberSystemDataView:memberSystemRadioGroup"`
}

type SubmitForm struct {
	SelectStartStation            string `json:"selectStartStation"`
	SelectDestinationStation      string `json:"selectDestinationStation"`
	BookingMethod                 string `json:"bookingMethod"`
	TripConTypesoftrip            string `json:"tripCon:typesoftrip"`
	ToTimeInputField              string `json:"toTimeInputField"`
	ToTimeTable                   string `json:"toTimeTable"`
	HomeCaptchaSecurityCode       string `json:"homeCaptcha:securityCode"`
	SeatConSeatRadioGroup         string `json:"seatCon:seatRadioGroup"`
	BookingS1FormHf0              string `json:"BookingS1Form:hf:0"`
	TrainConTrainRadioGroup       string `json:"trainCon:trainRadioGroup"`
	BackTimeInputField            string `json:"backTimeInputField"`
	BackTimeTable                 string `json:"backTimeTable"`
	ToTrainIDInputField           string `json:"toTrainIDInputField"`
	BackTrainIDInputField         string `json:"backTrainIDInputField"`
	TicketPanelRows0TicketAmount  string `json:"ticketPanel:rows:0:ticketAmount"`
	TicketPanelRows1TicketAmount  string `json:"ticketPanel:rows:1:ticketAmount"`
	TicketPannelRows2TicketAmount string `json:"ticketPannel:rows:2:ticketAmount"`
	TicketPanelRows3TicketAmount  string `json:"ticketPanel:rows:3:ticketAmount"`
	TicketPanelRows4TicketAmount  string `json:"ticketPanel:rows:4:ticketAmount"`
}

type TrainData struct {
	TrainCode     string
	Value         string
	Date          string
	DepartureTime string
	ArrivalTime   string
}

func (t TrainData) getHour() (int, error) {
	parts := strings.Split(t.DepartureTime, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("departure time \"%s\" format not correct", t.DepartureTime)
	}

	return strconv.Atoi(parts[0])
}

func (t TrainData) getMin() (int, error) {
	parts := strings.Split(t.DepartureTime, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("departure time \"%s\" format not correct", t.DepartureTime)
	}

	return strconv.Atoi(parts[1])
}
