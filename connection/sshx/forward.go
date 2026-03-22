package sshx

import (
	"fmt"
	"log"
	"net"
	
	"golang.org/x/crypto/ssh"
)

func forward(params *Params, listener net.Listener, sshClient *ssh.Client) {
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("%w: %v", errCloseListener, err)
		}
	}(listener)
	
	for {
		localConnection, err := listener.Accept()
		if err != nil {
			continue
		}
		
		remoteConnection, err := sshClient.Dial("tcp", fmt.Sprintf("%s:%s", params.RemoteHost, params.RemotePort))
		if err != nil {
			err = localConnection.Close()
			
			if err != nil {
				continue
			}
			
			continue
		}
		
		go proxyConn(localConnection, remoteConnection)
	}
}
