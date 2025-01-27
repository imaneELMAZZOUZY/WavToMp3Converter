package app

import "errors"

var ErrFFMPEGNotFound = errors.New("ffmpeg.exe not found")

var ErrSqliteNotFound = errors.New("sqlite executable not found")