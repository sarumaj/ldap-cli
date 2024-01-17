package attributes

// LDAP Matching Rules (extensibleMatch): https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/4e638665-f466-4597-93c4-12f2ebfabab5
const (
	LDAP_MATCHING_RULE_BIT_AND         MatchingRule = "1.2.840.113556.1.4.803"  // https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/6dd1d7b4-2b2f-4e55-b164-7047c4c5bb00
	LDAP_MATCHING_RULE_BIT_OR          MatchingRule = "1.2.840.113556.1.4.804"  // https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/4e5b2424-642a-40da-acb1-9fff381b46e4
	LDAP_MATCHING_RULE_IN_CHAIN        MatchingRule = "1.2.840.113556.1.4.1941" // https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/1e889adc-b503-4423-8985-c28d5c7d4887
	LDAP_MATCHING_RULE_TRANSITIVE_EVAL MatchingRule = "1.2.840.113556.1.4.1941" // https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/1e889adc-b503-4423-8985-c28d5c7d4887
	LDAP_MATCHING_RULE_DN_WITH_DATA    MatchingRule = "1.2.840.113556.1.4.2253" // https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/e5bfc285-05b9-494e-a123-c5c4341c450e
)

// MatchingRule is used to define an LDAP matching rule bit mask
type MatchingRule string
