# Infergo - A Simple Rule Engine in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/tarcisio/infergo)](https://goreportcard.com/report/github.com/yourusername/infergo)
[![Build Status](https://github.com/tarcisio/infergo/actions/workflows/go-ci.yml/badge.svg)](https://github.com/tarcisio/infergo/actions)

Infergo is a generic rule engine written in Go.
It allows you to define rules with conditions (when) and actions (then) to evaluate and
execute on a given payload. This engine is highly flexible, easy to extend, and useful for
applications that need dynamic rule evaluation.

## Features
- **Generic Rule Engine**: Supports any payload type using Go's generics.
- **Prioritized Rules**: Rules are executed based on a configurable priority system.
- **Cycle Limiting**: Avoid infinite loops with a configurable maximum cycle limit.
- **Simple API**: Define rules with conditions and actions, and execute them with minimal setup.

