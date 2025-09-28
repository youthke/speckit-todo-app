# Quickstart: API Health Endpoint Testing

**Feature**: Enhanced health endpoint for TODO application
**Testing Date**: 2025-09-27

## Prerequisites

1. **Backend server running**:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```
   Server should be available at `http://localhost:8080`

2. **Database initialized**: SQLite database should be created and connected

## Core User Scenarios

### Scenario 1: Healthy Service Check
**As a monitoring service, verify that all systems are operational**

**Steps**:
1. Open terminal or API client (curl/Postman)
2. Make GET request to health endpoint:
   ```bash
   curl -X GET http://localhost:8080/health
   ```

**Expected Result**:
- HTTP Status: 200 OK
- Response body:
  ```json
  {
    "status": "healthy",
    "database": "connected",
    "timestamp": "2025-09-27T10:30:00Z",
    "version": "1.0.0",
    "uptime": 3600
  }
  ```

**Validation**:
- [x] Status is "healthy"
- [x] Database is "connected"
- [x] Timestamp is current ISO 8601 format
- [x] Response time < 200ms

### Scenario 2: Database Connectivity Verification
**As a system administrator, verify database connectivity status**

**Steps**:
1. Ensure database is properly connected
2. Make health check request:
   ```bash
   curl -X GET http://localhost:8080/health -w "%{http_code}"
   ```

**Expected Result**:
- HTTP Status: 200
- Database field shows "connected"
- Response includes valid timestamp

**Validation**:
- [x] Database connectivity reported accurately
- [x] Health check doesn't impact app performance
- [x] Response includes diagnostic timestamp

### Scenario 3: Service Status for Monitoring
**As a monitoring tool, poll health endpoint regularly**

**Steps**:
1. Set up repeated requests (simulate monitoring):
   ```bash
   for i in {1..5}; do
     curl -X GET http://localhost:8080/health
     sleep 2
   done
   ```

**Expected Result**:
- All requests return 200 OK
- Consistent response structure
- Updated timestamps on each request
- No performance degradation

**Validation**:
- [x] Endpoint handles repeated requests reliably
- [x] Timestamps update correctly
- [x] No memory leaks or performance issues
- [x] Response time remains < 200ms

### Scenario 4: Service Version Information
**As an operations team, verify deployed service version**

**Steps**:
1. Make health check request
2. Verify version information is included:
   ```bash
   curl -X GET http://localhost:8080/health | jq '.version'
   ```

**Expected Result**:
- Version field present in response
- Valid semantic version format
- Consistent across requests

**Validation**:
- [x] Version information included
- [x] Version format is readable
- [x] Helps track deployments

### Scenario 5: Health Check Response Time
**As a monitoring system, ensure health checks are fast enough for frequent polling**

**Steps**:
1. Measure response time:
   ```bash
   curl -X GET http://localhost:8080/health -w "\nResponse time: %{time_total}s\n"
   ```
2. Repeat multiple times to get average

**Expected Result**:
- Response time consistently < 200ms
- Minimal variance between requests
- No timeout errors

**Validation**:
- [x] Response time meets monitoring requirements
- [x] Performance suitable for frequent polling
- [x] No blocking operations

## Edge Cases

### Edge Case 1: Database Connection Issues
**Test health endpoint behavior when database is unavailable**

**Simulation**:
1. Stop database or corrupt database file
2. Make health check request
3. Verify appropriate error response

**Expected Behavior**:
- HTTP Status: 503 Service Unavailable
- Response indicates database issue
- Service remains responsive

### Edge Case 2: High Load Testing
**Test health endpoint under load**

**Simulation**:
1. Generate multiple concurrent requests:
   ```bash
   ab -n 100 -c 10 http://localhost:8080/health
   ```

**Expected Behavior**:
- All requests succeed or fail gracefully
- No service crashes
- Consistent response format

### Edge Case 3: Startup Health Checks
**Test health endpoint during service startup**

**Simulation**:
1. Restart service
2. Immediately poll health endpoint
3. Verify startup behavior

**Expected Behavior**:
- Service reports status accurately during startup
- Database connection status reflects reality
- No false positives

## Manual Testing Checklist

### Functional Testing
- [ ] Health endpoint returns 200 when service is healthy
- [ ] Database connectivity is accurately reported
- [ ] Timestamp is current and properly formatted
- [ ] Version information is included and accurate
- [ ] Response structure matches contract specification

### Performance Testing
- [ ] Response time < 200ms under normal load
- [ ] Health checks don't impact main application performance
- [ ] Endpoint handles concurrent requests appropriately
- [ ] No memory leaks during repeated requests

### Error Handling
- [ ] Appropriate HTTP status codes for different scenarios
- [ ] Graceful handling of database connection issues
- [ ] Service remains responsive during health check failures
- [ ] Error responses include useful diagnostic information

### Integration Testing
- [ ] Health endpoint works with existing CORS configuration
- [ ] Logging middleware captures health check requests appropriately
- [ ] Security headers are applied consistently
- [ ] Frontend can successfully call enhanced health endpoint

## Troubleshooting

### Common Issues

**Issue**: Health endpoint returns 404
**Solution**: Verify server is running and route is configured

**Issue**: Database always shows "disconnected"
**Solution**: Check database initialization and GORM connection

**Issue**: Response time too slow
**Solution**: Verify database connection pooling and query optimization

**Issue**: Missing version information
**Solution**: Ensure version is properly configured in application build

### Debug Commands

```bash
# Check if server is running
curl -I http://localhost:8080/health

# Verbose health check with timing
curl -v -w "@curl-format.txt" http://localhost:8080/health

# Test with different HTTP methods
curl -X OPTIONS http://localhost:8080/health
```

## Success Criteria

✅ **All scenarios pass validation checks**
✅ **Response time consistently < 200ms**
✅ **Database connectivity accurately reported**
✅ **Monitoring-friendly response format**
✅ **Compatible with existing application architecture**