package sshx

import (
	"io"
	"net"
)

func proxyConn(localConn net.Conn, remoteConn net.Conn) {
	defer localConn.Close()
	defer remoteConn.Close()
	
	done := make(chan struct{}, 2)
	
	go func() {
		_, _ = io.Copy(remoteConn, localConn)
		closeWrite(remoteConn)
		done <- struct{}{}
	}()
	
	go func() {
		_, _ = io.Copy(localConn, remoteConn)
		closeWrite(localConn)
		done <- struct{}{}
	}()
	
	<-done
}

func closeWrite(conn net.Conn) {
	type closeWriter interface {
		CloseWrite() error
	}
	
	if cw, ok := conn.(closeWriter); ok {
		_ = cw.CloseWrite()
		return
	}
	
	_ = conn.Close()
}
