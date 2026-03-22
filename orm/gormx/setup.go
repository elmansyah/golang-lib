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

// connectDB handles database connection with optional SSH tunneling.
func connect(params *Params) (*Params, error) {
	// key for lookup
	key := fmt.Sprintf("%s:%s", params.DBLocation, params.DBTunnel)
	
	// whether SSH tunnel is needed
	tunnelNeeded := map[string]bool{
		"remote:ssh": true,
		"local:ssh":  true,
	}
	
	var sshClose func() error
	
	// default host/port (before possible SSH override)
	dbHost := params.Host
	dbPort := params.Port
	
	// handle SSH tunnel if needed
	if tunnelNeeded[key] {
		sshChan, err := sshx.StartSSHTunnel(&sshx.Params{
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
		
		sshClose = func() error {
			close(sshChan) // close SSH tunnel
			return nil
		}
		
		// redirect DB connection to localhost:localPort
		dbHost = "localhost"
		dbPort = params.SSHParams.LocalPort
	} else if params.DBTunnel != "none" {
		// unknown combination
		return nil, fmt.Errorf("%w: %s", errUnknownDBLocation, key)
	}
	
	// assign resolved host/port
	params.Host = dbHost
	params.Port = dbPort
	
	// open DB
	openDB, err := open(params)
	if err != nil {
		if sshClose != nil {
			_ = sshClose()
		}
		
		return nil, fmt.Errorf("%w: %w", errFailedToOpenDB, err)
	}
	
	if openDB == nil {
		if sshClose != nil {
			_ = sshClose()
		}
		
		return nil, fmt.Errorf("%w: %w", errOpenDBNil, err)
	}
	
	params.DB = openDB
	
	// unified close function
	params.Closed = func() error {
		if sshClose != nil {
			_ = sshClose() // close SSH tunnel if active
		}
		
		return params.Close() // close DB
	}
	
	return params, nil
}
