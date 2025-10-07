# Contribution Notes: Fix Schema Validation Errors

## Summary

This PR fixes schema validation errors in the MCP DigitalOcean server by embedding JSON schema files directly into the Go binary using Go's `embed` package.

## Problem

The MCP DigitalOcean server was experiencing schema validation errors when running in VSCode, Cursor, and other environments. The server would panic on startup with errors like:

```
panic: failed to load app create schema: failed to read schema file app-create-schema.json: 
  open /tmp/app-create-schema.json: no such file or directory
```

This occurred because the schema files for the apps and doks services were being loaded from the filesystem at runtime using `os.ReadFile()`. The fragile implementation relied on:

- Schema files being present in the same directory as the executable
- Successful resolution of the executable path via `os.Executable()`
- The binary being run from the expected directory

This approach failed in various real-world scenarios:
- Running the binary from different directories
- Packaged/containerized environments
- NPM distribution where schema files might not be co-located with the binary
- Symlinked executables

## Solution

Embedded the JSON schema files directly into the Go binary using Go's `embed` package.

### Before

```go
func (a *AppPlatformTool) Tools() []server.ServerTool {
    appCreateSchema, err := loadSchema("app-create-schema.json")
    if err != nil {
        panic(fmt.Errorf("failed to load app create schema: %w", err))
    }
    // ... use appCreateSchema
}

func loadSchema(file string) ([]byte, error) {
    executablePath, err := os.Executable()
    // ... fragile file I/O
}
```

### After

```go
//go:embed spec/app-create-schema.json
var appCreateSchemaJSON []byte

func (a *AppPlatformTool) Tools() []server.ServerTool {
    // ... use appCreateSchemaJSON directly
}
```

## Changes

### Files Modified

1. **internal/apps/apps.go**
   - Added `_ "embed"` import
   - Removed `os` and `path/filepath` imports
   - Added `//go:embed` directives for `app-create-schema.json` and `app-update-schema.json`
   - Removed `loadSchema()` function (18 lines)
   - Updated `Tools()` method to use embedded variables directly

2. **internal/doks/doks.go**
   - Added `_ "embed"` import
   - Removed `os` and `path/filepath` imports
   - Added `//go:embed` directives for `cluster-create-schema.json` and `node-pool-create-schema.json`
   - Removed `loadSchema()` function (18 lines)
   - Updated `Tools()` method to use embedded variables directly

### Statistics

- Total: +18 insertions, -60 deletions
- 2 files changed
- Removed 42 lines of fragile file I/O logic
- Added 6 lines of embed directives and variable declarations

## Benefits

✅ **Self-contained binary** - No external files required at runtime  
✅ **Location-independent** - Works from any directory  
✅ **Eliminates errors** - No more "schema file not found" panics  
✅ **Better performance** - Schemas loaded at compile time, not runtime  
✅ **Simpler distribution** - Single binary contains everything  
✅ **Cleaner code** - Removed 42 lines of fragile file I/O logic  
✅ **Backward compatible** - No breaking changes  

## Testing

All tests pass successfully:

### Unit Tests
```bash
make test
```
- ✅ All packages: `ok` (100% pass rate)
- ✅ Apps package tests pass
- ✅ DOKS package tests pass
- ✅ All other package tests pass

### Linting
```bash
make lint
```
- ✅ No issues found

### Code Formatting
```bash
make format-check
```
- ✅ All files properly formatted

### Build Testing
```bash
go build -o bin/mcp-digitalocean ./cmd/mcp-digitalocean
```
- ✅ Binary builds successfully

### Runtime Testing
Tested the binary from different directories:
```bash
cd /tmp
/path/to/mcp-digitalocean --services apps,doks --digitalocean-api-token test
```
- ✅ No schema file not found errors
- ✅ Server starts successfully
- ✅ Apps service initializes without panic
- ✅ DOKS service initializes without panic

## Verification Steps

To verify this fix works:

1. Build the binary:
   ```bash
   go build -o mcp-digitalocean ./cmd/mcp-digitalocean
   ```

2. Move the binary to a different directory:
   ```bash
   cp mcp-digitalocean /tmp/
   cd /tmp
   ```

3. Run the binary with apps or doks service:
   ```bash
   ./mcp-digitalocean --services apps
   ```

4. Verify no panic about missing schema files occurs

## Additional Notes

- The schema JSON files (`internal/apps/spec/*.json` and `internal/doks/spec/*.json`) are still present in the repository as source files
- They are used at compile time via `//go:embed` directive
- The NPM distribution's Makefile still copies these files to `scripts/npm/dist/`, but they are no longer required at runtime
- This change is fully backward compatible - no changes to the public API or behavior

## Contribution Checklist

- [x] Code compiles without errors
- [x] All tests pass
- [x] Linter passes
- [x] Code is properly formatted
- [x] No breaking changes
- [x] Runtime behavior verified
- [x] Documentation updated (this file)

## Related Issues

This fix addresses schema validation errors reported when running the MCP DigitalOcean server in various environments including VSCode, Cursor, and NPM installations.
