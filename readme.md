# Linux Bootstrap

A cross-distro Linux post-installation automation tool written in Go.

The goal of **Linux Bootstrap** is to automate the common setup tasks performed after installing a fresh Linux distribution, while detecting the system and adapting to its environment.

## Current Status

🚧 **Early development**

Currently supported:

* [x] Linux distribution detection
* [x] Distribution version detection
* [x] Package manager detection
* [x] Package manager abstraction
* [x] Initial APT implementation
* [x] Safe command execution
* [x] Dry-run mode
* [x] User confirmation before system changes
* [x] APT support
* [x] DNF support
* [x] Pacman support
* [x] Multimedia codec installation
* [x] NVIDIA driver detection and setup

Currently tested on:

* Pop!_OS 24.04

## Planned Features

* [ ] Development environment setup
* [ ] Docker setup
* [ ] System diagnostics
* [ ] Post-installation verification
* [ ] Additional Linux distributions

## Project Structure

```text
linux-bootstrap/
├── main.go
├── go.mod
└── internal/
    ├── system/
    │   └── os.go
    └── packages/
        ├── manager.go
        ├── package_manager.go
        ├── apt.go
        └── factory.go
```

## Development

Clone the repository and run:

```bash
go run .
```

The project is currently under active development. Features and supported distributions will be added incrementally.

## License

License to be decided.
