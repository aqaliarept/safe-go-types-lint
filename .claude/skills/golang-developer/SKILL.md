---
name: golang-developer
description: Use the skill when developing in Go 
---

# Rules
1. Try always use custom types for scalars and avoid language scalars as much as possible
`type Address string`
2. Every custom type must have constructor difined. Constructor name should start with `New` and type name. Constructor should contain type's invariant check and return error if invariant violated.
```
func NewAddress(val string) (Address, error) {
  if val == ""{
    return "", ErrIsEmpty
  }
  return Address(val), nil
}
```
3. Main design principle is **Tell, Don't Ask**
-  Prefer OOP style over procedural: Create methods to address type's data instead of creating functions accepting types.
-  Encapsulation: Don't expose data from structs. Fields must be private by default. Data mutation should be done via command methods. They should check preconditions and invariants. Expose getters only when really needed.
- Use public fields in structs when this required by other reasons, e.g. serialization requirements
