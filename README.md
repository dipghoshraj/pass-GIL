# pass-GIL

A proof of concept (PoC) demonstrating how to bypass Python's Global Interpreter Lock (GIL) by leveraging Go routines through C ABI calls. This project showcases a hybrid approach where computationally intensive tasks are offloaded to Go, which can utilize true parallelism via goroutines, while maintaining a Python interface.

## Overview

Python's GIL prevents multiple threads from executing Python bytecode simultaneously, which can be a bottleneck for CPU-bound tasks. This project circumvents this limitation by:

1. Implementing the core logic in Go, which doesn't have a GIL and supports true parallelism
2. Exposing Go functions through a C-compatible interface
3. Calling the Go code from Python using CFFI (C Foreign Function Interface)
4. Using FlatBuffers for efficient, schema-driven serialization between languages

## Architecture

The project consists of three main components:

### Go Engine (`goengine/`)
- **engine.go**: Main C interface with exported functions `CallEngine` and `FreeBuffer`
- **registry.go**: Function registry and dispatcher for different operations
- **hasher.go**: Parallel hashing implementation using goroutines
- **engine/**: FlatBuffers-generated Go structs for engine-level messages
- **simple/**: FlatBuffers-generated Go structs for hash-specific messages

### Python Engine (`pyengine/`)
- **hasher.py**: Python client that loads the Go DLL and performs hashing operations
- **engine/**: FlatBuffers-generated Python structs for engine messages
- **simple/**: FlatBuffers-generated Python structs for hash messages

### Schema (`schema/`)
- **engine.fbs**: FlatBuffers schema for generic request/response envelopes
- **simple.fbs**: FlatBuffers schema for hash-specific data structures

## Prerequisites

- Go 1.25.5 or later
- Python 3.12 or later
- FlatBuffers compiler (`flatc`)
- C compiler (for building the Go DLL)

## Dependencies

### Python
- `cffi`: For calling C functions from Python
- `flatbuffers`: For serialization/deserialization

### Go
- `github.com/google/flatbuffers/go`: Go FlatBuffers library

## Build Instructions

### 1. Generate FlatBuffers Code

First, generate the necessary FlatBuffers code from the schemas:

```bash
# Generate Go code
flatc --go -o goengine/engine schema/engine.fbs
flatc --go -o goengine/simple schema/simple.fbs

# Generate Python code
flatc --python -o pyengine/engine schema/engine.fbs
flatc --python -o pyengine/simple schema/simple.fbs
```

### 2. Build Go DLL

Navigate to the `goengine` directory and build the shared library:

```bash
cd goengine
go build -buildmode=c-shared -o ../bin/goengine.dll .
```

This creates `goengine.dll` (Windows) or `goengine.so` (Linux/macOS) in the `bin/` directory.

### 3. Set Up Python Environment

Create a virtual environment and install dependencies:

```bash
python -m venv venv
venv\Scripts\activate  # On Windows
pip install cffi flatbuffers
```

## Usage

### Python Example

```python
from pyengine.hasher import hash_data

# Example data to hash
data = ["hello", "world", "python", "golang"]

# Perform parallel hashing using Go routines
results = hash_data(data)

for original, hash_value in results:
    print(f"{original} -> {hash_value}")
```

### Direct DLL Usage

The `pyengine/hasher.py` file demonstrates direct usage of the Go DLL:

```python
import os
from cffi import FFI

# Load the Go DLL
lib_path = os.path.join("bin", "goengine.dll")
C = ffi.dlopen(lib_path)

# Prepare data and call the engine
# ... (see hasher.py for complete implementation)
```

## How It Works

1. **Python Side**: Data is serialized using FlatBuffers into a `Request` envelope containing:
   - Function ID (e.g., 1 for hashing)
   - Request ID
   - Payload (serialized `HashRequest`)

2. **C Interface**: The `CallEngine` function receives the serialized data, deserializes it, and dispatches to the appropriate handler.

3. **Go Processing**: The hash handler:
   - Deserializes the `HashRequest`
   - Spawns goroutines (up to 64 workers) to process data in parallel
   - Each goroutine computes SHA-256 hashes
   - Results are collected and serialized back

4. **Response**: Results are serialized as `HashOutPut` containing an array of `HashResponse` objects, wrapped in a `Response` envelope.

## Performance Benefits

- **True Parallelism**: Go goroutines can run on multiple CPU cores simultaneously
- **No GIL Contention**: CPU-bound work happens outside Python's GIL
- **Efficient Serialization**: FlatBuffers provides zero-copy, fast serialization
- **C ABI**: Direct C calls minimize overhead

## Project Structure

```
pass-GIL/
├── README.md
├── bin/
│   ├── goengine.dll    # Built Go shared library
│   └── goengine.h      # C header file
├── goengine/
│   ├── engine.go       # Main C interface
│   ├── registry.go     # Function dispatcher
│   ├── hasher.go       # Parallel hashing logic
│   ├── go.mod          # Go module file
│   ├── go.sum          # Go dependencies
│   ├── engine/         # Generated FlatBuffers Go code
│   └── simple/         # Hash-specific FlatBuffers Go code
├── pyengine/
│   ├── hasher.py       # Python client
│   ├── engine/         # Generated FlatBuffers Python code
│   └── simple/         # Hash-specific FlatBuffers Python code
├── schema/
│   ├── engine.fbs      # Generic message schema
│   └── simple.fbs      # Hash operation schema
└── venv/               # Python virtual environment
```

## Future Enhancements

- Add more function handlers beyond hashing
- Implement connection pooling for multiple concurrent requests
- Add error handling and recovery mechanisms
- Create Python package with proper setup.py
- Add benchmarking scripts to demonstrate performance gains

## Contributing

This is a proof of concept project. Contributions are welcome, especially for:
- Additional function handlers
- Performance optimizations
- Cross-platform build improvements
- Documentation enhancements
