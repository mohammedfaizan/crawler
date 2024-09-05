package main


import (
    "testing"
	"strings"
)


func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		errorContains string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
        // add more test cases here


		{
			name:     "http protocol",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},

		{
			name:     "mailto protocol",
			inputURL: "mailto://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},

		{
			name:     "uppercase check",
			inputURL: "mailto://blog.bOOt.dev/path",
			expected: "blog.boot.dev/path",
		},

		{
			name:     "last slash",
			inputURL: "mailto://blog.bOOt.dev/path/",
			expected: "blog.boot.dev/path",
		},

		{
			name:          "handle invalid URL",
			inputURL:      `:\\invalidURL`,
			expected:      "",
			errorContains: "couldn't parse URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}

			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
    tests := []struct {
        name       string
        htmlBody   string
        rawBaseURL string
        expected   []string
		
	}{
		{
			name:       "relative and absolute URLs",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
				<html>
					<body>
						<a href="/path/one">Link 1</a>
						<a href="https://other.com/path/two">Link 2</a>
					</body>
				</html>
			`,
			expected: []string{
				"https://blog.boot.dev/path/one",
				"https://other.com/path/two",
			},
		},

		{
			name:       "no URLs",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody:   "<html><body><p>No links here</p></body></html>",
			expected:   []string{},
		},
		{
			name:       "multiple occurrences",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
				<html>
					<body>
						<a href="/path/one">Link 1</a>
						<a href="/path/one">Link 1 Again</a>
						<a href="https://other.com/path/two">Link 2</a>
					</body>
				</html>
			`,
			expected: []string{
				"https://blog.boot.dev/path/one",
				"https://blog.boot.dev/path/one",
				"https://other.com/path/two",
			},
		},

		{
            name:       "relative URL conversion",
            rawBaseURL: "https://blog.boot.dev",
            htmlBody: `
                <html>
                    <body>
                        <a href="/">Home</a>
                        <a href="/path/page">Page</a>
                        <a href="relative">Relative</a>
                        <a href="//example.com">Protocol-relative</a>
                    </body>
                </html>
            `,
            expected: []string{
                "https://blog.boot.dev/",
                "https://blog.boot.dev/path/page",
                "https://blog.boot.dev/relative",
                "https://example.com",
            },
        },
	}

    for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.htmlBody, tc.rawBaseURL)
			if err != nil {
				t.Fatalf("Test '%s' failed with error: %v", tc.name, err)
			}

			if len(actual) != len(tc.expected) {
                t.Errorf("Test '%s' failed. Expected %d URLs, but got %d", tc.name, len(tc.expected), len(actual))
            } else {
                for i, url := range actual {
                    if url != tc.expected[i] {
                        t.Errorf("Test '%s' failed. URL at index %d doesn't match. Expected %s, but got %s", tc.name, i, tc.expected[i], url)
                    }
                }
            }

		})
	}

}
