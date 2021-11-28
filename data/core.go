package data

import (
	"bufio"
	"fmt"
	"os"
)

func WriteOnePair(pair string) error {
	err := writeOnePair(pair)
	return err
}

func RemoveOnePair(pair string) error {
	err := removeOnePair(pair)
	return err
}

func writeOnePair(pair string) error {
	path := absolutePath() + "/pairs.txt"
	pairs, _ := readLines(path)
	pairs = append(pairs, pair)
	err := writeLines(pairs, path)
	return err
}

func removeOnePair(pair string) error {
	path := absolutePath() + "/pairs.txt"
	pairs, _ := readLines(path)
	_pairs := []string{}
	for _, v := range pairs {
		if v != pair {
			_pairs = append(_pairs, v)
		}
	}
	fmt.Println(pairs)
	fmt.Println(_pairs)
	err := writeLines(_pairs, path)
	return err
}

func absolutePath() string {
	path, _ := os.UserCacheDir()
	return path
}

func readLines(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		writeLines(lines, path)
		return lines, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func readBytes(path string) ([]byte, error) {
	var bytes []byte
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = file.Read(bytes)

	return bytes, err
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func writeBytes(bytes []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	return err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isExistedPairs(p string, pairs []string) bool {
	for _, v := range pairs {
		if v == p {
			return true
		}
	}
	return false
}
