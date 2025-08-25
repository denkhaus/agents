---
name: test_prompt
schema:
  type: object
  properties:
    Name:
      type: string
  required:
    - Name
---
Hello, {{.Name}}!
