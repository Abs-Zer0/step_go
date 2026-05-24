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
	var filesList []string
	var err error
	if printFiles {
		filesList, err = getDirsListWithFiles(path)
	} else {
		filesList, err = getDirsList(path)
	}
	if err != nil {
		return err
	}

	// Maybe sort

	treeList := convertList(filesList, "")
	for _, treeLine := range treeList {
		fmt.Fprint(out, treeLine+"\n")
	}

	return nil
}

func getDirsList(rootPath string) ([]string, error) {
	dirsList := []string{}
	err := filepath.WalkDir(rootPath, func(currentPath string, d os.DirEntry, err error) error {
		return appendPathTrimmed(&dirsList, rootPath, currentPath)
	})

	return dirsList, err
}

func getDirsListWithFiles(rootPath string) ([]string, error) {
	filesList := []string{}
	err := filepath.Walk(rootPath, func(currentPath string, info os.FileInfo, err error) error {
		return appendPathTrimmed(&filesList, rootPath, currentPath)
	})

	return filesList, err
}

func appendPathTrimmed(filesList *[]string, rootPath, currentPath string) error {
	if rootPath == currentPath {
		return nil
	}

	*filesList = append(*filesList, strings.TrimPrefix(currentPath, rootPath+string(os.PathSeparator)))

	return nil
}

func convertList(filesList []string, prefix string) (treeList []string) {
	for i := 0; i < len(filesList); {
		pathPrefix := filesList[i]
		lastPathIndex, _ := sort.Find(len(filesList), func(idx int) int {
			path := filesList[idx]
			pathPref := strings.SplitN(path, string(os.PathSeparator), 2)[0]

			return strings.Compare(pathPrefix, pathPref)
		})

		isLast := lastPathIndex == len(filesList)-1
		var nextPrefix string
		if isLast {
			treeList = append(treeList, prefix+"└───"+pathPrefix)
			nextPrefix = prefix + "\t"
			if i < lastPathIndex {
				treeList = append(treeList, convertList(filesList[i+1:lastPathIndex], prefix+"\t")...)
			}
		} else {
			treeList = append(treeList, prefix+"├───"+pathPrefix)
			nextPrefix = prefix + "│\t"
			if i < lastPathIndex {
				treeList = append(treeList, convertList(filesList[i+1:lastPathIndex], prefix+"│\t")...)
			}
		}

		if i < lastPathIndex {
			treeList = append(treeList, convertList(removePrefix(filesList[i+1:lastPathIndex], pathPrefix), nextPrefix)...)
		}

		i = lastPathIndex + 1
	}

	return
}

func removePrefix(filesList []string, pathPrefix string) []string {
	result := make([]string, 0, len(filesList))
	for _, path := range filesList {
		result = append(result, strings.TrimPrefix(path, pathPrefix+string(os.PathSeparator)))
	}

	return result
}
