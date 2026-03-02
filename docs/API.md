# API Reference

## Base URL
http://localhost:8080

## Endpoints

### Submit Task
POST /task
```json
{
  "type": "hello_world",
  "params": {}
}
```

### Get Task
GET /task/:id

### List Tasks
GET /tasks
