# Flowt

This is a copy of the AWS Step Functions otherwise known as state machines. It's
goal is to implement the [Amazon States Language](https://states-language.net/spec.html) and allow for loading multiple
step function flows as well as executing them.  It provides some extensions to
make it more useful however, such as generic RPC and shell calls.

# Extension Coverage

The following new tasks are available:

- [x] shell commands: `x-cmd:cat /etc/passwd`
- [x] NATS publishing `x-nats:127.0.0.1:4222/the.nats.subject.to.publish.the.input.on`
- [x] NATS request/response `x-nats-req:127.0.0.1:4222/the.nats.subject.to.request.on`
- [ ] MQTT publishing `x-mqtt:the/mqtt/topic/to/publish/on`
- [ ] HTTP requests `x-http:https://the.url.to/send/input/to`
- [x] RPC request `x-rpc:127.0.0.1:4553/the_function_to_call`

# ASL Coverage

- [ ] context object
- [ ] lambda tasks

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