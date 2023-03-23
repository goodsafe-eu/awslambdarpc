// package client provides the function to make a RPC request to a
// lambda function and read the response.
package client

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

// Invoke makes the request to the local lambda function running
// on the address specified in addr.
// data is the payload used in the request, eg. a JSON to be passed
// to the lambda function as body.
// If the lambda returned an error then this function will return
// the error message in the error interface
func Invoke(addr string, data []byte, deadLineSeconds int64) ([]byte, error) {
	t := time.Now().Add(time.Duration(deadLineSeconds) * time.Second)
	deadline := &messages.InvokeRequest_Timestamp{
		Seconds: int64(t.Unix()),
		Nanos:   int64(t.Nanosecond()),
	}
	request := messages.InvokeRequest{Payload: data, Deadline: *deadline}
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	var r messages.InvokeResponse
	err = client.Call("Function.Invoke", request, &r)
	if err != nil {
		return nil, err
	}

	if r.Error != nil {
		return nil, fmt.Errorf("lambda returned error:\n%s", r.Error.Message)
	}

	return r.Payload, nil
}
