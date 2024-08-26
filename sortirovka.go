package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

func calcSumDirect(pathDirectory string) (int64, error) {

	//считывание содержания директории
	files, err := os.ReadDir(pathDirectory)
	if err != nil {
		return 0, err
	}
	//суммарный размер директории
	var sum int64
	//проход по кадому файлу
	for _, file := range files {
		//формирование нового пути к внутренней директории и рекурсивный вызов
		//с этим путем в качестве аргумента
		if file.IsDir() {
			dirSum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name()))
			if err != nil { // Обработка ошибок при получении информации о файле
				fmt.Printf("Не удалось получить размер директории'%s': %v\n", file.Name(), err)
				continue
			}
			sum += dirSum
			continue
		}
		info, _ := file.Info()
		sum += info.Size()
	}
	return sum, nil
}

// Функция parseFlags анализирует командные флаги для получения пути к директории.
func parseFlags() (string, bool) {
	sortFlag := flag.Bool("sort", false, "Сортировка файлов по размеру")
	includeCurrentDirPtr := flag.String("dst", "", "Путь к директории")
	flag.Parse()
	if *includeCurrentDirPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	return *includeCurrentDirPtr, *sortFlag
}

func GetDir(file os.DirEntry, pathDirectory string, filesData *[]FileEntry, wg *sync.WaitGroup) {
	var size int64
	fileInfo, _ := file.Info()
	defer wg.Done()
	directsum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name()))
	size += directsum //суммируем размеры
	if err != nil {   // Обработка ошибок при получении информации о файле
		fmt.Printf("Не удалось получить размер дирректории '%s': %v\n", file.Name(), err)
		return
	}
	size += fileInfo.Size()
	//сохранение данных о директории

	*filesData = append(*filesData, FileEntry{"Директория", size, file.Name()})

}

// Функция printFileDetails выводит детальную информацию о каждом файле в переданном списке объектов os.DirEntry.
func printFileDetails(mdir string) []FileEntry {
	files, _ := os.ReadDir(mdir)
		
	var wg sync.WaitGroup
	var fileisentry = []FileEntry{}
	fmt.Printf("%-15s %-15s %-30s\n", "Тип", "Размер", "Название") // Вывод заголовков столбцов
	for _, file := range files {                                   // Итерация по каждому файлу
		info, err := file.Info() // Получение информации о файле
		if err != nil {          // Обработка ошибок при получении информации о файле
			fmt.Printf("Не удалось получить информацию о файле '%s': %v\n", file.Name(), err)
			continue
		}
		isDirectory := info.IsDir() // Определение, является ли запись директорией
		// size := info.Size()         // Получение размера файла/директории
		// name := file.Name()         // Получение имени файла/директории
		// tip := " "                  // Тип
		if isDirectory {
			wg.Add(1)
			go GetDir(file, mdir, &fileisentry, &wg)
		} else {
			fileisentry= append(fileisentry, FileEntry{"Файл", info.Size(), file.Name()})
		}
		// formattedSize := formatSize(size)                                    // Форматирование размера для вывода
		// line := fmt.Sprintf("%-15s %-15s %-30s\n", tip, formattedSize, name) // Форматирование строки для вывода
		// fmt.Println(line)                                                    // Вывод строки
	}
	wg.Wait()
	return fileisentry
}

func main() {
	dirPath, shouldSort := parseFlags()

	startTime := time.Now()
	fmt.Println("Программа выполняется...")

	
	entryfiles := printFileDetails(dirPath)
	SortFileEntry(entryfiles,shouldSort)
	for _, f := range entryfiles{
		f.Print()

	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Время выполнения программы: %v\n", duration)
}
