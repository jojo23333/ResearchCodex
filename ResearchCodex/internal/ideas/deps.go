package ideas

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

// DependencyEntry mirrors one line in .rcodex/idea_deps.jsonl.
type DependencyEntry struct {
	Project   string  `json:"project"`
	IdeaPath  string  `json:"idea_path"`
	DependsOn *string `json:"depends_on"`
	CreatedAt string  `json:"created_at"`
	Note      string  `json:"note,omitempty"`
}

func AppendDependency(path string, entry DependencyEntry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(enc, '\n')); err != nil {
		return err
	}
	return nil
}

func LoadDependencies(path string) ([]DependencyEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []DependencyEntry{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var entries []DependencyEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry DependencyEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// ResolveChain returns dependency entries from base to target idea path.
func ResolveChain(entries []DependencyEntry, targetPath string) []DependencyEntry {
	index := make(map[string]DependencyEntry, len(entries))
	for _, e := range entries {
		index[e.IdeaPath] = e
	}

	var chain []DependencyEntry
	current := targetPath
	for {
		entry, ok := index[current]
		if !ok {
			break
		}
		chain = append(chain, entry)
		if entry.DependsOn == nil || *entry.DependsOn == "" {
			break
		}
		current = *entry.DependsOn
	}

	// reverse to chronological order
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain
}
