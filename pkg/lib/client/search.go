package client

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	ldap "github.com/go-ldap/ldap/v3"
	ldif "github.com/go-ldap/ldif"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
	progressbar "github.com/schollz/progressbar/v3"
)

// searchRangeRegex is a regular expression to match ranged attributes
var searchRangeRegex = regexp.MustCompile(`(\w+);range\=(?P<from>\d+)\-(?P<to>\d+)`)

// searchRecursivelyArguments is a struct to hold arguments for recursive search
type searchRecursivelyArguments struct {
	From, To, ID        int
	Path, AttributeName string
	Filter              filter.Filter
	Repeat              bool
}

// searchRecursively searches for ranged attributes recursively
func searchRecursively(conn *auth.Connection, args searchRecursivelyArguments, result map[string][]string) error {
	var rName string
	standardizedName := libutil.ToTitleNoLower(args.AttributeName)
	if !args.Repeat {
		// to retrieve remaining attribute values,
		// the limit must be enforced by server and not by the requestor
		rName = fmt.Sprintf("%s;range=%d-*", args.AttributeName, args.To+1)
	} else {
		// enforce limit as requestor
		rName = fmt.Sprintf("%s;range=%d-%d", args.AttributeName, args.To+1, 2*args.To-args.From+1)
	}
	if sr, err := conn.SearchWithPaging(ldap.NewSearchRequest(
		args.Path, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		int(conn.SizeLimit),           // SizeLimit
		int(conn.TimeLimit.Seconds()), // TimeLimit
		false,                         // TypesOnly
		args.Filter.String(),          // LDAP Filter
		[]string{rName},               // Attribute List
		nil,                           // []ldap.Control
	), 1000); err == nil {

		found := sr.Entries[args.ID].GetAttributeValues(rName)

		// make recursive call if count equals range boundaries
		if len(found)-1 == args.To-args.From {
			if err := searchRecursively(conn, searchRecursivelyArguments{
				From:          args.To + 1,
				To:            2*args.To - args.From + 1,
				ID:            args.ID,
				Path:          args.Path,
				AttributeName: args.AttributeName,
				Filter:        args.Filter,
				Repeat:        true,
			}, result); err != nil {

				return err
			}

		} else {
			// the count was 0, indicating the end of the range
			// to retrieve remaining members, invoke request without
			// limiting the upper range boundary
			if args.Repeat { // prevent endless looping :)
				if err := searchRecursively(conn, searchRecursivelyArguments{
					From:          args.From,
					To:            args.To,
					ID:            args.ID,
					Path:          args.Path,
					AttributeName: args.AttributeName,
					Filter:        args.Filter,
					Repeat:        false,
				}, result); err != nil {

					return err
				}

			}
		}

		if val, ok := result[standardizedName]; ok {
			result[standardizedName] = append(val, found...)
		} else {
			result[standardizedName] = found
		}

	} else {
		return err
	}

	return nil
}

// SearchArguments is a struct to hold arguments for search
type SearchArguments struct {
	Path       string
	Attributes attributes.Attributes
	Filter     filter.Filter
}

// Search searches for entries in AD
//
//gocyclo:ignore
func Search(conn *auth.Connection, args SearchArguments, bar *progressbar.ProgressBar) (results attributes.Maps, requests *ldif.LDIF, err error) {
	if conn == nil {
		return nil, nil, ldap.ErrNilConnection
	}

	if bar != nil {
		bar.Describe("searching")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for id, pagingControl := 0, ldap.NewControlPaging(1); true; id++ {
		request := ldap.NewSearchRequest(
			args.Path, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
			int(conn.SizeLimit),               // SizeLimit
			int(conn.TimeLimit.Seconds()),     // TimeLimit
			false,                             // TypesOnly
			args.Filter.String(),              // LDAP Filter
			args.Attributes.ToAttributeList(), // Attribute List
			[]ldap.Control{pagingControl},     // []ldap.Control
		)

		sr := conn.SearchAsync(ctx, request, 1)
		requests = &ldif.LDIF{}
		for s := (*ldap.Entry)(nil); sr.Next(); s = sr.Entry() {
			if s == nil {
				continue
			}

			if bar != nil {
				bar.Describe(fmt.Sprintf("[%d]: found: %q", id, s.DN))
				_ = bar.Add(1)
			}

			requests.Entries = append(requests.Entries, &ldif.Entry{Entry: s})
			result := make(map[string][]string)
			converted := make(attributes.Map)
			for _, attr := range s.Attributes {
				// retrieve ranged attribute recursively
				if searchRangeRegex.MatchString(attr.Name) {
					from, _ := strconv.Atoi(searchRangeRegex.FindStringSubmatch(attr.Name)[2])
					to, _ := strconv.Atoi(searchRangeRegex.FindStringSubmatch(attr.Name)[3])
					if err := searchRecursively(conn, searchRecursivelyArguments{
						From:          from,
						To:            to,
						ID:            id,
						Path:          args.Path,
						AttributeName: searchRangeRegex.FindStringSubmatch(attr.Name)[1],
						Filter:        args.Filter,
						Repeat:        true,
					}, result); err != nil {

						return nil, nil, err
					}

				} else {
					result[attr.Name] = attr.Values
				}
			}

			for k, v := range result {
				// parse registered attribute
				if attr := attributes.Lookup(k); attr != nil {
					attr.Parse(v, &converted)
					continue
				}

				// parse unknown attributes
				attributes.Raw(libutil.ToTitleNoLower(k), "", attributes.TypeRaw).Parse(v, &converted)

			}

			results = append(results, converted)
		}

		if err := sr.Err(); err != nil {
			return nil, nil, err
		}

		if ctrl, ok := ldap.FindControl(sr.Controls(), ldap.ControlTypePaging).(*ldap.ControlPaging); ok && ctrl != nil && len(ctrl.Cookie) > 0 {
			pagingControl.SetCookie(ctrl.Cookie)
		} else {
			break
		}
	}

	return
}
