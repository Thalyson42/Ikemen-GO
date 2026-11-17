package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// memFromProc reads a single "Key:   value kB" line from a /proc file.
// Returns an empty string on platforms without /proc (Windows, macOS).
func memFromProc(file, key string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	for _, l := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(l, key+":") {
			f := strings.Fields(l)
			if len(f) >= 2 {
				return fmt.Sprintf("%s=%dMB", key, atoiSafe(f[1])/1024)
			}
		}
	}
	return ""
}

func atoiSafe(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// logMem appends one line with memory usage to save/logs/memlog.txt.
// It opens and closes the file on every call, so the line is already
// written even if Android kills the app right after (out of memory).
func logMem(tag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	dir := filepath.Join(sys.baseDir, "save", "logs")
	os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(filepath.Join(dir, "memlog.txt"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s | %s | goheap=%dMB | %s | %s | %s\n",
		time.Now().Format("15:04:05"), tag, m.HeapAlloc>>20,
		memFromProc("/proc/self/status", "VmRSS"),
		memFromProc("/proc/self/status", "VmHWM"),
		memFromProc("/proc/meminfo", "MemAvailable"))
}