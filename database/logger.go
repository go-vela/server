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

type gormLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	DriverName            string
}

func NewGormLogger(driver string) *gormLogger {
	return &gormLogger{
		SkipErrRecordNotFound: true,
		DriverName:            driver,
	}
}

func (l *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).WithField("database", l.DriverName).Info(s, args)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).WithField("database", l.DriverName).Warn(s, args)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).WithField("database", l.DriverName).Error(s, args)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{
		"database": l.DriverName,
		"rows":     rows,
		"elapsed":  elapsed,
		"sql":      sql,
	}

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound) {
		logrus.WithContext(ctx).WithError(err).WithFields(fields).Error("gorm error")
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logrus.WithContext(ctx).WithFields(fields).Warn("slow query")
		return
	}

	logrus.WithContext(ctx).WithFields(fields).Infof("gorm trace")
}
