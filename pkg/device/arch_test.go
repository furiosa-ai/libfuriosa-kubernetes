package device

import (
	"errors"
	"testing"
)

func TestArchFromStr(t *testing.T) {
	tests := []struct {
		description  string
		arch         string
		rev          string
		expectedArch Arch
		expectedErr  error
	}{
		{
			description:  "WarboyA0",
			arch:         "Warboy",
			rev:          "A0",
			expectedArch: ArchWarboy,
			expectedErr:  nil,
		},
		{
			description:  "WarboyB0",
			arch:         "Warboy",
			rev:          "B0",
			expectedArch: ArchWarboy,
			expectedErr:  nil,
		},
		{
			description:  "Renegade",
			arch:         "Renegade",
			rev:          "",
			expectedArch: ArchRenegade,
			expectedErr:  nil,
		},
		{
			description:  "Wrong arch",
			arch:         "Wrong",
			rev:          "Wrong",
			expectedArch: "",
			expectedErr:  errors.New("unknown arch"),
		},
	}

	for _, tc := range tests {
		actualArch, actualErr := archFromStr(tc.arch, tc.rev)
		if tc.expectedErr != nil || actualErr != nil {
			//NOTE(bg): comparing error message is terrible idea, but this is only for testing.
			// we return typed error(UnknownArch) at the outside of this function.
			if tc.expectedErr.Error() != actualErr.Error() {
				t.Errorf("expected %s but got %s", tc.expectedErr, actualErr)
				continue
			}
		}

		if tc.expectedArch != actualArch {
			t.Errorf("expected %s but got %s", tc.expectedArch, actualArch)
			continue
		}
	}
}
