package flowt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/nats-io/nats.go"
)

type taskHandler func(state *gabs.Container) ([]byte, error)

func commandHandler(cmdstring string) taskHandler {
	return func(state *gabs.Container) ([]byte, error) {
		var bits []string
		stateInline := strings.Contains(cmdstring, "%s")

		if stateInline {
			cmdstring = fmt.Sprintf(cmdstring, state.String())
		}

		bits = strings.Split(cmdstring, " ")
		cmd := exec.Command(bits[0])

		if len(bits) > 1 {
			os.Args = append(os.Args, bits[1:]...)
		}

		if !stateInline {
			cmd.Stdin = bytes.NewReader(state.Bytes())
		}

		cmd.Stderr = os.Stderr
		data, err := cmd.Output()
		// exc.setLog(logbuf.String())
		if err != nil {
			return nil, err
		}

		return data, nil
	}
}

func waitForState(s Step, input *gabs.Container) error {
	duration, err := waitDuration(s, input)
	if err != nil {
		return err
	}

	time.Sleep(duration)
	return nil
}

func waitDuration(s Step, input *gabs.Container) (time.Duration, error) {
	switch {
	case s.Seconds > 0:
		return time.Duration(s.Seconds) * time.Second, nil

	case s.SecondsPath != "":
		jseconds, ok := input.Path(s.SecondsPath).Data().(json.Number)
		if !ok {
			return 0, errors.New("seconds at given path is not a number (or does not exist)")
		}

		seconds, err := jseconds.Int64()
		if err != nil {
			return 0, err
		}
		return time.Duration(seconds) * time.Second, nil

	case s.TimestampPath != "":
		timestamp, ok := input.Path(s.TimestampPath).Data().(string)
		if !ok {
			return 0, errors.New("timestamp at given path is not a string (or does not exist)")
		}

		t, err := time.Parse("2006-01-02T15:04:05Z0700", timestamp)
		if err != nil {
			return 0, err
		}
		return time.Until(t), nil

	case s.Timestamp != "":
		t, err := time.Parse("2006-01-02T15:04:05Z0700", s.Timestamp)
		if err != nil {
			return 0, err
		}
		return time.Until(t), nil
	}

	return 0, errors.New("failed to determine wait time")
}

func rpcHandler(addr string) taskHandler {
	return func(state *gabs.Container) ([]byte, error) {
		bits := strings.Split(addr, "/")
		if len(bits) != 2 {
			return nil, errors.New("bad format for RPC resource")
		}

		addr = bits[0]
		proc := bits[1]

		cl, err := rpc.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}

		var data string
		err = cl.Call(proc, state.String(), &data)
		if err != nil {
			return nil, err
		}

		return []byte(data), nil
	}
}

func natsReqHandler(addr string) taskHandler {
	return func(state *gabs.Container) ([]byte, error) {
		bits := strings.Split(addr, "/")
		if len(bits) != 2 {
			return nil, errors.New("bad format for NATS resource")
		}

		addr = bits[0]
		subj := bits[1]

		cl, err := nats.Connect(addr)
		if err != nil {
			return nil, err
		}

		msg, err := cl.RequestWithContext(context.Background(), subj, state.Bytes())
		if err != nil {
			return nil, err
		}

		return msg.Data, nil
	}
}

func natsHandler(addr string) taskHandler {
	return func(state *gabs.Container) ([]byte, error) {
		bits := strings.Split(addr, "/")
		if len(bits) != 2 {
			return nil, errors.New("bad format for NATS resource")
		}

		addr = bits[0]
		subj := bits[1]

		cl, err := nats.Connect(addr)
		if err != nil {
			return nil, err
		}

		err = cl.Publish(subj, state.Bytes())
		return []byte("{}"), err
	}
}

func httpHandler(addr string) taskHandler {
	return func(state *gabs.Container) ([]byte, error) {
		res, err := http.Post(addr, "application/json", bytes.NewBuffer(state.Bytes()))
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		res.Body.Close()

		b := map[string]interface{}{}
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, err
		}

		x := map[string]interface{}{
			"code":    res.StatusCode,
			"body":    b,
			"headers": res.Header,
		}

		return json.Marshal(x)
	}
}
