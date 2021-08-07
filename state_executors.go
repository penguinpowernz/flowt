package flowt

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Jeffail/gabs"
)

type stateExecutor func(*gabs.Container) (string, []byte, error)

func passExecutor(step Step) stateExecutor {
	return func(state *gabs.Container) (string, []byte, error) {
		data, err := json.Marshal(step.Result)
		return step.Next, data, err
	}
}

func waitExecutor(step Step) stateExecutor {
	return func(state *gabs.Container) (string, []byte, error) {
		err := waitForState(step, state)
		return step.Next, nil, err
	}
}

func choiceExecutor(step Step) stateExecutor {
	return func(state *gabs.Container) (string, []byte, error) {
		for _, ch := range step.Choices {
			if ch.IsSatisfied(state) {
				return ch.Next, nil, nil
			}
		}

		// TODO: throw error if no choice nor default is found

		return step.Default, nil, nil
	}
}

func taskExecutor(step Step) stateExecutor {
	typ := strings.Split(step.Resource, ":")[0]
	resource := strings.TrimPrefix(step.Resource, typ+":")
	var hdlr taskHandler

	switch typ {
	case "x-cmd":
		hdlr = commandHandler(resource)
	// case "arn":
	// TODO: handle lambda
	case "x-http":
		hdlr = httpHandler(resource)
	case "x-rpc":
		hdlr = rpcHandler(resource)
	case "x-nats":
		hdlr = natsHandler(resource)
	case "x-nats-req":
		hdlr = natsReqHandler(resource)
	}

	return func(state *gabs.Container) (string, []byte, error) {
		if hdlr == nil {
			return "", nil, errors.New("unsupported resource type: " + typ)
		}

		data, err := hdlr(state)
		return step.Next, data, err
	}
}
