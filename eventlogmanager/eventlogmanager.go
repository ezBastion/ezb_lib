// This file is part of ezBastion.

//     ezBastion is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     ezBastion is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.

//     You should have received a copy of the GNU Affero General Public License
//     along with ezBastion.  If not, see <https://www.gnu.org/licenses/>.

// Package eventlogmanager add helper for eventlogs on windows
package eventlogmanager

import (
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

// EventLog management
var elog debug.Log
var status int
var eid int

func init() {
	status = 0
}

// Open open a eventlog specified by name, returning nil or an error
func Open(name string) error {
	var err error

	elog, err = eventlog.Open(name)
	if err != nil {
		status = 255
		return err
	}

	return nil
}

// Info logs an info event into the windows eventlog system
func Info(logline string) error {
	if status == 0 {
		elog.Info(1, logline)
	}

	return nil
}

// Error logs an error event into the windows eventlog system
func Error(logline string) error {
	if status == 0 {
		elog.Error(1, logline)
	}

	return nil
}

// Warning logs an warning event into the windows eventlog system
func Warning(logline string) error {
	if status == 0 {
		elog.Warning(1, logline)
	}

	return nil
}
