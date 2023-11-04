package client

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestSearch(t *testing.T) {
	libutil.SkipOAT(t)

	user := os.Getenv("AD_DEFAULT_USER")
	conn, err := auth.Bind(
		auth.NewBindParameters().SetType(auth.SIMPLE).SetUser(user).SetPassword(os.Getenv("AD_DEFAULT_PASS")),
		auth.NewDialOptions().SetURL(os.Getenv("AD_CW01_URL")).SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).SetSizeLimit(-1).SetTimeLimit(time.Minute*5),
	)
	if err != nil {
		t.Error(err)
	}

	dom, uid, _ := strings.Cut(user, `\\`)
	t.Log(uid, dom)

	result, _, err := Search(
		conn,
		SearchArguments{
			Attributes: attributes.Attributes{attributes.Any()},
			Path:       fmt.Sprintf("DC=%s,DC=contiwan,DC=com", dom),
			Filter:     filter.Filter{Attribute: attributes.SamAccountName(), Value: uid},
		},
	)
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
