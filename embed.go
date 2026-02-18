package nebula

import "embed"

// COMPILER DIRECTIVE
//
//go:embed all:terraform/* all:ansible/*
var ProjectFiles embed.FS
