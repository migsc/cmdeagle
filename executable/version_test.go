package executable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckVersionCompatibility(t *testing.T) {
	tests := []struct {
		name            string
		foundVersion    string
		declaredVersion string
		wantCompatible  bool
		wantReason      string
	}{
		// Wildcard tests
		{
			name:            "wildcard accepts any version",
			foundVersion:    "1.2.3",
			declaredVersion: "*",
			wantCompatible:  true,
			wantReason:      "",
		},

		// Caret (^) tests
		{
			name:            "caret - compatible major version",
			foundVersion:    "1.2.3",
			declaredVersion: "^1.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "caret - incompatible major version",
			foundVersion:    "2.0.0",
			declaredVersion: "^1.0.0",
			wantCompatible:  false,
			wantReason:      "Major version mismatch",
		},

		// Tilde (~) tests
		{
			name:            "tilde - compatible patch version",
			foundVersion:    "1.2.3",
			declaredVersion: "~1.2.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "tilde - incompatible minor version",
			foundVersion:    "1.3.0",
			declaredVersion: "~1.2.0",
			wantCompatible:  false,
			wantReason:      "Major or minor version mismatch",
		},

		// Greater than or equal tests
		{
			name:            "greater than or equal - valid",
			foundVersion:    "2.0.0",
			declaredVersion: ">=1.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "greater than or equal - invalid",
			foundVersion:    "1.0.0",
			declaredVersion: ">=2.0.0",
			wantCompatible:  false,
			wantReason:      "Version too low",
		},

		// Greater than tests
		{
			name:            "greater than - valid",
			foundVersion:    "2.0.0",
			declaredVersion: ">1.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "greater than - invalid",
			foundVersion:    "1.0.0",
			declaredVersion: ">1.0.0",
			wantCompatible:  false,
			wantReason:      "Version too low",
		},

		// Less than or equal tests
		{
			name:            "less than or equal - valid",
			foundVersion:    "1.0.0",
			declaredVersion: "<=2.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "less than or equal - invalid",
			foundVersion:    "3.0.0",
			declaredVersion: "<=2.0.0",
			wantCompatible:  false,
			wantReason:      "Version too high",
		},

		// Less than tests
		{
			name:            "less than - valid",
			foundVersion:    "1.0.0",
			declaredVersion: "<2.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "less than - invalid",
			foundVersion:    "2.0.0",
			declaredVersion: "<2.0.0",
			wantCompatible:  false,
			wantReason:      "Version too high",
		},

		// Range tests
		{
			name:            "range - version within range",
			foundVersion:    "1.5.0",
			declaredVersion: "1.0.0 - 2.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "range - version below range",
			foundVersion:    "0.9.0",
			declaredVersion: "1.0.0 - 2.0.0",
			wantCompatible:  false,
			wantReason:      "Version below range",
		},
		{
			name:            "range - version above range",
			foundVersion:    "2.1.0",
			declaredVersion: "1.0.0 - 2.0.0",
			wantCompatible:  false,
			wantReason:      "Version above range",
		},

		// Exact version tests
		{
			name:            "exact version match",
			foundVersion:    "1.2.3",
			declaredVersion: "1.2.3",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "exact version mismatch",
			foundVersion:    "1.2.4",
			declaredVersion: "1.2.3",
			wantCompatible:  false,
			wantReason:      "Version does not match exactly",
		},

		// Pre-release version tests
		{
			name:            "pre-release version comparison",
			foundVersion:    "1.2.3-beta",
			declaredVersion: "1.2.3",
			wantCompatible:  true,
			wantReason:      "",
		},

		// Invalid format tests
		{
			name:            "invalid found version format",
			foundVersion:    "1.2",
			declaredVersion: "1.2.3",
			wantCompatible:  false,
			wantReason:      "Invalid found version format",
		},
		{
			name:            "invalid declared version format",
			foundVersion:    "1.2.3",
			declaredVersion: "1.2",
			wantCompatible:  false,
			wantReason:      "Invalid declared version format",
		},

		// Double-digit version tests
		{
			name:            "greater than or equal - double digit major version",
			foundVersion:    "10.2.3",
			declaredVersion: ">=9.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "greater than or equal - double digit minor version",
			foundVersion:    "1.10.3",
			declaredVersion: ">=1.9.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "greater than or equal - double digit patch version",
			foundVersion:    "1.2.10",
			declaredVersion: ">=1.2.9",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "less than or equal - double digit major version",
			foundVersion:    "9.2.3",
			declaredVersion: "<=10.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "less than or equal - double digit minor version",
			foundVersion:    "1.9.3",
			declaredVersion: "<=1.10.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "less than or equal - double digit patch version",
			foundVersion:    "1.2.9",
			declaredVersion: "<=1.2.10",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "range - double digit versions within range",
			foundVersion:    "10.5.0",
			declaredVersion: "9.0.0 - 11.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "range - double digit version below range",
			foundVersion:    "9.9.9",
			declaredVersion: "10.0.0 - 11.0.0",
			wantCompatible:  false,
			wantReason:      "Version below range",
		},
		{
			name:            "range - double digit version above range",
			foundVersion:    "12.0.0",
			declaredVersion: "10.0.0 - 11.0.0",
			wantCompatible:  false,
			wantReason:      "Version above range",
		},
		{
			name:            "caret - compatible double digit major version",
			foundVersion:    "10.2.3",
			declaredVersion: "^10.0.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "caret - incompatible double digit major version",
			foundVersion:    "11.0.0",
			declaredVersion: "^10.0.0",
			wantCompatible:  false,
			wantReason:      "Major version mismatch",
		},
		{
			name:            "tilde - compatible double digit minor version",
			foundVersion:    "1.10.3",
			declaredVersion: "~1.10.0",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "tilde - incompatible double digit minor version",
			foundVersion:    "1.11.0",
			declaredVersion: "~1.10.0",
			wantCompatible:  false,
			wantReason:      "Major or minor version mismatch",
		},
		{
			name:            "exact match - all double digit versions",
			foundVersion:    "10.11.12",
			declaredVersion: "10.11.12",
			wantCompatible:  true,
			wantReason:      "",
		},
		{
			name:            "exact match - double digit version mismatch",
			foundVersion:    "10.11.12",
			declaredVersion: "10.11.13",
			wantCompatible:  false,
			wantReason:      "Version does not match exactly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCompatible, gotReason := CheckVersionCompatibility(tt.foundVersion, tt.declaredVersion)
			assert.Equal(t, tt.wantCompatible, gotCompatible, "compatibility mismatch")
			assert.Equal(t, tt.wantReason, gotReason, "reason mismatch")
		})
	}
}
