package models

import "fmt"

type Station int

const (
	NangangStation Station = 1
	TaipeiStation  Station = 2
	Banqiao        Station = 3
	Taoyuan        Station = 4
	Hsinchu        Station = 5
	Miaoli         Station = 6
	Taichung       Station = 7
	Changhua       Station = 8
	Yunlin         Station = 9
	Chiayi         Station = 10
	Tainan         Station = 11
	Zuoying        Station = 12
)

func (s Station) Code() string {
	return fmt.Sprintf("%d", s)
}
