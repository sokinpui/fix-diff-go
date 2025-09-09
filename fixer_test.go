package fixdiff

import (
	"strings"
	"testing"
)

const originalFile = `line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10`

const incorrectDiff = `--- a/file.txt
+++ b/file.txt
@@ -99,5 +99,5 @@
 line 2
 line 3
-line 4
+line four
 line 5
 line 6
`

const expectedCorrectDiff = `--- a/file.txt
+++ b/file.txt
@@ -1,6 +1,6 @@
 line 1
 line 2
 line 3
-line 4
+line four
 line 5
 line 6
`

func TestFix(t *testing.T) {
	t.Run("corrects a diff with incorrect hunk headers", func(t *testing.T) {
		got, err := Fix(incorrectDiff, originalFile)
		if err != nil {
			t.Fatalf("Fix() returned an unexpected error: %v", err)
		}

		// Normalize newlines and compare line by line for easier debugging.
		gotLines := strings.Split(strings.ReplaceAll(got, "\r\n", "\n"), "\n")
		expectedLines := strings.Split(strings.ReplaceAll(expectedCorrectDiff, "\r\n", "\n"), "\n")

		if len(gotLines) != len(expectedLines) {
			t.Fatalf("Fix() returned diff with wrong number of lines.\ngot:\n%s\n\nwant:\n%s", got, expectedCorrectDiff)
		}

		for i := range gotLines {
			if gotLines[i] != expectedLines[i] {
				t.Errorf("Line %d mismatch:\ngot:  %q\nwant: %q", i+1, gotLines[i], expectedLines[i])
			}
		}
	})

	t.Run("returns error for unpatchable diff", func(t *testing.T) {
		unpatchableDiff := `--- a/file.txt
+++ b/file.txt
@@ -1,3 +1,3 @@
-nonexistent line 1
+new line 1
`
		_, err := Fix(unpatchableDiff, originalFile)
		if err == nil {
			t.Error("Fix() was expected to return an error for an unpatchable diff, but it did not")
		}
	})

	t.Run("returns error for invalid diff format", func(t *testing.T) {
		invalidDiff := `this is not a diff`
		_, err := Fix(invalidDiff, originalFile)
		if err == nil {
			t.Error("Fix() was expected to return an error for an invalid diff format, but it did not")
		}
	})
}
