package sshx

import (
	"context"
	"fmt"
	"net"
	
	"golang.org/x/crypto/ssh"
)

func listenLocal(params *Params, sshClient *ssh.Client) (net.Listener, error) {
	listenConfig := net.ListenConfig{}
	
	listener, err := listenConfig.Listen(context.Background(), "tcp", fmt.Sprintf("localhost:%d", params.LocalPort))
	if err != nil {
		if closeErr := sshClient.Close(); closeErr != nil {
			return nil, fmt.Errorf("listen error: %w; also failed to close SSH: %w", err, closeErr)
		}
		
		return nil, fmt.Errorf("failed to listen on local port: %w", err)
	}
	
	return listener, nil
}
