# Data Model: API Health Endpoint

**Feature**: Enhanced health endpoint for TODO application
**Design Date**: 2025-09-27

## Entities

### HealthStatus
**Purpose**: Represents the current operational state of the service

**Fields**:
- `status` (string, required): Overall service health ("healthy", "degraded", "unhealthy")
- `database` (string, required): Database connectivity status ("connected", "disconnected", "error")
- `timestamp` (string, required): ISO 8601 timestamp when check was performed
- `version` (string, optional): Application version identifier
- `uptime` (number, optional): Service uptime in seconds

**Validation Rules**:
- `status` must be one of: "healthy", "degraded", "unhealthy"
- `database` must be one of: "connected", "disconnected", "error"
- `timestamp` must be valid ISO 8601 format
- `version` should be semantic version if provided
- `uptime` must be positive number if provided

**State Transitions**:
- `healthy`: All systems operational (database connected)
- `degraded`: Some issues detected (database problems but service running)
- `unhealthy`: Critical failures (service cannot operate)

### HealthResponse
**Purpose**: HTTP response wrapper for health check results

**Fields**:
- `health` (HealthStatus, required): Current health status
- `message` (string, optional): Human-readable status message
- `details` (object, optional): Additional diagnostic information

**JSON Schema**:
```json
{
  "status": "healthy|degraded|unhealthy",
  "database": "connected|disconnected|error",
  "timestamp": "2025-09-27T10:30:00Z",
  "version": "1.0.0",
  "uptime": 3600
}
```

## Relationships

### Service → Database
- Service checks database connectivity
- Database status influences overall service status
- One-to-one relationship (single database)

### HealthCheck → HealthStatus
- Each health check produces one status snapshot
- Stateless operation (no persistence required)
- Real-time evaluation on each request

## Implementation Notes

### Go Model Structure
```go
type HealthStatus struct {
    Status    string  `json:"status"`
    Database  string  `json:"database"`
    Timestamp string  `json:"timestamp"`
    Version   string  `json:"version,omitempty"`
    Uptime    int64   `json:"uptime,omitempty"`
}
```

### Database Integration
- Use existing GORM connection from storage package
- Test connection with `db.DB().Ping()`
- No new database tables required
- Leverages existing database initialization

### Response Patterns
- HTTP 200: Healthy status
- HTTP 503: Degraded/unhealthy status
- HTTP 500: Unexpected errors
- Consistent JSON structure across all responses