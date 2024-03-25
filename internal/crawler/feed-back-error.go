package crawler

import (
	"strings"
)

type FeedBackError struct {
	Type   FeedBackErrorType
	Detail string
}

func (e FeedBackError) Error() string {
	return e.Detail
}

type FeedBackErrorType string

const (
	Unknow                    FeedBackErrorType = "未知錯誤"
	CaptchaNotCorrect         FeedBackErrorType = "檢測碼輸入錯誤"
	DepartureDateNotAvaliable FeedBackErrorType = "您所選擇的日期超過目前開放預訂之日期"
)

func GetFeedBackError(text string) FeedBackError {
	if strings.Contains(text, string(CaptchaNotCorrect)) {
		return FeedBackError{
			Type:   CaptchaNotCorrect,
			Detail: text,
		}
	}

	if strings.Contains(text, string(DepartureDateNotAvaliable)) {
		return FeedBackError{
			Type:   DepartureDateNotAvaliable,
			Detail: text,
		}
	}

	return FeedBackError{
		Type:   Unknow,
		Detail: text,
	}
}
