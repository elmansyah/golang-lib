package sshx

import (
	"errors"
	"fmt"
	"os"
	"strings"
	
	"golang.org/x/crypto/ssh"
)

func dialSSHClient(params *Params) (*ssh.Client, error) {
	if err := validateParams(params); err != nil {
		return nil, fmt.Errorf("%w: %w", errValidateParams, err)
	}
	
	signer, err := loadPrivateKey(params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errLoadPrivateKey, err)
	}
	
	config := newSSHClientConfig(params, signer)
	
	sshConnection, err := dialSSHServer(params, config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDialSSHServer, err)
	}
	
	return sshConnection, nil
}

func validateParams(params *Params) error {
	var err []string
	
	if params.User == "" {
		err = append(err, fmt.Sprintf("SSH_USER environment variable is not set: %s", params.User))
	}
	
	if params.KeyPath == "" {
		err = append(err, fmt.Sprintf("SSH_KEY environment variable is not set: %s", params.KeyPath))
	}
	
	if params.RemoteHost == "" {
		err = append(err, fmt.Sprintf("DB_HOST_MASTER environment variable is not set: %s", params.RemoteHost))
	}
	
	if params.RemotePort == 0 {
		err = append(err, fmt.Sprintf("SSH_PORT environment variable is not set: %s", params.RemotePort))
	}
	
	if len(err) > 0 {
		return errors.New(strings.Join(err, ", ")) //nolint:err113
	}
	
	return nil
}

func loadPrivateKey(params *Params) (ssh.Signer, error) {
	key, err := os.ReadFile(params.KeyPath)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	
	return signer, nil
}

func newSSHClientConfig(params *Params, signer ssh.Signer) (sshConfig *ssh.ClientConfig) {
	return &ssh.ClientConfig{
		User: params.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		BannerCallback:    nil,
		ClientVersion:     "",
		HostKeyAlgorithms: nil,
		Timeout:           0,
	}
}

func dialSSHServer(params *Params, sshConfig *ssh.ClientConfig) (sshConnection *ssh.Client, err error) {
	return ssh.Dial("tcp", fmt.Sprintf("%s:%s", params.RemoteHost, params.RemotePort), sshConfig) //nolint:wrapcheck
}
