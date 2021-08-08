# Flowt

This is a copy of the AWS Step Functions otherwise known as state machines. It's
goal is to implement the [Amazon States Language](https://states-language.net/spec.html) and allow for loading multiple
step function flows as well as executing them.  It provides some extensions to
make it more useful however, such as generic RPC and shell calls.

Flowt aims to be a drop in replacement for Step Functions and takes Amazon States Language (with some extensions).

## Why?

While I like the step functions idea, it's implementation as the AWS product is a little restrictive in what it lets you do. 

# Extension Coverage

The following new tasks are available:

- [x] shell commands: `x-cmd:cat /etc/passwd`
- [x] NATS publishing `x-nats:127.0.0.1:4222/the.nats.subject.to.publish.the.input.on`
- [x] NATS request/response `x-nats-req:127.0.0.1:4222/the.nats.subject.to.request.on`
- [ ] MQTT publishing `x-mqtt:the/mqtt/topic/to/publish/on`
- [ ] HTTP requests `x-http:https://the.url.to/send/input/to`
- [x] RPC request `x-rpc:127.0.0.1:4553/the_function_to_call`
- [ ] TCP requests `x-tcp:127.0.0.1:5555`

# ASL Coverage

- [ ] context object
- [ ] lambda tasks
- [ ] activity tasks
- [ ] activities

- [ ] parameters
- [x] input paths
- [x] result paths
- [ ] result selector
- [x] output paths

- [ ] choice state (started)
- [x] pass state
- [x] fail state
- [x] succeed state
- [x] task state
- [ ] parallel state (this one looks hard)
- [ ] map state (this one looks really hard)

- [ ] catch
- [ ] retry

- [ ] intrinsic functions
- [ ] predefined errors

## Building

    make

## Usage


Start `flowt` with an Amazon States file, JSON or YAML is OK:

    flowt -f flows.yml -nats nats://localhost:4222 -bind :9090

You can start an execution with a given JSON input via the API or the NATS channel:

    curl -XPOST localhost:9090/api/v1/executions -d '{"hello": "world"}'
    < {"id": "8e53b604-f774-11eb-8628-6bbb552b9e4b"}
    
    natspub -s "flowt.execution.create.<flow-name>" '{"hello": "world"}'

When starting via NATS you can use the request/reply method if you want to get the ID:

```golang
reply, err := nc.RequestWithContext(ctx, "flowt.execution.create.<flow-name>", data)
// prints: 8e53b604-f774-11eb-8628-6bbb552b9e4b
fmt.Println(string(reply.Data))
```

You can get an execution like so:

    curl localhost:9090/api/v1/executions/8e53b604-f774-11eb-8628-6bbb552b9e4b
    
    
