package gormx

import (
	"fmt"
	
	"github.com/elmansyah/golang-lib/connection/sshx"
)

// Setup initializes the database connection layer.
//
// This is the entry point for establishing a DB connection.
// It delegates the actual connection logic to `connect()`,
// which handles validation, SSH tunneling (if needed),
// and database initialization.
//
// Returns:
// - *Params with active DB connection
// - error if setup fails at any stage.
func (params *Params) Setup() (*Params, error) {
	setupDB, err := connect(params)
	if err != nil {
		return nil, err
	}
	
	return setupDB, nil
}

// connect handles database connection with optional SSH tunneling.
//
// Flow:
// 1. Validate DB_LOCATION and DB_TUNNEL values
// 2. Determine whether an SSH tunnel is needed
// 3. Start an SSH tunnel (if enabled)
// 4. Resolve actual DB host/port (direct or via tunnel)
// 5. Open a database connection
// 6. Register cleanup handler (DB + SSH)
//
// Supported combinations:
// - local  + none → direct local connection
// - remote + none → direct remote connection
// - local  + ssh  → SSH tunnel to local (rare, mostly for testing)
// - remote + ssh  → SSH tunnel to remote DB (common production case).
func connect(params *Params) (*Params, error) { //nolint:revive
	validLocation := map[string]bool{
		"local":  true,
		"remote": true,
	}
	
	validTunnel := map[string]bool{
		"none": true,
		"ssh":  true,
	}
	
	if !validLocation[params.DBLocation] {
		return nil, fmt.Errorf(
			"%w: %s (allowed: local, remote)",
			errInvalidDBLocation,
			params.DBLocation,
		)
	}
	
	if !validTunnel[params.DBTunnel] {
		return nil, fmt.Errorf(
			"%w: %s (allowed: none, ssh)",
			errInvalidDBTunnel,
			params.DBTunnel,
		)
	}
	
	// Determine whether connection should go through an SSH tunnel.
	// If true, DB connection will be redirected via localhost:localPort.
	useTunnel := params.DBTunnel == "ssh"
	
	// Default DB target (direct connection).
	// These values may be overridden if SSH tunneling is enabled.
	dbHost := params.Host
	dbPort := params.Port
	
	// sshClose is a cleanup function to close the SSH tunnel when done.
	var sshClose func()
	
	// Start SSH tunnel if required.
	if useTunnel {
		// Start the SSH tunnel to the target server.
		//
		// This creates a secure connection from:
		//   localhost:<LocalPort> → RemoteHost:<RemotePort>
		//
		// After the tunnel is established, the application will connect
		// to the database via the local forwarded port instead of directly
		// accessing the remote host.
		//
		// Example:
		//   localhost:5000 → db-server:5432
		//
		// This is commonly used when the database is not publicly accessible
		// and must be reached through an SSH jump host (bastion server).
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
		
		// Register SSH tunnel close function.
		// This will be executed when the application shuts down.
		sshClose = func() {
			close(sshChan)
		}
		
		// When using an SSH tunnel, the DB is accessed via:
		// localhost:<localPort> instead of the original remote host.
		dbHost = "localhost"
		dbPort = params.SSHParams.LocalPort
	}
	
	// Apply resolved host and port to params.
	// This ensures the DB driver uses the correct endpoint.
	params.Host = dbHost
	params.Port = dbPort
	
	// Open a database connection using the resolved configuration.
	openDB, err := open(params)
	if err != nil {
		if sshClose != nil {
			sshClose()
		}
		
		return nil, fmt.Errorf("%w: %w", errFailedToOpenDB, err)
	}
	
	// Safety check: DB should never be nil if no error occurred.
	if openDB == nil {
		if sshClose != nil {
			sshClose()
		}
		
		return nil, fmt.Errorf("%w", errOpenDBNil)
	}
	
	// Assign active DB connection to params.
	params.DB = openDB
	
	// Register a unified cleanup function.
	// This ensures:
	// - an SSH tunnel is closed (if used)
	// - DB connection is properly closed
	params.Closed = func() error {
		if sshClose != nil {
			sshClose()
		}
		
		return params.Close()
	}
	
	return params, nil
}
