package sshx

import (
	"io"
	"net"
	"sync"
)

func proxyConn(localConn net.Conn, remoteConn net.Conn) {
	defer localConn.Close()  //nolint:errcheck
	defer remoteConn.Close() //nolint:errcheck
	
	var group sync.WaitGroup
	
	group.Add(2) //nolint:mnd
	
	go func() {
		defer group.Done()
		
		_, _ = io.Copy(remoteConn, localConn) //nolint:errcheck
		
		closeWrite(remoteConn)
	}()
	
	go func() {
		defer group.Done()
		
		_, _ = io.Copy(localConn, remoteConn) //nolint:errcheck
		
		closeWrite(localConn)
	}()
	
	group.Wait()
}

func closeWrite(conn net.Conn) {
	type closeWriter interface {
		CloseWrite() error
	}
	
	if cw, ok := conn.(closeWriter); ok {
		err := cw.CloseWrite()
		if err != nil {
			return
		}
		
		return
	}
	
	err := conn.Close()
	if err != nil {
		return
	}
}
