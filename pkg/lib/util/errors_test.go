package util

import (
	"errors"
	"testing"

	ldap "github.com/go-ldap/ldap/v3"
)

func TestHandle(t *testing.T) {
	pass := ldap.NewError(321, errors.New("pass"))
	for _, tt := range []struct {
		name string
		args error
		want error
	}{
		{"test#1", ldap.NewError(ldap.LDAPResultSuccess, errors.New("")), nil},
		{"test#2", ldap.NewError(ldap.LDAPResultOperationsError, errors.New("")), ErrOperationFailed},
		{"test#3", ldap.NewError(ldap.LDAPResultProtocolError, errors.New("")), ErrOperationFailed},
		{"test#4", ldap.NewError(ldap.LDAPResultTimeLimitExceeded, errors.New("")), ErrQuotaExceeded},
		{"test#5", ldap.NewError(ldap.LDAPResultSizeLimitExceeded, errors.New("")), ErrQuotaExceeded},
		{"test#6", ldap.NewError(ldap.LDAPResultAuthMethodNotSupported, errors.New("")), ErrAuthenticationFailed},
		{"test#7", ldap.NewError(ldap.LDAPResultStrongAuthRequired, errors.New("")), ErrAuthenticationFailed},
		{"test#8", ldap.NewError(ldap.LDAPResultReferral, errors.New("")), ErrOperationFailed},
		{"test#9", ldap.NewError(ldap.LDAPResultAdminLimitExceeded, errors.New("")), ErrQuotaExceeded},
		{"test#10", ldap.NewError(ldap.LDAPResultUnavailableCriticalExtension, errors.New("")), ErrOperationFailed},
		{"test#11", ldap.NewError(ldap.LDAPResultConfidentialityRequired, errors.New("")), ErrAuthenticationFailed},
		{"test#12", ldap.NewError(ldap.LDAPResultSaslBindInProgress, errors.New("")), ErrAuthenticationFailed},
		{"test#13", ldap.NewError(ldap.LDAPResultNoSuchAttribute, errors.New("")), ErrOperationFailed},
		{"test#14", ldap.NewError(ldap.LDAPResultUndefinedAttributeType, errors.New("")), ErrOperationFailed},
		{"test#15", ldap.NewError(ldap.LDAPResultInappropriateMatching, errors.New("")), ErrOperationFailed},
		{"test#16", ldap.NewError(ldap.LDAPResultConstraintViolation, errors.New("")), ErrOperationFailed},
		{"test#17", ldap.NewError(ldap.LDAPResultAttributeOrValueExists, errors.New("")), nil},
		{"test#18", ldap.NewError(ldap.LDAPResultInvalidAttributeSyntax, errors.New("")), ErrOperationFailed},
		{"test#19", ldap.NewError(ldap.LDAPResultNoSuchObject, errors.New("")), ErrOperationFailed},
		{"test#20", ldap.NewError(ldap.LDAPResultAliasProblem, errors.New("")), ErrOperationFailed},
		{"test#21", ldap.NewError(ldap.LDAPResultInvalidDNSyntax, errors.New("")), ErrOperationFailed},
		{"test#22", ldap.NewError(ldap.LDAPResultIsLeaf, errors.New("")), ErrOperationFailed},
		{"test#23", ldap.NewError(ldap.LDAPResultAliasDereferencingProblem, errors.New("")), ErrOperationFailed},
		{"test#24", ldap.NewError(ldap.LDAPResultInappropriateAuthentication, errors.New("")), ErrAuthenticationFailed},
		{"test#25", ldap.NewError(ldap.LDAPResultInvalidCredentials, errors.New("")), ErrAuthenticationFailed},
		{"test#26", ldap.NewError(ldap.LDAPResultInsufficientAccessRights, errors.New("")), ErrOperationFailed},
		{"test#27", ldap.NewError(ldap.LDAPResultBusy, errors.New("")), ErrOperationFailed},
		{"test#28", ldap.NewError(ldap.LDAPResultUnavailable, errors.New("")), ErrOperationFailed},
		{"test#29", ldap.NewError(ldap.LDAPResultUnwillingToPerform, errors.New("")), ErrOperationFailed},
		{"test#30", ldap.NewError(ldap.LDAPResultLoopDetect, errors.New("")), ErrOperationFailed},
		{"test#31", ldap.NewError(ldap.LDAPResultSortControlMissing, errors.New("")), ErrOperationFailed},
		{"test#32", ldap.NewError(ldap.LDAPResultOffsetRangeError, errors.New("")), ErrOperationFailed},
		{"test#33", ldap.NewError(ldap.LDAPResultNamingViolation, errors.New("")), ErrOperationFailed},
		{"test#34", ldap.NewError(ldap.LDAPResultObjectClassViolation, errors.New("")), ErrOperationFailed},
		{"test#35", ldap.NewError(ldap.LDAPResultNotAllowedOnNonLeaf, errors.New("")), ErrOperationFailed},
		{"test#36", ldap.NewError(ldap.LDAPResultNotAllowedOnRDN, errors.New("")), ErrOperationFailed},
		{"test#37", ldap.NewError(ldap.LDAPResultEntryAlreadyExists, errors.New("")), ErrOperationFailed},
		{"test#38", ldap.NewError(ldap.LDAPResultObjectClassModsProhibited, errors.New("")), ErrOperationFailed},
		{"test#39", ldap.NewError(ldap.LDAPResultResultsTooLarge, errors.New("")), ErrQuotaExceeded},
		{"test#40", ldap.NewError(ldap.LDAPResultAffectsMultipleDSAs, errors.New("")), ErrOperationFailed},
		{"test#41", ldap.NewError(ldap.LDAPResultVirtualListViewErrorOrControlError, errors.New("")), ErrOperationFailed},
		{"test#42", ldap.NewError(ldap.LDAPResultServerDown, errors.New("")), ErrNetworkFailure},
		{"test#43", ldap.NewError(ldap.LDAPResultLocalError, errors.New("")), ErrOperationFailed},
		{"test#44", ldap.NewError(ldap.LDAPResultEncodingError, errors.New("")), ErrOperationFailed},
		{"test#45", ldap.NewError(ldap.LDAPResultDecodingError, errors.New("")), ErrOperationFailed},
		{"test#46", ldap.NewError(ldap.LDAPResultTimeout, errors.New("")), ErrQuotaExceeded},
		{"test#47", ldap.NewError(ldap.LDAPResultAuthUnknown, errors.New("")), ErrAuthenticationFailed},
		{"test#48", ldap.NewError(ldap.LDAPResultFilterError, errors.New("")), ErrInvalidFilter},
		{"test#49", ldap.NewError(ldap.LDAPResultUserCanceled, errors.New("")), ErrOperationFailed},
		{"test#50", ldap.NewError(ldap.LDAPResultParamError, errors.New("")), ErrOperationFailed},
		{"test#51", ldap.NewError(ldap.LDAPResultNoMemory, errors.New("")), ErrQuotaExceeded},
		{"test#52", ldap.NewError(ldap.LDAPResultConnectError, errors.New("")), ErrNetworkFailure},
		{"test#53", ldap.NewError(ldap.LDAPResultNotSupported, errors.New("")), ErrOperationFailed},
		{"test#54", ldap.NewError(ldap.LDAPResultControlNotFound, errors.New("")), ErrOperationFailed},
		{"test#55", ldap.NewError(ldap.LDAPResultNoResultsReturned, errors.New("")), ErrNothingReturned},
		{"test#56", ldap.NewError(ldap.LDAPResultMoreResultsToReturn, errors.New("")), nil},
		{"test#57", ldap.NewError(ldap.LDAPResultClientLoop, errors.New("")), ErrOperationFailed},
		{"test#58", ldap.NewError(ldap.LDAPResultReferralLimitExceeded, errors.New("")), ErrQuotaExceeded},
		{"test#59", ldap.NewError(ldap.LDAPResultInvalidResponse, errors.New("")), ErrOperationFailed},
		{"test#60", ldap.NewError(ldap.LDAPResultAmbiguousResponse, errors.New("")), ErrOperationFailed},
		{"test#61", ldap.NewError(ldap.LDAPResultTLSNotSupported, errors.New("")), ErrAuthenticationFailed},
		{"test#62", ldap.NewError(ldap.LDAPResultUnknownType, errors.New("")), ErrOperationFailed},
		{"test#63", ldap.NewError(ldap.LDAPResultCanceled, errors.New("")), ErrOperationFailed},
		{"test#64", ldap.NewError(ldap.LDAPResultNoSuchOperation, errors.New("")), ErrOperationFailed},
		{"test#65", ldap.NewError(ldap.LDAPResultTooLate, errors.New("")), ErrOperationFailed},
		{"test#66", ldap.NewError(ldap.LDAPResultCannotCancel, errors.New("")), ErrOperationFailed},
		{"test#67", ldap.NewError(ldap.LDAPResultAssertionFailed, errors.New("")), ErrOperationFailed},
		{"test#68", ldap.NewError(ldap.LDAPResultAuthorizationDenied, errors.New("")), ErrAuthenticationFailed},
		{"test#69", ldap.NewError(ldap.LDAPResultSyncRefreshRequired, errors.New("")), ErrOperationFailed},
		{"test#70", ldap.NewError(ldap.ErrorNetwork, errors.New("")), ErrNetworkFailure},
		{"test#71", ldap.NewError(ldap.ErrorFilterCompile, errors.New("")), ErrInvalidFilter},
		{"test#72", ldap.NewError(ldap.ErrorFilterDecompile, errors.New("")), ErrInvalidFilter},
		{"test#73", ldap.NewError(ldap.ErrorDebugging, errors.New("")), ErrOperationFailed},
		{"test#74", ldap.NewError(ldap.ErrorUnexpectedMessage, errors.New("")), ErrOperationFailed},
		{"test#75", ldap.NewError(ldap.ErrorUnexpectedResponse, errors.New("")), ErrOperationFailed},
		{"test#76", ldap.NewError(ldap.ErrorEmptyPassword, errors.New("")), ErrAuthenticationFailed},
		{"test#77", nil, nil},
		{"test#78", pass, pass},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Handle(tt.args)
			if !errors.Is(got, tt.want) {
				t.Errorf(`Handle(%v) failed: got: %v, want: %v`, tt.args, got, tt.want)
			}
		})
	}
}
