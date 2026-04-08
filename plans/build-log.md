# Build Log

| Date | Phase | Scenario | Status | Commit |
|------|-------|----------|--------|--------|
| 2026-04-08 | Phase 1 | struct field with raw string is flagged | ✅ | 2c3ead8 |
| 2026-04-08 | Phase 1 | struct field with raw int is flagged | ✅ | 0b88d8a |
| 2026-04-08 | Phase 1 | struct field with custom type is not flagged | ✅ | f0658b7 |
| 2026-04-08 | Phase 1 | underlying type in a type definition is not flagged | ✅ | 3189443 |
| 2026-04-08 | Phase 1 | function parameter with raw scalar is not flagged | ✅ | dc4d84c |
| 2026-04-08 | Phase 1 | function return type with raw scalar is not flagged | ✅ | 8033a5b |
| 2026-04-08 | Phase 1 | all scalar types in struct fields are flagged | ✅ | 3464598 |
| 2026-04-08 | Phase 1 | refactor + binary entrypoint | ✅ | 6f21097 |
| 2026-04-08 | Phase 2 | custom type with no constructor is flagged | ✅ | 25b12ff |
| 2026-04-08 | Phase 2 | custom type with valid exported constructor is not flagged | ✅ | 62db3cc |
| 2026-04-08 | Phase 2 | custom type with valid unexported constructor is not flagged | ✅ | 347151b |
| 2026-04-08 | Phase 2 | constructor missing error return is not recognized | ✅ | 545d47a |
