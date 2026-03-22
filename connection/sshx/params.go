package sshx

import (
	"errors"
)

var (
	errDialSSH         = errors.New("failed to dial SSH client")
	errListenLocalPort = errors.New("failed to listen on local port")
	errValidateParams  = errors.New("invalid SSH connection parameters")
	errLoadPrivateKey  = errors.New("failed to load private key")
	errDialSSHServer   = errors.New("failed to dial SSH server")
	errCloseListener   = errors.New("failed to close listener")
)

type Params struct {
	User       string
	Password   string
	RemoteHost string
	KeyPath    string
	LocalPort  int
	RemotePort int
}
