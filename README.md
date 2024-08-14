# Hardware Monitor

## Overview

The Hardware Monitor project is a lightweight Go-based system that provides real-time monitoring of a computer's hardware metrics, including CPU, memory, and disk usage. The information is displayed on a web interface that is updated dynamically using WebSockets.

![Application Screenshot](https://imgur.com/JOOrs9e.png)

## Features

- **System Information Monitoring**: Captures and displays system information such as hostname, total and used memory, operating system, disk space, and CPU details.
- **Real-time Updates**: Utilizes WebSockets to push hardware metrics to the frontend in real-time.
- **Responsive UI**: A user-friendly and responsive web interface designed with modern CSS practices.
- **Templating**: Custom Go templates are used to render hardware data into HTML.

## Project Structure

```
hardware-monitor/
├── bin/                  # Compiled binaries (optional, .gitignored)
├── cmd/                  
│   └── main.go           # Entry point for the application
├── internal/             
│   ├── hardware/         
│   │   └── hardware.go   # Hardware information retrieval
│   ├── server/           
│   │   ├── metrics.go    # Metrics struct definition
│   │   ├── server.go     # Web server implementation
│   │   ├── server_test.go # Unit tests for server components
│   │   └── subscriber.go # WebSocket subscriber handling
│   └── templates/        
│       ├── monitor.templ # Go template for rendering hardware data
│       └── monitor_templ.go # Template integration with Go code
├── static/               
│   ├── assets/           # Static assets such as images, icons
│   ├── index.html        # Main HTML file served by the server
│   └── styles.css        # CSS styling for the UI
├── .gitignore            # Ignored files and directories
├── Makefile              # Automation tasks (build, run, etc.)
├── TODO.md               # Task list and future improvements
├── go.mod                # Go module dependencies
└── go.sum                # Dependency checksums
```

## Getting Started

### Prerequisites

- Go 1.16+ installed
- Git for version control
- Access to the terminal/command prompt

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/hardware-monitor-go.git
   cd hardware-monitor-go
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the project:

   ```bash
   make build
   ```

4. Run the project:

   ```bash
   make run
   ```

   The server will start on `http://localhost:8080`.

### Usage

- **Web Interface**: Open your web browser and navigate to `http://localhost:8080`. The page will display real-time hardware metrics.
- **Customizing the UI**: Modify `static/styles.css` to customize the look and feel of the interface.
- **Extend the Metrics**: Add additional metrics by editing `internal/hardware/hardware.go` and corresponding Go template files.

## Testing

Run the tests with:

```bash
make test
```

The tests cover the server components and ensure the correct functionality of WebSocket communication and hardware data retrieval.

## TODO

- Add support for monitoring GPU usage.
- Implement user authentication for accessing the monitoring page.
- Enhance error handling and logging.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

## Acknowledgements

- [gopsutil](https://github.com/shirou/gopsutil) for hardware metrics retrieval.
- [htmx](https://htmx.org) for modern HTML extensions.

---

Happy Monitoring!

