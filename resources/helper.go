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
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

const (
	// FinderCheckInterval defines how often (in seconds) the Finder will check for updated resources.
	FinderCheckInterval = 300 // 5 minutes
)

/**
 * FindLatestSuccessDir finds the directory with the latest numeric timestamp that contains a SUCCESS file.
 *
 * This function scans all subdirectories under `dir`, where each subdirectory name must be a numeric timestamp.
 * Only directories containing a file named "SUCCESS" are considered valid.
 *
 * @param dir Root directory to scan.
 * @return Path to the latest valid directory, or error if none found.
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

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()

		// Parse directory name as timestamp (int64)
		timestamp, err := strconv.ParseInt(strings.TrimSpace(dirName), 10, 64)
		if err != nil {
			continue
		}

		subDirPath := filepath.Join(dir, dirName)
		successFilePath := filepath.Join(subDirPath, "SUCCESS")

		// Check if SUCCESS file exists
		if _, err := os.Stat(successFilePath); err != nil {
			continue
		}

		// Update if this timestamp is newer
		if timestamp > latestTimestamp {
			latestTimestamp = timestamp
			latestDir = subDirPath
		}
	}

	if latestDir == "" {
		return "", fmt.Errorf("no timestamp directory with SUCCESS file found in %s", dir)
	}

	return latestDir, nil
}

// Finder monitors a directory for the latest timestamp-based resources.
type Finder struct {
	dir        string // Root directory containing timestamp subdirectories
	creator    func(string) (model.Resource, error)
	stopCh     chan struct{}
	isWatching atomic.Bool
	resource   atomic.Value
}

// NewFinder creates a new Finder, initializes it with the latest resource, and starts watching.
//
// @param dir Directory to monitor.
// @param creator Function to load a resource from a given path.
// @return Finder instance or error if initialization fails.
func NewFinder(dir string, creator func(string) (model.Resource, error)) (*Finder, error) {
	pStat := prome.NewStat("NewFinder")
	defer pStat.End()

	zlog.LOG.Info("Finder: creating new instance", zap.String("dir", dir))

	latestDir, err := FindLatestSuccessDir(dir)
	if err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Finder: failed to find latest SUCCESS directory",
			zap.String("dir", dir),
			zap.Error(err))
		return nil, err
	}

	// Load initial resource using latest SUCCESS path
	resource, err := creator(latestDir)
	if err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Finder: failed to load initial resource",
			zap.String("latest_dir", latestDir),
			zap.Error(err))
		return nil, err
	}

	f := &Finder{
		dir:     dir,
		creator: creator,
		stopCh:  make(chan struct{}),
	}

	f.resource.Store(resource)

	zlog.LOG.Info("Finder: initialized successfully",
		zap.String("dir", dir),
		zap.String("initial_latest_dir", latestDir),
		zap.Int("interval_seconds", FinderCheckInterval))

	f.start()
	return f, nil
}

// Get returns the current loaded resource.
func (f *Finder) Get() model.Resource {
	return f.resource.Load().(model.Resource)
}

// start launches the watch loop if not already running.
func (f *Finder) start() {
	if !f.isWatching.CompareAndSwap(false, true) {
		zlog.LOG.Warn("Finder: already watching, ignoring duplicate start",
			zap.String("dir", f.dir))
		return
	}

	go f.watchLoop()
	zlog.LOG.Info("Finder: started watching directory",
		zap.String("dir", f.dir),
		zap.Int("interval_seconds", FinderCheckInterval))
}

// Stop gracefully stops the monitoring goroutine.
// Safe to call multiple times.
func (f *Finder) Stop() {
	if f.isWatching.CompareAndSwap(true, false) {
		close(f.stopCh)
		zlog.LOG.Info("Finder: stopped watching directory", zap.String("dir", f.dir))
	} else {
		zlog.LOG.Debug("Finder: stop called but not watching", zap.String("dir", f.dir))
	}
}

// watchLoop periodically checks for updated resources.
func (f *Finder) watchLoop() {
	zlog.LOG.Debug("Finder: watch loop started", zap.String("dir", f.dir))
	ticker := time.NewTicker(time.Duration(FinderCheckInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			zlog.LOG.Debug("Finder: periodic check triggered", zap.String("dir", f.dir))
			f.checkAndUpdate()
		case <-f.stopCh:
			zlog.LOG.Debug("Finder: watch loop stopped", zap.String("dir", f.dir))
			return
		}
	}
}

// checkAndUpdate reloads the resource if a new latest SUCCESS directory appears.
func (f *Finder) checkAndUpdate() {
	latestPath, err := FindLatestSuccessDir(f.dir)
	if err != nil {
		zlog.LOG.Error("Finder: failed to find latest SUCCESS dir",
			zap.String("dir", f.dir),
			zap.Error(err))
		return
	}

	current := f.Get()
	if current != nil {
		url := current.GetURL()
		if url == latestPath {
			zlog.LOG.Debug("Finder: resource up-to-date",
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

	oldPath := "none"
	if current != nil {
		oldPath = current.GetURL()
	}

	f.resource.Store(next)
	zlog.LOG.Info("Finder: successfully updated resource",
		zap.String("dir", f.dir),
		zap.String("new_path", latestPath),
		zap.String("old_path", oldPath))
}
