# OData Migration Plan - Executive Summary

## Overview

This document provides a high-level summary of the plan to migrate the Clubs backend from custom REST APIs to OData v4 APIs using `github.com/NLstn/go-odata`.

## Current State

- **15 API domains** covering clubs, members, events, fines, notifications, etc.
- **70+ REST endpoints** with custom routing and handling
- **17 primary data models** with complex relationships
- **JWT-based authentication** with magic link and OAuth support
- **Azure integration** for blob storage and email services

## Target State

- **OData v4 API** at `/api/v2/` with full protocol support
- **Standard query capabilities:** $filter, $select, $expand, $orderby, $top, $skip, $count, $search
- **Custom operations** for non-CRUD workflows (actions and functions)
- **100% feature parity** with existing REST API
- **Parallel operation** during migration period
- **Backward compatibility** through gradual deprecation

## Key Benefits

1. **Standardization** - Industry-standard protocol with extensive tooling
2. **Flexibility** - Clients control exactly what data they need
3. **Reduced Code** - Eliminate custom query endpoints and boilerplate
4. **Better Performance** - Selective field retrieval and efficient data loading
5. **Future-Proof** - Easy to extend without breaking changes

## Migration Phases

### Phase 1: Foundation (Weeks 1-2)
- Set up OData service infrastructure
- Annotate models with OData tags
- Register all entities
- Basic CRUD operations working

### Phase 2: Authorization (Weeks 3-4)
- Implement authentication middleware
- Add read/write hooks for authorization
- Test permission enforcement
- Handle soft delete logic

### Phase 3: Custom Operations (Weeks 5-6)
- Implement OData actions (state-changing)
- Implement OData functions (queries)
- Handle special cases (file upload, magic link)
- Test all operations

### Phase 4: Advanced Features (Week 7)
- Enable change tracking for delta queries
- Configure full-text search
- Test complex expansions
- Performance optimization

### Phase 5: Testing (Week 8)
- Integration testing
- OData compliance testing
- Performance testing
- Security audit

### Phase 6: Frontend Migration (Weeks 9-12)
- Update API client
- Leverage OData features
- Test all workflows
- User acceptance testing

### Phase 7: Deprecation (Week 13)
- Mark old endpoints deprecated
- Monitor usage
- Communicate timeline
- Celebrate success! 🎉

## API Mapping Examples

### Before (REST)
```http
GET /api/v1/clubs
GET /api/v1/clubs/{id}
GET /api/v1/clubs/{id}/members
GET /api/v1/clubs/{id}/events
POST /api/v1/clubs/{id}/leave
```

### After (OData)
```http
GET /api/v2/Clubs
GET /api/v2/Clubs('{id}')
GET /api/v2/Members?$filter=clubId eq '{id}'&$expand=User
GET /api/v2/Events?$filter=clubId eq '{id}'&$orderby=startTime
POST /api/v2/Clubs('{id}')/Leave
```

### OData Query Power
```http
# Get club with members and upcoming events in one call
GET /api/v2/Clubs('{id}')?$expand=Members($expand=User),Events($filter=startTime gt now())

# Complex filtering
GET /api/v2/Events?$filter=Club/name eq 'Soccer Club' and startTime gt now()&$orderby=startTime&$top=10

# Count without data
GET /api/v2/Notifications/$count?$filter=userId eq '{id}' and isRead eq false

# Aggregation
GET /api/v2/Fines?$apply=groupby((clubId),aggregate(amount with sum as total))
```

## Technical Approach

### Entity Registration
```go
// Register all entities with OData service
service.RegisterEntity(&Club{})
service.RegisterEntity(&Member{})
service.RegisterEntity(&Event{})
// ... etc
```

### Authentication Integration
```go
// Wrap OData service with JWT middleware
handler := ODataAuthMiddleware(service)
mux.Handle("/api/v2/", handler)
```

### Authorization via Hooks
```go
// Read hook: Filter by user membership
service.RegisterReadHook("Clubs", func(ctx context.Context, query *gorm.DB) (*gorm.DB, error) {
    userID := ctx.Value("userID").(string)
    return query.
        Joins("INNER JOIN members ON members.club_id = clubs.id").
        Where("members.user_id = ?", userID), nil
})
```

### Custom Operations
```go
// Action: Accept invite
service.RegisterAction("Invite", "Accept", acceptInviteHandler)

// Function: Check admin rights
service.RegisterFunction("Club", "IsAdmin", checkAdminHandler)
```

## Risk Mitigation

| Risk | Mitigation |
|------|-----------|
| Learning curve | Comprehensive docs, training, gradual rollout |
| OData complexity | Start simple, provide helpers, examples |
| Performance | Proper indexing, query optimization, monitoring |
| File uploads | Custom endpoints with clear documentation |
| Breaking changes | Parallel operation, gradual deprecation, versioning |

## Success Criteria

✅ API response time < 200ms for simple queries  
✅ Test coverage > 90% for OData endpoints  
✅ 100% feature parity with REST API  
✅ Pass all OData v4 compliance tests  
✅ Positive developer feedback  

## Timeline

- **Weeks 1-2:** Foundation
- **Weeks 3-4:** Authorization
- **Weeks 5-6:** Custom operations
- **Week 7:** Advanced features
- **Week 8:** Testing & validation
- **Weeks 9-12:** Frontend migration
- **Week 13:** Deprecation & cleanup

**Total Duration:** 13 weeks

## Resources Required

- **Backend Development:** 1-2 developers full-time
- **Frontend Development:** 1-2 developers (weeks 9-12)
- **Testing/QA:** 1 tester (ongoing)
- **DevOps:** Support for deployment and monitoring

## Progress Status

### ✅ Completed

- **Phase 1: Foundation** - Entity registration, basic CRUD operations
- **Phase 2: Authorization** - Authentication middleware, read/write hooks  
- **Phase 3: Core CRUD** - Basic entity operations tested
- **Phase 4: Custom Operations** - All actions and functions implemented! 🎉
  - 8 OData Actions (state-changing operations)
  - 8 OData Functions (read-only queries)

### 🚧 In Progress

- **Phase 5: Advanced Features** - Change tracking, full-text search, optimization

### ⏳ Upcoming

- **Phase 6: Testing & Validation** - Integration and compliance testing
- **Phase 7: Frontend Migration** - Client updates and testing
- **Phase 8: Deprecation** - Gradual phase-out of old endpoints

## Recent Updates (December 2025)

**Phase 4 Complete! ✅**

- ✅ Updated go-odata to v0.5.1 with improved documentation
- ✅ Implemented 8 custom actions (Accept/Reject invites, join requests, leave club, etc.)
- ✅ Implemented 8 custom functions (IsAdmin, GetOwnerCount, dashboard queries, search)
- ✅ All backend quality checks passing (go mod verify, go build)
- ✅ Actions and functions registered and ready for use
- ⏳ Integration testing needed for actions and functions

## Next Steps

1. ✅ Review migration plan
2. ✅ Get team approval
3. ✅ Phase 1-3 complete
4. 🚧 Complete Phase 4 (functions + testing)
5. ⏳ Set up weekly progress reviews
6. ⏳ Celebrate milestones along the way!

---

**For detailed technical information, see:** [OData_Migration_Plan.md](./OData_Migration_Plan.md)
