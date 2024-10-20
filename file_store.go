package cache

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileStore A in file key-value based cache store, implements [Store], this store
// should not be considered thread-safe by its own
type FileStore struct {
	cacheFilePath string
}

var _ Store = (*FileStore)(nil)

// NewFileStore Creates a new usable [FileStore]
func NewFileStore(dir string) *FileStore {
	return &FileStore{
		cacheFilePath: filepath.Join(dir, "file_store.kv"),
	}
}

func (p FileStore) openCacheFile() (*os.File, error) {
	return os.OpenFile(p.cacheFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
}

func (p *FileStore) Del(key any) {
	tempFile, _ := os.CreateTemp("", "file_store_*.kv")

	file, err := p.openCacheFile()
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")

		if line[0] == key {
			continue
		}

		fmt.Fprintln(tempFile, scanner.Text())
	}

	os.Rename(tempFile.Name(), p.cacheFilePath)
}

func (p *FileStore) Get(key any) (value any, exists bool) {
	file, err := p.openCacheFile()
	if err != nil {
		return nil, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		expiration := line[2]
		expiresAt, _ := time.Parse(time.RFC3339, expiration)

		if line[0] == key && expiresAt.Before(time.Now()) {
			p.Del(key)
			return nil, false
		}

		if line[0] == key {
			return line[1], true
		}
	}

	return nil, false
}

func (p *FileStore) Set(key any, value any, expiresAt time.Time) {
	file, err := p.openCacheFile()
	if err != nil {
		return
	}
	defer file.Close()

	line := fmt.Sprintf(`%s,%s,%s`, key, value, expiresAt.Format(time.RFC3339))
	fmt.Fprintln(file, line)
}
