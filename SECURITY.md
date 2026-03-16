# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | Yes                |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public GitHub issue
2. Use [GitHub Security Advisories](https://github.com/justinthevoid/recipe/security/advisories/new) to report privately
3. Or email the maintainer directly

## Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 1 week
- **Fix timeline**: Depends on severity, typically within 2 weeks for critical issues

## Scope

Recipe processes files entirely client-side (browser via WebAssembly, or local CLI). There are no servers, databases, or network services. Security concerns are primarily:

- File parsing vulnerabilities (malformed NP3/XMP input)
- WASM sandbox escapes (mitigated by browser sandboxing)
- Supply chain risks (Go/npm dependencies)
