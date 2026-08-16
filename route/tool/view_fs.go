package tool

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	builtin_views "opennamu/views"
)

func Clean_view_path(data string) string {
	data = strings.ReplaceAll(data, "\\", "/")
	data = strings.TrimPrefix(data, "/")
	data = path.Clean(data)

	if data == "." || data == ".." || strings.HasPrefix(data, "../") {
		return ""
	}

	return data
}

func Read_view_file(data string) ([]byte, error) {
	data = Clean_view_path(data)
	if data == "" || strings.HasSuffix(data, ".go") {
		return nil, fs.ErrNotExist
	}

	raw, err := builtin_views.Read(data)
	if err == nil {
		return raw, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	for _, base := range []string{"views"} {
		raw, err := os.ReadFile(filepath.Join(base, filepath.FromSlash(data)))
		if err == nil {
			return raw, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	return nil, fs.ErrNotExist
}

func List_view_dir(data string) ([]fs.DirEntry, error) {
	data = Clean_view_path(data)

	entry_map := map[string]fs.DirEntry{}

	builtin_path := data
	if builtin_path == "" {
		builtin_path = "."
	}

	entries, err := builtin_views.ReadDir(builtin_path)
	if err == nil {
		for _, entry := range entries {
			entry_map[entry.Name()] = entry
		}
	}

	for _, base := range []string{"views"} {
		entries, external_err := os.ReadDir(filepath.Join(base, filepath.FromSlash(data)))
		if external_err == nil {
			for _, entry := range entries {
				if _, ok := entry_map[entry.Name()]; !ok {
					entry_map[entry.Name()] = entry
				}
			}
		}
	}

	if len(entry_map) == 0 && err != nil {
		return nil, err
	}

	entry_list := []fs.DirEntry{}
	for _, entry := range entry_map {
		entry_list = append(entry_list, entry)
	}
	sort.Slice(entry_list, func(a int, b int) bool {
		return entry_list[a].Name() < entry_list[b].Name()
	})

	return entry_list, nil
}
