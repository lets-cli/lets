---
id: example_js
title: Example for JavaScript/Node.js
---

**`lets.yaml`**

```yaml
shell: bash

commands:
  run:
    description: Run node server
    cmd: npm run server

  webpack:
    description: Run webpack
    cmd: 
      - npm 
      - run 
      - webpack  
      
  tests:
    cmd: 
      - npm 
      - run 
      - test
```


Examples of usage:

- `lets run` - run server
- `lets webpack -w` - cmd is an array so all arguments will be appended to that array
- `lets test` - run all tests
- `lets test src/server/__tests__` - run only tests in particular directory