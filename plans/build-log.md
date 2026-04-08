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
| 2026-04-08 | Phase 2 | constructor with extra return values is not recognized | ✅ | 34964ba |
| 2026-04-08 | Phase 2 | constructor with wrong name prefix is not recognized | ✅ | ba31245 |
| 2026-04-08 | Phase 2 | constructor in a different package is not recognized | ✅ | 8f02b00 |
| 2026-04-08 | Phase 2 | custom type derived from another custom type requires its own constructor | ✅ | 468044e |
| 2026-04-08 | Phase 2 | refactor no-constructor implementation | ✅ | 95d52a2 |
| 2026-04-08 | Phase 3 | bare var declaration of custom type is flagged | ✅ | 92a4bce |
| 2026-04-08 | Phase 3 | custom type obtained via constructor is not flagged | ✅ | 8b10d63 |
| 2026-04-08 | Phase 3 | struct field of custom type with no initializer is flagged | ✅ | 3a20ddb |
| 2026-04-08 | Phase 3 | explicit cast to custom type outside constructor is flagged | ✅ | 6a1b73d |
| 2026-04-08 | Phase 3 | cast inside constructor body is not flagged | ✅ | 52a72c2 |
| 2026-04-08 | Phase 3 | reverse conversion from custom type to scalar is not flagged | ✅ | 8c1b651 |
| 2026-04-08 | Phase 3 | untyped string literal assigned to custom type variable is flagged | ✅ | 8ea1f2d |
| 2026-04-08 | Phase 3 | untyped literal passed as custom type argument is flagged | ✅ | 3ece636 |
| 2026-04-08 | Phase 3 | same-package constant is not flagged | ✅ | 6d32d71 |
| 2026-04-08 | Phase 3 | constant from same package used as value is not flagged | ✅ | a464123 |
| 2026-04-08 | Phase 3 | refactor phase 3 implementation | ✅ | ed4075b |
