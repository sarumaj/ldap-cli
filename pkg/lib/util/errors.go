package util

import (
	"errors"

	ldap "github.com/go-ldap/ldap/v3"
)

var ErrInvalidFilter = errors.New("invalid filter syntax")
var ErrOperationFailed = errors.New("query failed")
var ErrQuotaExceeded = errors.New("quota limit exceed")
var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrNetworkFailure = errors.New("network error")
var ErrNothingReturned = errors.New("nothing returned")
var ErrUbiquitousResults = errors.New("ubiquitous response")

var errorMapping = map[uint16]error{
	ldap.LDAPResultSuccess:                            nil,
	ldap.LDAPResultOperationsError:                    ErrOperationFailed,
	ldap.LDAPResultProtocolError:                      ErrOperationFailed,
	ldap.LDAPResultTimeLimitExceeded:                  ErrQuotaExceeded,
	ldap.LDAPResultSizeLimitExceeded:                  ErrQuotaExceeded,
	ldap.LDAPResultAuthMethodNotSupported:             ErrAuthenticationFailed,
	ldap.LDAPResultStrongAuthRequired:                 ErrAuthenticationFailed,
	ldap.LDAPResultReferral:                           ErrOperationFailed,
	ldap.LDAPResultAdminLimitExceeded:                 ErrQuotaExceeded,
	ldap.LDAPResultUnavailableCriticalExtension:       ErrOperationFailed,
	ldap.LDAPResultConfidentialityRequired:            ErrAuthenticationFailed,
	ldap.LDAPResultSaslBindInProgress:                 ErrAuthenticationFailed,
	ldap.LDAPResultNoSuchAttribute:                    ErrOperationFailed,
	ldap.LDAPResultUndefinedAttributeType:             ErrOperationFailed,
	ldap.LDAPResultInappropriateMatching:              ErrOperationFailed,
	ldap.LDAPResultConstraintViolation:                ErrOperationFailed,
	ldap.LDAPResultAttributeOrValueExists:             nil,
	ldap.LDAPResultInvalidAttributeSyntax:             ErrOperationFailed,
	ldap.LDAPResultNoSuchObject:                       ErrOperationFailed,
	ldap.LDAPResultAliasProblem:                       ErrOperationFailed,
	ldap.LDAPResultInvalidDNSyntax:                    ErrOperationFailed,
	ldap.LDAPResultIsLeaf:                             ErrOperationFailed,
	ldap.LDAPResultAliasDereferencingProblem:          ErrOperationFailed,
	ldap.LDAPResultInappropriateAuthentication:        ErrAuthenticationFailed,
	ldap.LDAPResultInvalidCredentials:                 ErrAuthenticationFailed,
	ldap.LDAPResultInsufficientAccessRights:           ErrOperationFailed,
	ldap.LDAPResultBusy:                               ErrOperationFailed,
	ldap.LDAPResultUnavailable:                        ErrOperationFailed,
	ldap.LDAPResultUnwillingToPerform:                 ErrOperationFailed,
	ldap.LDAPResultLoopDetect:                         ErrOperationFailed,
	ldap.LDAPResultSortControlMissing:                 ErrOperationFailed,
	ldap.LDAPResultOffsetRangeError:                   ErrOperationFailed,
	ldap.LDAPResultNamingViolation:                    ErrOperationFailed,
	ldap.LDAPResultObjectClassViolation:               ErrOperationFailed,
	ldap.LDAPResultNotAllowedOnNonLeaf:                ErrOperationFailed,
	ldap.LDAPResultNotAllowedOnRDN:                    ErrOperationFailed,
	ldap.LDAPResultEntryAlreadyExists:                 ErrOperationFailed,
	ldap.LDAPResultObjectClassModsProhibited:          ErrOperationFailed,
	ldap.LDAPResultResultsTooLarge:                    ErrQuotaExceeded,
	ldap.LDAPResultAffectsMultipleDSAs:                ErrOperationFailed,
	ldap.LDAPResultVirtualListViewErrorOrControlError: ErrOperationFailed,
	ldap.LDAPResultServerDown:                         ErrNetworkFailure,
	ldap.LDAPResultLocalError:                         ErrOperationFailed,
	ldap.LDAPResultEncodingError:                      ErrOperationFailed,
	ldap.LDAPResultDecodingError:                      ErrOperationFailed,
	ldap.LDAPResultTimeout:                            ErrQuotaExceeded,
	ldap.LDAPResultAuthUnknown:                        ErrAuthenticationFailed,
	ldap.LDAPResultFilterError:                        ErrInvalidFilter,
	ldap.LDAPResultUserCanceled:                       ErrOperationFailed,
	ldap.LDAPResultParamError:                         ErrOperationFailed,
	ldap.LDAPResultNoMemory:                           ErrQuotaExceeded,
	ldap.LDAPResultConnectError:                       ErrNetworkFailure,
	ldap.LDAPResultNotSupported:                       ErrOperationFailed,
	ldap.LDAPResultControlNotFound:                    ErrOperationFailed,
	ldap.LDAPResultNoResultsReturned:                  ErrNothingReturned,
	ldap.LDAPResultMoreResultsToReturn:                nil,
	ldap.LDAPResultClientLoop:                         ErrOperationFailed,
	ldap.LDAPResultReferralLimitExceeded:              ErrQuotaExceeded,
	ldap.LDAPResultInvalidResponse:                    ErrOperationFailed,
	ldap.LDAPResultAmbiguousResponse:                  ErrOperationFailed,
	ldap.LDAPResultTLSNotSupported:                    ErrAuthenticationFailed,
	ldap.LDAPResultUnknownType:                        ErrOperationFailed,
	ldap.LDAPResultCanceled:                           ErrOperationFailed,
	ldap.LDAPResultNoSuchOperation:                    ErrOperationFailed,
	ldap.LDAPResultTooLate:                            ErrOperationFailed,
	ldap.LDAPResultCannotCancel:                       ErrOperationFailed,
	ldap.LDAPResultAssertionFailed:                    ErrOperationFailed,
	ldap.LDAPResultAuthorizationDenied:                ErrAuthenticationFailed,
	ldap.LDAPResultSyncRefreshRequired:                ErrOperationFailed,
	ldap.ErrorNetwork:                                 ErrNetworkFailure,
	ldap.ErrorFilterCompile:                           ErrInvalidFilter,
	ldap.ErrorFilterDecompile:                         ErrInvalidFilter,
	ldap.ErrorDebugging:                               ErrOperationFailed,
	ldap.ErrorUnexpectedMessage:                       ErrOperationFailed,
	ldap.ErrorUnexpectedResponse:                      ErrOperationFailed,
	ldap.ErrorEmptyPassword:                           ErrAuthenticationFailed,
}

// Handle maps LDAP errors to more specific errors
func Handle(err error) error {
	if err == nil {
		return nil
	}

	for code, v := range errorMapping {
		if ldap.IsErrorWithCode(err, code) {
			if v != nil {
				return ldap.NewError(code, errors.Join(v, errors.New(ldap.LDAPResultCodeMap[code]), err))
			}

			return nil
		}
	}

	return err
}
