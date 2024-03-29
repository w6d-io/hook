# Hook

[![codecov](https://codecov.io/gh/w6d-io/hook/branch/main/graph/badge.svg?token=HL9LOYYCWI)](https://codecov.io/gh/w6d-io/hook)

It is a package that send in one way term any payload to one or more targets.
It works by subscription and scope.

## subscribe a http target

```go
package main

import (
    "context"
    "os"

    "github.com/w6d-io/hook"
    "github.com/w6d-io/x/logx"
)

type payload struct {
    Name string `json:"name"`
    Kind string `json:"kind"`
}

func main() {
    ctx := context.Background()
    log := logx.WithName(ctx, "Main")
    URL := "http://localhost:8080/test"
    // add a target for the payload for all scope
    if err := hook.Subscribe(ctx, URL, "test"); err != nil {
        log.Error(err, "Subscription failed", "target", URL)
        os.Exit(1)
    }
    kafka := "kafka://localhost:9092?topic=TEST&messagekey=msgid"
    // add a target for the payload for all scope
    if err := hook.Subscribe(ctx, kafka, ".*"); err != nil {
        log.Error(err, "Subscription failed", "target", kafka)
        os.Exit(1)
    }

    p := &payload{
        Name: "payload",
        Kind: "test",
    }
    // Send payload with test as scope the payload is send to http and kafka
    if err := hook.Send(ctx, p, "test"); err != nil {
        log.Error(err, "Send failed")
        os.Exit(1)
    }

    // Send payload with failed as scope. Only sends to kafka due to the scope
    if err := hook.Send(ctx, p, "failed"); err != nil {
        log.Error(err, "Send failed")
        os.Exit(1)
    }

}
```

The subscribe function gets the scope in regex format
