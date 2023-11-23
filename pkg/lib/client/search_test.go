package client

import (
	"testing"
	"time"

	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestSearch(t *testing.T) {
	libutil.SkipOAT(t)

	conn, err := auth.Bind(
		auth.NewBindParameters().SetType(auth.SIMPLE).SetUser("cn=admin,dc=mock,dc=ad,dc=com").SetPassword("admin"),
		auth.NewDialOptions().SetSizeLimit(10).SetTimeLimit(time.Minute*5),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	result, _, err := Search(
		conn,
		SearchArguments{
			Attributes: attributes.Attributes{attributes.Any()},
			Path:       "dc=mock,dc=ad,dc=com",
			Filter:     filter.Filter{Attribute: attributes.CommonName(), Value: "uix00001"},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Error("empty result set")
	}

	result, _, err = Search(
		conn,
		SearchArguments{
			Attributes: attributes.Attributes{attributes.Any()},
			Path:       "dc=mock,dc=ad,dc=com",
			Filter:     filter.Filter{Attribute: attributes.CommonName(), Value: "group01"},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Error("empty result set")
	}
}
