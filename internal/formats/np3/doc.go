// Package np3 provides parsing, verification, and generation for Nikon Picture Control (.np3) binary files.
//
// This package is at the leaf level of the import graph. It provides core functionality
// to read strictly-typed NP3 structures, perform checksum validations, and act as the single source
// of truth for data mapping. It is imported by both `cmd/np3tool` (for the VS Code extension backend)
// and `internal/formats/nksc` (to embed parsed NP3 structures).
//
// Its primary responsibilities include secure binary read/write operations without memory explosion,
// calculating embedded magic hash checksums, and safely mapping internal byte structures
// to JSON-compatible exports for Inter-Process Communication (IPC).
//
// Functions operating on binary payloads within this package are fully stateless
// providing strong thread-safety guarantees. Optimized buffer reads ensure
// processing characteristics remain well under sub-millisecond latencies per file.
package np3
