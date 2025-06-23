package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DoneTaskListDays      int           `yaml:"done_task_list_days"`
	DatabaseBackupCount   int           `yaml:"database_backup_count"`
	NotificationTimeout   time.Duration `yaml:"notification_timeout"`
	JournalUpdateInterval time.Duration `yaml:"journal_update_interval"`
	PollInterval          time.Duration `yaml:"poll_interval"`
}

func DefaultConfig() Config {
	return Config{
		DoneTaskListDays:      14,
		DatabaseBackupCount:   2,
		NotificationTimeout:   5 * time.Second,
		JournalUpdateInterval: 1 * time.Minute,
		PollInterval:          30 * time.Second,
	}
}

func ensureConfigFile(path string, cfg Config) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Load() *Config {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "tyn", "tyn.yml")
		err := ensureConfigFile(configPath, cfg)
		if err == nil {
			if data, err := ioutil.ReadFile(configPath); err == nil {
				yaml.Unmarshal(data, &cfg)
			}
		}
	}
	doneTaskListDays := flag.Int("done-task-list-days", intVal("TYN_DONE_TASK_LIST_DAYS", cfg.DoneTaskListDays), "How many days of done tasks to show in lists (default: 14)")
	dbBackupCount := flag.Int("db-backup-count", intVal("TYN_DB_BACKUP_COUNT", cfg.DatabaseBackupCount), "How many automatic database backups to keep (default: 2)")
	notificationTimeout := flag.Duration("notification-timeout", durationVal("TYN_NOTIFICATION_TIMEOUT", cfg.NotificationTimeout), "Notification timeout (e.g. 5s, 10s)")
	journalUpdateInterval := flag.Duration("journal-update-interval", durationVal("TYN_JOURNAL_UPDATE_INTERVAL", cfg.JournalUpdateInterval), "How often to update the journal (e.g. 1m, 10m)")
	pollInterval := flag.Duration("poll-interval", durationVal("TYN_POLL_INTERVAL", cfg.PollInterval), "How often to poll for notifications and periodic tasks (e.g. 30s, 60s)")

	flag.Parse()

	cfg.DoneTaskListDays = *doneTaskListDays
	cfg.DatabaseBackupCount = *dbBackupCount
	cfg.NotificationTimeout = *notificationTimeout
	cfg.JournalUpdateInterval = *journalUpdateInterval
	cfg.PollInterval = *pollInterval

	return &cfg
}

func envVal(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func intVal(key string, fallback int) int {
	v := os.Getenv(key)
	if v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func durationVal(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}
