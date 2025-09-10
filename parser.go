package fixdiff

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type LineType int

const (
	ContextLine LineType = iota
	AddedLine
	RemovedLine
)

type Line struct {
	Content string
	Type    LineType
}

// Hunk represents a single "@@ ... @@" block in a diff.
// It contains the lines of code that are part of the change.
type Hunk struct {
	Lines []Line
}

type UnifiedDiff struct {
	FromFile string
	ToFile   string
	Hunks    []*Hunk
}

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
		case currentHunk != nil && len(line) == 0:
			currentHunk.Lines = append(currentHunk.Lines, Line{Content: "", Type: ContextLine})
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
