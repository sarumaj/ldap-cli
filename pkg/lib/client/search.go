package client

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/auth"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

var searchRangeRegex = regexp.MustCompile(`(\w+);range\=(?P<from>\d+)\-(?P<to>\d+)`)

type searchRecursivelyArguments struct {
	From, To, ID        int
	Path, AttributeName string
	Filter              filter.Filter
	Repeat              bool
}

func searchRecursively(conn *auth.Connection, args searchRecursivelyArguments, result map[string][]string) error {
	var rName string
	standardizedName := strings.ToLower(args.AttributeName)
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

type SearchArguments struct {
	Path       string
	Attributes attributes.Attributes
	Filter     filter.Filter
}

func Search(conn *auth.Connection, args SearchArguments) (results attributes.Maps, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	if conn == nil {
		return nil, util.ErrBindFirst
	}

	sr, err := conn.SearchWithPaging(ldap.NewSearchRequest(
		args.Path, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		int(conn.SizeLimit),               // SizeLimit
		int(conn.TimeLimit.Seconds()),     // TimeLimit
		false,                             // TypesOnly
		args.Filter.String(),              // LDAP Filter
		args.Attributes.ToAttributeList(), // Attribute List
		nil,                               // []ldap.Control
	), 1000)
	if err != nil {
		return nil, err
	}

	for id, s := range sr.Entries {
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

					return nil, err
				}

			} else {
				result[strings.ToLower(attr.Name)] = attr.Values
			}
		}

		for k, v := range result {
			attr := attributes.Lookup(k)
			if attr == nil {
				continue
			}

			attr.Parse(v, &converted)
		}

		results = append(results, converted)
	}

	return
}
