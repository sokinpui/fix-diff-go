package fixdiff

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// LineType defines the type of a line in a diff hunk.
type LineType int

const (
	// ContextLine represents a line that is part of the context.
	ContextLine LineType = iota
	// AddedLine represents a line that was added.
	AddedLine
	// RemovedLine represents a line that was removed.
	RemovedLine
)

// Line represents a single line within a diff hunk.
// It contains the content of the line and its type.
type Line struct {
	Content string
	Type    LineType
}

// Hunk represents a single "@@ ... @@" block in a diff.
// It contains the lines of code that are part of the change.
type Hunk struct {
	Lines []Line
}

// UnifiedDiff represents a parsed unified diff file.
// It holds the original and modified file paths and a slice of hunks.
type UnifiedDiff struct {
	FromFile string
	ToFile   string
	Hunks    []*Hunk
}

// Parse reads a unified diff from a reader and converts it into a UnifiedDiff struct.
// It intentionally ignores the line numbers in the hunk headers (`@@ -a,b +c,d @@`)
// as they are assumed to be potentially incorrect.
func Parse(reader io.Reader) (*UnifiedDiff, error) {
	scanner := bufio.NewScanner(reader)
	diff := &UnifiedDiff{}
	var currentHunk *Hunk

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "--- "):
			diff.FromFile = strings.TrimSpace(strings.TrimPrefix(line, "--- "))
		case strings.HasPrefix(line, "+++ "):
			diff.ToFile = strings.TrimSpace(strings.TrimPrefix(line, "+++ "))
		case strings.HasPrefix(line, "@@ "):
			currentHunk = &Hunk{}
			diff.Hunks = append(diff.Hunks, currentHunk)
		case strings.HasPrefix(line, "+") && currentHunk != nil:
			currentHunk.Lines = append(currentHunk.Lines, Line{Content: line[1:], Type: AddedLine})
		case strings.HasPrefix(line, "-") && currentHunk != nil:
			currentHunk.Lines = append(currentHunk.Lines, Line{Content: line[1:], Type: RemovedLine})
		case strings.HasPrefix(line, " ") && currentHunk != nil:
			currentHunk.Lines = append(currentHunk.Lines, Line{Content: line[1:], Type: ContextLine})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading diff content: %w", err)
	}

	if diff.FromFile == "" && diff.ToFile == "" && len(diff.Hunks) == 0 {
		return nil, fmt.Errorf("input does not appear to be a valid diff")
	}

	return diff, nil
}
