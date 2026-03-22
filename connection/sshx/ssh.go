package sshx

import (
	"fmt"
)

func StartSSHTunnel(params *Params) (chan bool, error) {
	sshClient, err := dialSSHClient(params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDialSSH, err)
	}
	
	listener, err := listenLocal(params, sshClient)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errListenLocalPort, err)
	}
	
	done := make(chan bool)
	
	go forward(params, listener, sshClient)
	
	return done, nil
}
