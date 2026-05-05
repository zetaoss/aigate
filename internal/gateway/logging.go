package gateway

import (
	"log"
	"strings"
)

type logLevel int

const (
	logLevelDebug logLevel = iota
	logLevelInfo
	logLevelError
)

func parseLogLevel(value string) logLevel {
	switch normalizeLogLevel(value) {
	case "debug":
		return logLevelDebug
	case "error":
		return logLevelError
	default:
		return logLevelInfo
	}
}

func normalizeLogLevel(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (g *Gateway) shouldLog(level logLevel) bool {
	current := parseLogLevel(g.cfg.Server.LogLevel)
	return level >= current
}

func (g *Gateway) logDebug(format string, args ...any) {
	if g.shouldLog(logLevelDebug) {
		log.Printf("[debug] "+format, args...)
	}
}

func (g *Gateway) logInfo(format string, args ...any) {
	if g.shouldLog(logLevelInfo) {
		log.Printf("[info] "+format, args...)
	}
}

func (g *Gateway) logError(format string, args ...any) {
	if g.shouldLog(logLevelError) {
		log.Printf("[error] "+format, args...)
	}
}
