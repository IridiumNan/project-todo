// Package cache is defined for cache data dir path on global cache file
// When init command is exec, it will add it to the cache path file
//
// If data dir not found for current search dir
// It will use fzf command for user select a data dir doing current command
package cache

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"

	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/utils"
)

const dirCacheFileName = "data_dir.txt"

func XDGCacheDir() string {
	home, err := models.GetHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(home, ".cache", models.AppName)
}

var DefaultDirCacheFilePath = path.Join(XDGCacheDir(), dirCacheFileName)

// DirCache manage the cache file which store all path that has be init
// [DirCache.SelectWithFzf] return a dir by fzf command selected
// [DirCache.PushNewDir] add a new dir into cache file
// [DirCache.Dirs] return all dirs that has been cached
type DirCache struct {
	cacheFilePath string
}

// SelectWithFzf read the cache file then use fzf to select a cache dir then return this dir
func (c *DirCache) SelectWithFzf() (dataDirPath string) {
	byteData, err := c.byteData()
	if err != nil || len(byteData) == 0 {
		slog.Warn("error when read byte data from cache file", "file_path", c.cacheFilePath, "err", err)
		return models.EmptyStr
	}
	buf := bytes.NewBuffer(byteData)

	fzfCmd := exec.Command("fzf")

	var selectedBuf *bytes.Buffer

	fzfCmd.Stdin = buf
	fzfCmd.Stdout = selectedBuf

	return selectedBuf.String()
}

// PushNewDir add a new data dir path into cache file
func (c *DirCache) PushNewDir(dataDirPath string) error {
	oldDirs, err := c.Dirs()
	if err != nil {
		return fmt.Errorf("error when get old dirs, err: %s", err.Error())
	}

	if slices.Contains(oldDirs, dataDirPath) {
		slog.Warn("data dir has already exist, ignore it", "dir_path", dataDirPath)
		return nil
	}

	// use append because the cache file will be trunc when call updateCacheFile function
	newDirs := append(oldDirs, dataDirPath)

	err = c.updateCacheFile(newDirs)
	if err != nil {
		return fmt.Errorf("error when update cache file with new dirs, path: %s, err: %s", c.cacheFilePath, err.Error())
	}

	return nil
}

// updateCacheFile update the cache file with newDirs
// It will trunc this cache file, remember to add old dirs on newDirs
func (c *DirCache) updateCacheFile(newDirs []string) error {
	strData := strings.Builder{}

	cleanDirs := c.delNotExistDirs(newDirs)

	for idx := range cleanDirs {
		str := cleanDirs[idx] + "\n"
		strData.WriteString(str)
	}

	slog.Info("checking strData", "strData", strData.String())

	cacheFile, err := os.OpenFile(c.cacheFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("error when opening cache file, file path: %s, err: %s", c.cacheFilePath, err.Error())
	}
	_, err = cacheFile.WriteString(strData.String())
	if err != nil {
		return fmt.Errorf("error when write cache dir into cache file, file path: %s, err: %s", c.cacheFilePath, err.Error())
	}

	err = cacheFile.Sync()
	if err != nil {
		return fmt.Errorf("error when sync cache file into disk, file path: %s, err: %s", c.cacheFilePath, err.Error())
	}

	return cacheFile.Close()
}

// delNotExistDirs return the clean dirs which not contains any not exist dir or not readable dir
func (c *DirCache) delNotExistDirs(dirs []string) []string {
	cleanDirs := []string{}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			cleanDirs = append(cleanDirs, dir)
		}
	}
	slog.Info("checking clean dirs", "dirs", cleanDirs)
	return cleanDirs
}

// byteData read data on [DirCache.CacheFilePath] then return then raw byteData
func (c *DirCache) byteData() ([]byte, error) {
	return os.ReadFile(c.cacheFilePath)
}

// Dirs return all dirs that has been cached on file
func (c *DirCache) Dirs() ([]string, error) {
	data, err := c.byteData()
	if err != nil {
		return nil, fmt.Errorf("cache: error when read byte data from file, file_path: %s, err: %s", c.cacheFilePath, err.Error())
	}

	dataStr := string(data)

	dirs := strings.Split(dataStr, "\n")
	cleanDirs := []string{}

	for idx := range dirs {
		if dirs[idx] != models.EmptyStr {
			cleanDirs = append(cleanDirs, dirs[idx])
		}
	}

	// slog.Info("Dirs: load dir", "cleanDirs", cleanDirs)

	return cleanDirs, nil
}

// NewCache create a [Cache]
func NewCache(cacheFilePath string) (*DirCache, error) {
	err := utils.EnsureFileExist(cacheFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when create cache file, cache file path: %s, err: %s", cacheFilePath, err.Error())
	}

	return &DirCache{
		cacheFilePath: cacheFilePath,
	}, nil
}

// DefaultCache return the DirCache with cacheFilePath = [DefaultDirCacheFilePath]
func DefaultCache() (*DirCache, error) {
	return NewCache(DefaultDirCacheFilePath)
}
