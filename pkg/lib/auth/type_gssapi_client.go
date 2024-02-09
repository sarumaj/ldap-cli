package auth

import "github.com/go-ldap/ldap/v3"

// not supported yet, to be implemented
type GSSAPIClient struct{}

var _ ldap.GSSAPIClient = (*GSSAPIClient)(nil)

func (*GSSAPIClient) InitSecContext(target string, token []byte) (outputToken []byte, needContinue bool, err error) {
	return nil, false, nil
}

func (*GSSAPIClient) NegotiateSaslAuth(token []byte, authzID string) ([]byte, error) {
	return nil, nil
}

func (*GSSAPIClient) DeleteSecContext() error {
	return nil
}
