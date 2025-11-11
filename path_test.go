package mappath_test

import (
	"testing"

	"github.com/gekatateam/mappath"
)

func TestCurrentKey(t *testing.T) {
	tests := map[string]struct {
		path         string
		expectedCurr string
		expectedLeft string
	}{
		// Basic cases
		"simple path": {
			path:         `a.b.c`,
			expectedCurr: `a`,
			expectedLeft: `b.c`,
		},
		"single node": {
			path:         `single`,
			expectedCurr: `single`,
			expectedLeft: ``,
		},
		"empty string": {
			path:         ``,
			expectedCurr: ``,
			expectedLeft: ``,
		},

		// Escaped dot
		"escaped dot in middle": {
			path:         `first\.part.second`,
			expectedCurr: `first.part`,
			expectedLeft: `second`,
		},
		"multiple escaped dots": {
			path:         `a\.b\.c.d`,
			expectedCurr: `a.b.c`,
			expectedLeft: `d`,
		},
		"escaped dot at end of node": {
			path:         `test\..next.and.another`,
			expectedCurr: `test.`,
			expectedLeft: `next.and.another`,
		},
		"all escaped dots": {
			path:         `a\.b\.c`,
			expectedCurr: `a.b.c`,
			expectedLeft: ``,
		},

		// Escaped backslash
		"escaped backslash": {
			path:         `path\\with.backslash`,
			expectedCurr: `path\with`,
			expectedLeft: `backslash`,
		},
		"mixed escaping": {
			path:         `mixed\\\.escaping.next`,
			expectedCurr: `mixed\.escaping`,
			expectedLeft: `next`,
		},

		// Some corner cases
		"backslash at end": {
			path:         `end\\`,
			expectedCurr: `end\`,
			expectedLeft: ``,
		},
		"only escaped dot": {
			path:         `\.`,
			expectedCurr: `.`,
			expectedLeft: ``,
		},
		"only backslash": {
			path:         `\\`,
			expectedCurr: `\`,
			expectedLeft: ``,
		},
		"path starting with dot": {
			path:         `.start`,
			expectedCurr: ``,
			expectedLeft: `start`,
		},
		"path starting with escapsed dot": {
			path:         `\.start`,
			expectedCurr: `.start`,
			expectedLeft: ``,
		},
		"path ending with dot": {
			path:         `end.`,
			expectedCurr: `end`,
			expectedLeft: ``,
		},
		"path ending with escapsed dot": {
			path:         `end\.`,
			expectedCurr: `end.`,
			expectedLeft: ``,
		},
		"consecutive dots": {
			path:         `a..b`,
			expectedCurr: `a`,
			expectedLeft: `.b`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			curr, left := mappath.CurrentKey(test.path)

			if curr != test.expectedCurr {
				t.Errorf("unexpected current with input %q - want: %q, got: %q", test.path, curr, test.expectedCurr)
			}
			if left != test.expectedLeft {
				t.Errorf("unexpected left with input %q - want: %q, got: %q", test.path, left, test.expectedLeft)
			}
		})
	}
}
