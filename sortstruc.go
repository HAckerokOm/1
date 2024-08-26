package main

import (
	"fmt"
	"sort"
)

type FileEntry struct {
	FType string
	FSize int64
	FName string
}

// Print - печать в консоли
func (fileEntry *FileEntry) Print() {
	newSize := formatSize(fileEntry.FSize)
	fmt.Printf("%-15s %-16s %-10s\n", fileEntry.FType, newSize, fileEntry.FName)
}

func formatSize(bytes int64) string {
	switch {
	case bytes >= 1000*1000*1000: // Если размер в гигабайтах или больше
		return fmt.Sprintf("%.2f Гигабайта", float64(bytes)/float64(1000*1000*1000))
	case bytes >= 1000*1000: // Если размер в мегабайтах
		return fmt.Sprintf("%.2f Мегабайта", float64(bytes)/float64(1000*1000))
	case bytes >= 1000: // Если размер в килобайтах
		return fmt.Sprintf("%.2f Килобайта", float64(bytes)/float64(1000))
	default: // Если размер меньше килобайта
		return fmt.Sprintf("%d Байта", bytes)
	}
}
func SortFileEntry(dataFiles []FileEntry, ask bool) {
	if !ask {
		sort.Slice(dataFiles, func(i, j int) bool {
			return dataFiles[i].FSize < dataFiles[j].FSize
		})
	} else {
		sort.Slice(dataFiles, func(i, j int) bool {
			return dataFiles[i].FSize > dataFiles[j].FSize
		})
	}
}
