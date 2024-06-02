// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// GormLogger is a custom logger for Gorm.
type GormLogger struct {
	slowThreshold         time.Duration
	skipErrRecordNotFound bool
	showSQL               bool
	entry                 *logrus.Entry
}

// NewGormLogger creates a new Gorm logger.
func NewGormLogger(logger *logrus.Entry, slowThreshold time.Duration, skipNotFound, showSQL bool) *GormLogger {
	return &GormLogger{
		skipErrRecordNotFound: skipNotFound,
		slowThreshold:         slowThreshold,
		showSQL:               showSQL,
		entry:                 logger,
	}
}

// LogMode sets the log mode for the logger.
func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

// Info implements the logger.Interface.
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.entry.WithContext(ctx).Info(msg, args)
}

// Warn implements the logger.Interface.
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.entry.WithContext(ctx).Warn(msg, args)
}

// Error implements the logger.Interface.
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.entry.WithContext(ctx).Error(msg, args)
}

// Trace implements the logger.Interface.
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{
		"rows":    rows,
		"elapsed": elapsed,
		"source":  utils.FileWithLineNum(),
	}

	if l.showSQL {
		fields["sql"] = sql
	}

	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.skipErrRecordNotFound) {
		l.entry.WithContext(ctx).WithError(err).WithFields(fields).Error("gorm error")
		return
	}

	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		l.entry.WithContext(ctx).WithFields(fields).Warnf("gorm warn SLOW QUERY >= %s, took %s", l.slowThreshold, elapsed)
		return
	}

	l.entry.WithContext(ctx).WithFields(fields).Infof("gorm info")
}
