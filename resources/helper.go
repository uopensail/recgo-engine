package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

/**
 * @brief Finds the directory with the latest timestamp that contains a SUCCESS file
 *
 * This function searches through subdirectories of the given directory path,
 * where each subdirectory name is expected to be a numeric timestamp (integer).
 * It returns the path of the subdirectory with the highest timestamp value
 * that also contains a file named "SUCCESS".
 *
 * @param dir The root directory path to search in
 * @return The path of the latest directory containing SUCCESS file, empty string if not found
 * @return Error if directory doesn't exist, can't be read, or no valid directory found
 *
 * @note Only directories with numeric names (parseable as int64) are considered
 * @note The SUCCESS file must exist directly in the timestamp directory
 *
 * @example
 * @code
 * latestDir, err := FindLatestSuccessDir("/data/jobs")
 * if err != nil {
 *     log.Printf("Error: %v", err)
 * } else {
 *     fmt.Printf("Latest SUCCESS directory: %s", latestDir)
 * }
 * @endcode
 */
func FindLatestSuccessDir(dir string) (string, error) {
	// Check if the root directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist: %s", dir)
	}

	// Read all entries in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %v", err)
	}

	var latestTimestamp int64 = -1
	var latestDir string

	// Iterate through all entries
	for _, entry := range entries {
		// Only process directories
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()

		// Try to parse directory name as timestamp (integer)
		timestamp, err := strconv.ParseInt(strings.TrimSpace(dirName), 10, 64)
		if err != nil {
			// Skip if directory name is not a valid integer
			continue
		}

		// Build the full path of the subdirectory
		subDirPath := filepath.Join(dir, dirName)

		// Check if SUCCESS file exists in this directory
		successFilePath := filepath.Join(subDirPath, "SUCCESS")
		if _, err := os.Stat(successFilePath); err != nil {
			// SUCCESS file doesn't exist, skip this directory
			continue
		}

		// Update record if this is the latest timestamp found so far
		if timestamp > latestTimestamp {
			latestTimestamp = timestamp
			latestDir = subDirPath
		}
	}

	// Return error if no valid directory was found
	if latestDir == "" {
		return "", fmt.Errorf("no timestamp directory with SUCCESS file found")
	}

	return latestDir, nil
}

type Finder struct {
	dir        string // Directory to monitor for timestamp subdirectories
	interval   int    // Check interval in seconds
	creator    func(string) (model.Resource, error)
	stopCh     chan struct{} // Channel to stop the monitoring goroutine
	isWatching atomic.Bool   // Flag to indicate if monitoring is active
	resource   atomic.Value
}

func NewFinder(dir string, creator func(string) (model.Resource, error)) (*Finder, error) {
	zlog.LOG.Info("Finder: creating new helper", zap.String("dir", dir))

	filePath, err := FindLatestSuccessDir(dir)
	if err != nil {
		zlog.LOG.Error("Finder: failed to find latest success directory during initialization",
			zap.String("dir", dir),
			zap.Error(err))
		return nil, err
	}

	resource, err := creator(dir)
	if err != nil {
		zlog.LOG.Error("Finder: failed to load resource during initialization",
			zap.String("file_path", filePath),
			zap.String("dir", dir),
			zap.Error(err))
		return nil, err
	}

	finder := &Finder{
		dir:      dir,
		stopCh:   make(chan struct{}), // Initialize the channel
		interval: interval,
		creator:  creator,
	}

	finder.resource.Store(resource)

	zlog.LOG.Info("Finder: successfully created helper",
		zap.String("dir", dir),
		zap.String("initial_file_path", filePath),
		zap.Int("interval_seconds", interval))

	finder.start()
	return finder, nil
}

func (f *Finder) Get() model.Resource {
	return f.resource.Load().(model.Resource)
}

func (f *Finder) start() {
	// Check if already watching using CAS
	if !f.isWatching.CompareAndSwap(false, true) {
		zlog.LOG.Warn("Finder: already watching, ignoring duplicate start request",
			zap.String("dir", f.dir))
		return
	}

	go f.watchLoop()
	zlog.LOG.Info("Finder: started watching directory",
		zap.String("dir", f.dir),
		zap.Int("interval_seconds", f.interval))
}

/**
 * @brief Stops the directory monitoring goroutine
 *
 * This method gracefully stops the background monitoring goroutine.
 * It's safe to call multiple times.
 *
 * @note This method is non-blocking
 */
func (f *Finder) Stop() {
	if f.isWatching.CompareAndSwap(true, false) {
		close(f.stopCh)
		zlog.LOG.Info("Finder: stopped watching directory",
			zap.String("dir", f.dir))
	} else {
		zlog.LOG.Debug("InvertedIndexHelper: stop called but not watching",
			zap.String("dir", f.dir))
	}
}

/**
 * @brief Internal method that runs the directory monitoring loop
 */
func (f *Finder) watchLoop() {
	zlog.LOG.Debug("InvertedIndexHelper: watch loop started", zap.String("dir", f.dir))

	ticker := time.NewTicker(time.Duration(f.interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			zlog.LOG.Debug("InvertedIndexHelper: periodic check triggered", zap.String("dir", f.dir))
			f.checkAndUpdate()
		case <-f.stopCh:
			zlog.LOG.Debug("InvertedIndexHelper: watch loop stopped", zap.String("dir", f.dir))
			return
		}
	}
}

/**
 * @brief Internal method to check for updates and reload resource if necessary
 */
func (f *Finder) checkAndUpdate() {
	// Find latest SUCCESS directory
	latestPath, err := FindLatestSuccessDir(f.dir)
	if err != nil {
		zlog.LOG.Error("Finder: failed to find latest success dir",
			zap.String("dir", f.dir),
			zap.Error(err))
		return
	}

	current := f.Get()
	if current != nil {
		url, err := current.GetURL()
		if err == nil && url == latestPath {
			zlog.LOG.Debug("Finder: resource already up to date",
				zap.String("dir", f.dir),
				zap.String("path", latestPath))
			return
		}
	}

	next, err := f.creator(latestPath)
	if err != nil {
		zlog.LOG.Error("Finder: failed to load resource from latest path",
			zap.String("dir", f.dir),
			zap.String("path", latestPath),
			zap.Error(err))
		return
	}

	// Atomically update
	oldPath := "none"
	if current != nil {
		oldPath, _ = current.GetURL()
	}

	f.resource.Store(next)
	zlog.LOG.Info("Find: successfully updated resource",
		zap.String("dir", f.dir),
		zap.String("new_path", latestPath),
		zap.String("old_path", oldPath))
}
