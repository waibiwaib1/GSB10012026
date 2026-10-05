# Producer API

## Produce message

Producer API allows you to send one or more messages to a topic, optionally specifying key or partition for each record.

* Endpoint : `http://hostip:port/topics/{topic_name}`
* Request method : `POST`
* Path params : `topic_name`
* Request body params : `records` , `records[i].partition [optional]` , `records[i].key [optional]` , `records[i].value`

Request body should be in JSON format.

Sample request for sending three records:

The records specify a key, a partition, or only a value.

```
POST /topics/kafka-bridge HTTP/1.1
Host: localhost:8080
Content-Type: application/application/json
{
    "records": [
        {
            "key": "my-key",
            "value": "Hi this is kafka-bridge (with key)"
        },
        {
            "value": "Hi this is kafka-bridge (with partition)",
            "partition": 1
        },
        {
            "value": "Hi this is kafka-bridge"
        }
    ]
}
```

Response is a JSON object containing metadata for all produced records, in the same order as the request.

Sample response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "offsets": [
        {
            "partition": 2,
            "offset": 0
        },
        {
            "partition": 1,
            "offset": 0
        },
        {
            "partition": 2,
            "offset": 1
        }
    ]
}
```

In case of errors, the corresponding offset entry contains `error_code` and `error` fields.
