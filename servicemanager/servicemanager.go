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

// +build windows

package servicemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ezbastion/ezb_vault/server"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

// EventLog management
var elog debug.Log
var exPath string

type myservice struct{}

func init() {
	ex, _ := os.Executable()
	exPath = filepath.Dir(ex)
}

func exePath() (string, error) {
	log.Debugln("DBTP:Entering func exePath")
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

// StartService starts the windows service targeted by name
func StartService(name string) error {
	log.Debugln(fmt.Sprintf("DBTP:Entering func startService name : %s", name))

	m, err := mgr.Connect()
	if err != nil {
		log.Errorln(fmt.Sprintf("could not connect the service control manager error : %s", err.Error()))
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not access service (OpenService): %s", name))
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		log.Errorln(fmt.Sprintf("could not start service %s, error : %s", name, err.Error()))
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

// ControlService controls the service targetede by name
func ControlService(name string, c svc.Cmd, to svc.State) error {
	log.Debugln(fmt.Sprintf("DBTP:Entering func controlService name : %s", name))

	m, err := mgr.Connect()
	if err != nil {
		log.Errorln(fmt.Sprintf("could not connect the service %s, error : %s", name, err.Error()))
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not access service (OpenService): %s", name))
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not send control=%d: %s", c, err.Error()))
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			log.Errorln(fmt.Sprintf("timeout waiting for service to go to state=%d", to))
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			log.Errorln(fmt.Sprintf("could not send control=%d: %s", c, err.Error()))
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}

func (m *myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {

	log.Debugln("DBTP:Entering func Execute")

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	serverchan := make(chan bool)
	go server.MainGin(&serverchan)
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				close(serverchan)
				break loop
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

// RunService runs the service targeted by name
func RunService(name string, isDebug bool) {
	log.Debugln(fmt.Sprintf("DBTP:Entering func runService name :%s, isDebug : %t", name, isDebug))

	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	log.Debugln(fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}

	err = run(name, &myservice{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		log.Debugln(fmt.Sprintf("%s service failed", name))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
	log.Debugln(fmt.Sprintf("%s service stoped", name))
}

// InstallService installs the service targeted by name
func InstallService(name, desc string) error {
	log.Debugln(fmt.Sprintf("DBTP:Entering func installService with name : %s, desc : %s", name, desc))
	exepath, err := exePath()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}
	s, err = m.CreateService(name, exepath, mgr.Config{DisplayName: desc}, "is", "auto-started")
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

// RemoveService remove the service trageted by name
func RemoveService(name string) error {
	log.Debugln(fmt.Sprintf("DBTP:Entering func removeService with name : %s", name))

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}
