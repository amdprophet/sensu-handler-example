# Sensu Go Handler Template

The [Sensu Go][1] handler is a [Sensu Event Handler][2] that sends a payload, generated from a Sensu
event, to an HTTP API.

## Configuration

Example Sensu Go handler definition:

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "handler-template"
    },
    "spec": {
        "type": "pipe",
        "command": "sensu-handler-template --url mycompany.tld -t 15",
        "timeout": 10,
        "filters": [
            "is_incident"
        ]
    }
}
```

[1]: https://github.com/sensu/sensu-go
[2]: https://docs.sensu.io/sensu-go/latest/reference/handlers/#how-do-sensu-handlers-work
