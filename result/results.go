package result

import (
	"fmt"
)

type Result interface {
	MarshalJSON() ([]byte, error)
	TcpFailCount() int
	SuccesCount() int
	GetHost() string
}

type PayloadResult struct {
	ComResult
	PayloadInfo []SinglePayload
	prepos      bool
}

func (s PayloadResult) MarshalJSON() ([]byte, error) {
	//TODO: implemet later
	st := fmt.Sprintf(
		`{
		"Host": "%s",
		"TotalTcpFail": %d,
		"Success": %d,
		"Maybe": %d,
		"MaxSpeed": %f,
		"Err": "%s"
		`, 
		s.ComResult.Host, s.ComResult.TotalTcpFail, s.ComResult.Success, s.ComResult.Maybe, s.ComResult.MaxSpeed, s.Err,
	)

	if len(s.PayloadInfo) > 0 {
		st = st + ","+ `"PayloadInfo": [`
		for i, item := range s.PayloadInfo {
			st = st + fmt.Sprintf(
				`{
					"Success": %t,
					"Maybe": %t,
					"Error": "%s",
					"PayloadName": "%s",
					"TpcFailed": %t,
					"Tls": {
						"Failed": %t,
						"Error": "%s",
						"Servername": "%s"
					},
					"MaxSpeed": %f
				}`, item.Success, item.Maybe, item.Error, item.PayloadName, item.TcpFailed, item.Tls.Failed, item.Tls.Error, item.Tls.Servername, item.MaxSpeed,
			)
			if i < len(s.PayloadInfo)-1 {
				st = st + ","
			}

		}
		st = st + "]"
	}
	st = st + "}"

	return []byte(st), nil

}
func (s PayloadResult) TcpFailCount() int {
	return s.TotalTcpFail
}
func (s PayloadResult) SuccesCount() int {
	return s.Success
}
func (s PayloadResult) GetHost() string {
	return s.ComResult.Host
}

func (s *PayloadResult) PreProcess() {
	if s.prepos {
		return
	}
	for _, v := range s.PayloadInfo {
		if v.MaxSpeed > s.MaxSpeed {
			s.MaxSpeed = v.MaxSpeed
		}
		if v.TcpFailed {
			s.TotalTcpFail++
		}
		if v.Success {
			s.Success++
		}
		if v.Maybe {
			s.Maybe++
		}

	}
	s.prepos = true
}

type SinglePayload struct {
	Success     bool
	Maybe       bool
	Error       string
	PayloadName string
	TcpFailed   bool
	Tls         TlsInfo
	MaxSpeed    float64
}

type TlsInfo struct {
	Failed     bool
	Error      string
	Servername string
}

type ComResult struct {
	Host         string
	TotalTcpFail int
	Success      int
	Maybe        int
	MaxSpeed     float64
	Err          string
}