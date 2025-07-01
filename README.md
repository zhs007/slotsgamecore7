# slotsgamecore7

A modular, extensible slot game core framework written in Go.  
Supports rapid development, simulation, and analysis of slot games with modern math toolsets and flexible configuration.

## Features

- **Core Slot Game Engine**: Modular architecture for reels, symbols, paytables, player state, and more.
- **Math Toolset**: Tools for reel generation, RTP/statistics analysis, and configuration validation.
- **Low-code Support**: Quickly build and test slot games with low-code modules.
- **Simulation & Analytics**: Built-in simulation, statistics, and result analysis for game balancing.
- **Flexible Configuration**: Supports JSON, YAML, and XLSX for game and math configs.
- **Extensible Protocols**: gRPC, HTTP, and WebSocket server support for integration and deployment.
- **Relax/Third-party Engine Compatibility**: Generate configs for Relax and other engines.
- **GATI/Plugin System**: Easy integration with GATI and plugin-based extensions.

## Getting Started

1. **Clone the repository**
   ```sh
   git clone https://github.com/zhs007/slotsgamecore7.git
   cd slotsgamecore7
   ```

2. **Install dependencies**
   ```sh
   go mod tidy
   ```

3. **Build and run**
   ```sh
   go build ./...
   # or run a specific tool/server
   ./app/runsimserv.sh
   ```

4. **Run tests**
   ```sh
   go test ./...
   ```

## Directory Structure

- `game/`         — Core slot game logic and math modules
- `mathtoolset2/` — Advanced math tools for reel/statistics/simulation
- `asciigame/`    — ASCII-based slot game demo
- `app/`          — Scripts and entrypoints for various tools/servers
- `data/`         — Example configs and simulation data
- `gati/`         — GATI protocol and plugin support
- `grpcserv/`     — gRPC server implementation
- `http/`         — HTTP server implementation
- `stats/`        — Statistics and analytics modules

## Example: Generate Reels

```go
import "github.com/zhs007/slotsgamecore7/mathtoolset2"

reels, err := mathtoolset2.GenReels(reader, "SEP_3,SC,WL")
if err != nil {
    // handle error
}
```

## License

This project is licensed under the MIT License.