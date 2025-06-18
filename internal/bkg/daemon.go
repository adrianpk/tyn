package bkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	DaemonPIDFile = ".tyn/daemon.pid"
	DaemonLogFile = ".tyn/daemon.log"
)

func EnsureDaemon() error {
	running, err := IsDaemonRunning()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not check if daemon is running: %v\n", err)
	}

	if running {
		return nil
	}

	pid, err := StartDetachedDaemon()
	if err != nil {
		return fmt.Errorf("error starting daemon: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Started daemon process with PID %d\n", pid)
	return nil
}

func IsDaemonRunning() (bool, error) {
	pidFile, err := getPidFilePath()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return false, nil
	}

	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		return false, fmt.Errorf("error reading PID file: %w", err)
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return false, fmt.Errorf("invalid PID in file: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false, nil
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, nil
	}

	if executableMatch, err := isExecutableMatch(pid); err == nil && !executableMatch {
		return false, nil
	}

	return true, nil
}

func StartDetachedDaemon() (int, error) {
	exe, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("error getting executable path: %w", err)
	}

	logFile, err := getLogFilePath()
	if err != nil {
		return 0, err
	}

	logDir := filepath.Dir(logFile)
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		return 0, fmt.Errorf("error creating log directory: %w", err)
	}

	logFd, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("error opening log file: %w", err)
	}

	cmd := exec.Command(exe, "serve", "--daemon")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	cmd.Stdout = logFd
	cmd.Stderr = logFd

	err = cmd.Start()
	if err != nil {
		logFd.Close()
		return 0, fmt.Errorf("error starting daemon process: %w", err)
	}

	pid := cmd.Process.Pid

	err = writePidFile(pid)
	if err != nil {
		logFd.Close()
		return 0, fmt.Errorf("error writing PID file: %w", err)
	}

	err = cmd.Process.Release()
	if err != nil {
		logFd.Close()
		return 0, fmt.Errorf("error releasing daemon process: %w", err)
	}

	logFd.Close()
	return pid, nil
}

func getPidFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}
	return filepath.Join(home, DaemonPIDFile), nil
}

func getLogFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}
	return filepath.Join(home, DaemonLogFile), nil
}

func writePidFile(pid int) error {
	pidFile, err := getPidFilePath()
	if err != nil {
		return err
	}

	pidDir := filepath.Dir(pidFile)
	err = os.MkdirAll(pidDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating PID directory: %w", err)
	}

	err = os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("error writing PID to file: %w", err)
	}

	return nil
}

func isExecutableMatch(pid int) (bool, error) {
	procPath := fmt.Sprintf("/proc/%d/exe", pid)
	target, err := os.Readlink(procPath)
	if err != nil {
		return false, err
	}

	selfPath, err := os.Executable()
	if err != nil {
		return false, err
	}

	return target == selfPath, nil
}
