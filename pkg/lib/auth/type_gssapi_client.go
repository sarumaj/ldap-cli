package auth

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-ldap/ldap/v3"
	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/spnego"
)

// GSSAPIClient implements ldap.GSSAPIClient using gokrb5 for Kerberos authentication
type GSSAPIClient struct {
	krbClient   *client.Client
	spnegoCtx   *spnego.SPNEGO
	gssCtx      context.Context
	mu          sync.Mutex
	established bool
}

var _ ldap.GSSAPIClient = (*GSSAPIClient)(nil)

// InitSecContext initiates the GSSAPI context establishment for SASL bind
func (c *GSSAPIClient) InitSecContext(target string, token []byte) (outputToken []byte, needContinue bool, err error) {
	return c.InitSecContextWithOptions(target, token, nil)
}

// InitSecContextWithOptions initiates the GSSAPI context with options (unused for now)
func (c *GSSAPIClient) InitSecContextWithOptions(target string, token []byte, options []int) (outputToken []byte, needContinue bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.krbClient == nil {
		return nil, false, ErrNoKerberosClient
	}

	// Create SPNEGO context if not already
	if c.spnegoCtx == nil {
		c.spnegoCtx = spnego.SPNEGOClient(c.krbClient, target)
	}

	// Acquire credentials if needed
	if err := c.spnegoCtx.AcquireCred(); err != nil {
		return nil, false, err
	}

	// If no token, start negotiation
	if len(token) == 0 {
		ctxToken, err := c.spnegoCtx.InitSecContext()
		if err != nil {
			return nil, false, err
		}
		marshaled, err := ctxToken.Marshal()
		if err != nil {
			return nil, false, err
		}
		c.gssCtx = ctxToken.Context()
		c.established = false
		return marshaled, true, nil
	}

	// Continue negotiation with server token
	spnegoToken := &spnego.SPNEGOToken{}
	if err := spnegoToken.Unmarshal(token); err != nil {
		return nil, false, err
	}
	ok, _, status := c.spnegoCtx.AcceptSecContext(spnegoToken)
	if !ok {
		return nil, false, fmt.Errorf("SPNEGO AcceptSecContext failed: %v", status)
	}
	c.gssCtx = spnegoToken.Context()
	c.established = true
	return nil, false, nil
}

// NegotiateSaslAuth finalizes the SASL negotiation (not used for most LDAP GSSAPI flows)
func (c *GSSAPIClient) NegotiateSaslAuth(token []byte, authzID string) ([]byte, error) {
	// For most LDAP GSSAPI, this is a no-op or returns an empty token
	return nil, nil
}

// DeleteSecContext cleans up the GSSAPI context
func (c *GSSAPIClient) DeleteSecContext() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gssCtx = nil
	c.spnegoCtx = nil
	c.established = false
	return nil
}

// ErrNoKerberosClient is returned if the GSSAPIClient is not initialized with a Kerberos client
var ErrNoKerberosClient = fmt.Errorf("GSSAPIClient: no Kerberos client configured")
