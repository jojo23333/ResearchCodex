package ideas

import (
	"os"
	"path/filepath"
	"sort"
)

// ListIdeaDirs returns the folder names under a project that contain idea.md.
func ListIdeaDirs(projectDir string) ([]string, error) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(projectDir, e.Name(), "idea.md")); err == nil {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}

// LatestIdeaDir returns the lexicographically last idea directory name (slug).
func LatestIdeaDir(projectDir string) (string, bool, error) {
	dirs, err := ListIdeaDirs(projectDir)
	if err != nil {
		return "", false, err
	}
	if len(dirs) == 0 {
		return "", false, nil
	}
	return dirs[len(dirs)-1], true, nil
}
