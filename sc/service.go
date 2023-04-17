package sc

import (
	"errors"

	"github.com/kardianos/service"
)

var (
	s service.Service
)

type Action struct {
	start func()
	stop  func()
}

type Preferences = service.Config

func (a Action) Start(service.Service) error {
	if a.start == nil {
		return nil
	}

	go a.start()
	return nil
}

func (a Action) Stop(service.Service) error {
	if a.stop == nil {
		return nil
	}

	go a.stop()
	return nil
}

func Init(start, stop func(), conf *Preferences) error {
	sv, err := service.New(Action{start: start, stop: stop}, conf)
	if err != nil {
		return err
	}

	s = sv
	return nil
}

func Get() service.Service {
	return s
}

func Install() bool {
	return s.Install() != nil
}

func Start() error {
	return s.Run()
}

func Reinstall() error {
	if !NotInstalled() {
		if err := s.Uninstall(); err != nil {
			return err
		}
	}

	return s.Install()
}

func UnInstall() bool {
	err := s.Uninstall()
	return err != nil
}

func NotInstalled() bool {
	_, err := s.Status()
	return errors.Is(err, service.ErrNotInstalled)
}

func Status() service.Status {
	status, err := s.Status()
	if err != nil {
		if errors.Is(err, service.ErrNotInstalled) {
			return service.StatusUnknown
		}
	}

	return status
}
