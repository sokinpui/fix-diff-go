package fixdiff

import (
	"fmt"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

func Fix(diffContent, originalContent string) (string, error) {
	parsedDiff, err := Parse(strings.NewReader(diffContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse diff: %w", err)
	}

	modifiedContent, err := applyPatches(originalContent, parsedDiff.Hunks)
	if err != nil {
		return "", fmt.Errorf("failed to apply patches: %w", err)
	}

	return generateCorrectDiff(originalContent, modifiedContent, parsedDiff.FromFile, parsedDiff.ToFile)
}

// applyPatches constructs the modified file content by applying the changes from each hunk
// to the original file content.
func applyPatches(originalContent string, hunks []*Hunk) (string, error) {
	originalLines := difflib.SplitLines(originalContent)
	var modifiedLines []string
	lastIndex := 0

	for _, hunk := range hunks {
		searchSnippet := buildSearchSnippet(hunk)
		if len(searchSnippet) == 0 {
			continue // Pure-addition hunks cannot be located this way.
		}

		startIndex, err := findSnippetIndex(originalLines, searchSnippet, lastIndex)
		if err != nil {
			return "", fmt.Errorf("could not locate hunk in original file: %w", err)
		}

		// Add the lines from the original content before the patch.
		modifiedLines = append(modifiedLines, originalLines[lastIndex:startIndex]...)

		// Add the "new" lines from the hunk.
		modifiedLines = append(modifiedLines, buildModifiedSnippet(hunk)...)

		// Update the index to skip over the "old" lines in the original content.
		lastIndex = startIndex + len(searchSnippet)
	}

	// Append any remaining content from the original file.
	if lastIndex < len(originalLines) {
		modifiedLines = append(modifiedLines, originalLines[lastIndex:]...)
	}

	return strings.Join(modifiedLines, ""), nil
}

// findSnippetIndex searches for a sequence of lines (snippet) within a larger set of lines,
// starting from a given index. It returns the starting line number of the found snippet.
func findSnippetIndex(lines, snippet []string, startIndex int) (int, error) {
	for i := startIndex; i <= len(lines)-len(snippet); i++ {
		match := true
		for j := 0; j < len(snippet); j++ {
			if lines[i+j] != snippet[j] {
				match = false
				break
			}
		}
		if match {
			return i, nil
		}
	}
	return -1, fmt.Errorf("snippet not found")
}

// buildSearchSnippet creates a slice of strings representing the original state of a hunk,
// containing only context and removed lines.
func buildSearchSnippet(hunk *Hunk) []string {
	var snippet []string
	for _, line := range hunk.Lines {
		if line.Type == ContextLine || line.Type == RemovedLine {
			// difflib.SplitLines keeps the newline characters, so we add them back for an accurate search.
			snippet = append(snippet, line.Content+"\n")
		}
	}
	return snippet
}

// buildModifiedSnippet creates a slice of strings representing the modified state of a hunk,
// containing only context and added lines.
func buildModifiedSnippet(hunk *Hunk) []string {
	var snippet []string
	for _, line := range hunk.Lines {
		if line.Type == ContextLine || line.Type == AddedLine {
			snippet = append(snippet, line.Content+"\n")
		}
	}
	return snippet
}

// generateCorrectDiff uses a diffing library to create a clean, unified diff
// between the original and modified content.
func generateCorrectDiff(originalContent, modifiedContent, fromFile, toFile string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(originalContent),
		B:        difflib.SplitLines(modifiedContent),
		FromFile: fromFile,
		ToFile:   toFile,
		Context:  3,
	}
	return difflib.GetUnifiedDiffString(diff)
}
