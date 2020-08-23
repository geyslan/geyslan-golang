package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type LogTime struct {
	time.Time
}

func (t LogTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format("\"02/Jan/2006:15:04:05 -0700\"")), nil
}

type clf struct {
	// Common Log Format - https://httpd.apache.org/docs/2.4/logs.html
	ClientHost string  `json:"client_host"`  // Remote Host (1st field)
	RFC1413    string  `json:"-"`            // set when IdentityCheck is on (2nd field)
	RemoteUser string  `json:"remote_user"`  // userid (3rd field)
	DateTime   LogTime `json:"request_time"` // Received request date-time (4th field)
	Method     string  `json:"method"`       // Method (first part of 5th field)
	Resource   string  `json:"request"`      // Resource (second part of 5th field)
	Protocol   string  `json:"protocol"`     // Protocol (third part of 5th field)
	Status     uint16  `json:"status"`       // Status code (6th field)
	Size       uint64  `json:"size"`         // Size of the object returned (7th field)

	// Combined Long Format
	Referer   string `json:"referer"`    // Referer (sic) HTTP request header (8th field)
	UserAgent string `json:"user_agent"` // User-Agent HTTP request header (9th field)
}

func parse(line string) (log clf, err error) {
	const dateLayout = "02/Jan/2006:15:04:05 -0700"
	var tmpDate, tmpRequest, tmpTime, tmpStatus, tmpSize string

	n, err := fmt.Sscanf(line, "%s %s %s [%s %5s] %q %s %s %q %q",
		&log.ClientHost,
		&log.RFC1413,
		&log.RemoteUser,
		&tmpDate,
		&tmpTime,
		&tmpRequest,
		&tmpStatus,
		&tmpSize,
		&log.Referer,
		&log.UserAgent)
	if err != nil {
		return clf{}, errors.New("Invalid CLF entry (field " + strconv.Itoa(n+1) + ")")
	}

	log.DateTime.Time, err = time.Parse(dateLayout, tmpDate+" "+tmpTime)
	if err != nil {
		return clf{}, errors.New("Invalid CLF entry (DateTime)")
	}

	n, err = fmt.Sscanf(tmpRequest, "%s %s %s", &log.Method, &log.Resource, &log.Protocol)
	if err != nil {
		return clf{}, errors.New("Invalid CLF entry (Request field " + strconv.Itoa(n+1) + ")")
	}

	i64, _ := strconv.ParseUint(tmpStatus, 10, 16)
	log.Status = uint16(i64)
	log.Size, _ = strconv.ParseUint(tmpSize, 10, 64)

	return
}
