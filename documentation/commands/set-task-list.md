### `set-task-list`

Stores a tasklist to a controller.


```
Request:

topic: <root>/<requests>/device/tasklist:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
            "tasks": "array of task",
            "task.task": "string",
            "task.door": "uint8",
            "task.start-date": "date",
            "task.end-date": "date",
            "task.weekdays": "string list of weekday",
            "task.start": "time",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
tasks        list of task records
task.task    task type
task.door    door for task
task.start-date date from which task is enabled (inclusive)
task.end-date date until which task is enabled (inclusive)
task.weekdays weekdays on which time profile is enabled
task.start   task start time (HHmm)
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "set-task-list",
      "response": {
            "device-id": "uint32",
            "warnings": "array of string",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
warnings     list of warning messages
```


Example:
```
topic: uhppoted/gateway/requests/device/tasklist:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "tasks": [
          {
            "task": "trigger once",
            "door": 3,
            "start-date": "2021-01-01",
            "end-date": "2021-12-31",
            "weekdays": "Monday,Wednesday,Friday",
            "start": "08:27"
          }
        ]
      }
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-time-profile",
      "response": {
        "device-id": 405419896,
        "warnings": []
      }
    }
  }
}
```
