# SLF4GO

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/github/v/tag/MariusSchmidt/slf4go?label=Version)](https://github.com/MariusSchmidt/slf4go/releases)
[![codecov](https://codecov.io/gh/IhrUsername/slf4go/branch/main/graph/badge.svg)](https://codecov.io/gh/IhrUsername/slf4go)

SLF4GO is a Simple Logging Facade for Go, inspired by [Java's SLF4J](https://github.com/qos-ch/slf4j). It provides a clean, flexible, and extensible logging abstraction that decouples your application code from specific logging implementations.

## Overview

SLF4GO offers:

- A simple, consistent logging API
- Multiple logging levels (Fatal, Panic, Error, Warn, Info, Debug, Trace)
- Structured logging with tags
- Ability to create loggers with default tags
- Format string support (similar to fmt.Printf)
- Separation of logging API from implementation

By using SLF4GO, you can write code that logs messages without being tied to a specific logging library. This allows you to change the underlying logging implementation without modifying your application code.

## Prerequisites

The following runtimes and tools are required to build and run this software:

| Requirements | Version  | Installation                    |
|--------------|----------|--------------------------------|
| Golang       | `>=1.23` | See [Golang](https://go.dev/dl/) |

## Installation

To add SLF4GO to your project, run:

```bash
go get github.com/MariusSchmidt/slf4go
```

## Usage

### Basic Usage

### Log Levels

SLF4GO supports the following log levels (in descending order of severity):

1. **Fatal** - Logs the message and terminates the program with Exit(1)
2. **Panic** - Logs the message and then calls panic()
3. **Error** - For errors that should definitely be noted
4. **Warn** - For non-critical events that deserve attention
5. **Info** - For general operational entries about what's going on
6. **Debug** - Usually only enabled during development, produces verbose logging
7. **Trace** - Even finer-grained informational events than Debug

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.