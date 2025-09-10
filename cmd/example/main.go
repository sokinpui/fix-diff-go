package main

import (
	"fmt"
	"log"

	fixdiff "github.com/sokinpui/fix-diff-go"
)

// originalFileContent represents the complete, original content of the source file.
const originalFileContent = `line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10`

// incorrectDiffContent is a unified diff with incorrect hunk headers.
// For example, `@@ -99,5 +99,5 @@` is wrong, as the change happens near the top of the file.
const incorrectDiffContent = `--- a/test/file.txt
+++ b/test/file.txt
@@ -99,5 +99,5 @@
 line 2
 line 3
-line 4
+line four
 line 5
 line 6
`

func main() {
	// Call the Fix function with the incorrect diff and the original file content.
	correctedDiff, err := fixdiff.Fix(incorrectDiffContent, originalFileContent)
	if err != nil {
		log.Fatalf("Failed to fix the diff: %v", err)
	}

	// Print the newly generated, correct diff.
	fmt.Println("--- Original Incorrect Diff ---")
	fmt.Println(incorrectDiffContent)
	fmt.Println("\n--- Corrected Diff ---")
	fmt.Println(correctedDiff)
}
