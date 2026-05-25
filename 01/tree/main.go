package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	filesList, err := getFilesList(path, printFiles)
	if err != nil {
		return err
	}

	sort.Strings(filesList)

	treeList := convertList(filesList, "")
	for _, treeLine := range treeList {
		fmt.Fprint(out, treeLine+"\n")
	}

	return nil
}

func getFilesList(rootPath string, printFiles bool) ([]string, error) {
	filesList := []string{}
	err := filepath.WalkDir(rootPath, func(currentPath string, d os.DirEntry, err error) error {
		if currentPath == rootPath {
			return nil
		}

		if !d.IsDir() && !printFiles {
			return nil
		}

		filesList = append(filesList, strings.TrimPrefix(currentPath, rootPath+string(os.PathSeparator)))

		return nil
	})

	return filesList, err
}

func convertList(filesList []string, prefix string) (treeList []string) {
	for i := 0; i < len(filesList); {
		path := filesList[i]
		lastEntryIndex := lastEntryIndex(filesList, path, i)
		isPathLast := lastEntryIndex == len(filesList)-1
		fmt.Println("Path:", path, "; Last entry index:", lastEntryIndex, "; Path is last:", isPathLast)

		var nextPrefix string
		if isPathLast {
			treeList = append(treeList, prefix+"└───"+path)
			nextPrefix = prefix + "\t"
		} else {
			treeList = append(treeList, prefix+"├───"+path)
			nextPrefix = prefix + "│\t"
		}

		if i < lastEntryIndex {
			subFilesList := filesList[i+1 : lastEntryIndex+1]
			removePrefix(subFilesList, path)
			treeList = append(treeList, convertList(subFilesList, nextPrefix)...)
		}

		i = lastEntryIndex + 1
	}

	return
}

func lastEntryIndex(filesList []string, root string, startIndex int) int {
	var i int
	for i = startIndex; i < len(filesList); i++ {
		if !strings.HasPrefix(filesList[i], root+string(os.PathSeparator)) {
			return i - 1
		}
	}

	return i - 1
}

func removePrefix(filesList []string, removePrefix string) {
	for i := range filesList {
		filesList[i] = strings.TrimPrefix(filesList[i], removePrefix+string(os.PathSeparator))
	}
}
