package main

import (
	"fmt"
	"os"

	fixdiff "github.com/sokinpui/fix-diff-go"
)

// originalFileContent represents the complete, original content of the source file.
// read from file "patcher.go"

// incorrectDiffContent is a unified diff with incorrect hunk headers.
// For example, `@@ -99,5 +99,5 @@` is wrong, as the change happens near the top of the file.
const incorrectDiffContent = `--- a/patcher
+++ b/patcher
@@ -8,6 +8,8 @@
 	"path/filepath"
 	"regexp"
 	"strings"
+
+	"github.com/sokinpui/fix-diff-go"
 )

 // filePathRegex extracts the file path from a '+++ b/...' line.
@@ -52,16 +54,15 @@
 // CorrectDiff prepares a valid patch from a raw diff block.
 func CorrectDiff(diff model.DiffBlock, resolver *fs.PathResolver, extensions []string) (string, error) {
 	sourcePath := resolver.ResolveExisting(diff.FilePath)
-	var sourceLines []string
+	var sourceContent string
 	if sourcePath != "" {
 		content, err := os.ReadFile(sourcePath)
 		if err == nil {
-			sourceLines = strings.Split(string(content), "\n")
+			sourceContent = string(content)
 		}
 	}

-	return correctDiffHunks(sourceLines, diff.RawContent, diff.FilePath)
+	return fixdiff.Fix(diff.RawContent, sourceContent)
 }

 func applyPatch(filePath, patchContent string, resolver *fs.PathResolver) ([]string, error) {
`

func main() {
	// Call the Fix function with the incorrect diff and the original file content.

	originalFileContent, err := os.ReadFile("patcher")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read original file 'patcher.go': %v\n", err)
		os.Exit(1)
	}

	correctedDiff, err := fixdiff.Fix(incorrectDiffContent, string(originalFileContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fix the diff: %v\n", err)
		os.Exit(1)
	}

	// Print the newly generated, correct diff.
	fmt.Println("--- Original Incorrect Diff ---")
	fmt.Println(incorrectDiffContent)
	fmt.Println("--- Corrected Diff ---")
	fmt.Println(correctedDiff)
}
