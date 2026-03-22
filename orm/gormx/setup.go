package gormx

import (
	"fmt"
	
	"github.com/elmansyah/golang-lib/connection/sshx"
)

func (params *Params) Setup() (*Params, error) {
	setupDB, err := connect(params)
	if err != nil {
		return nil, err
	}
	
	return setupDB, nil
}

func connect(params *Params) (*Params, error) {
	if params.AppMode == "dev" {
		startSSH, err := sshx.StartSSHTunnel(&sshx.Params{
			User:       params.SSHParams.User,
			Password:   params.SSHParams.Password,
			RemoteHost: params.SSHParams.RemoteHost,
			KeyPath:    params.SSHParams.KeyPath,
			LocalPort:  params.SSHParams.LocalPort,
			RemotePort: params.SSHParams.RemotePort,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errSSHTunnel, err)
		}
		
		defer close(startSSH)
	}
	
	openDB, err := open(params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errFailedToOpenDB, err)
	}
	
	if openDB == nil {
		return nil, fmt.Errorf("%w: %w", errOpenDBNil, err)
	}
	
	params.DB = openDB
	params.Closed = func() error {
		return params.Close()
	}
	
	return params, nil
}
