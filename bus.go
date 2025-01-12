package core

import (
	"fmt"

	"github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user send
func send(busMessage uint64)

//go:wasmimport extism:host/user send_reply
func sendReply(busMessage uint64)

//go:wasmimport extism:host/user wait_for_reply
func waitForReply(busMessage uint64, timeoutMS uint64) uint64

// Send a BusMessage to the host
func Send(msg *BusMessage) error {
	mem, err := MarshalArg(msg)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}
	send(mem.Offset())
	mem.Free()
	return nil
}

// SendReply sends a BusMessage to the host that is a reply to a message
// received. The ReplyTo field should be set to the value from the received
// message.
func SendReply(msg *BusMessage) error {
	mem, err := MarshalArg(msg)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}
	sendReply(mem.Offset())
	mem.Free()
	return nil
}

// WaitForReply sends a message to the host and waits for a reply. If a is not
// received within timeoutMS milliseconds, the returned BusMessage.Error.Code
// will be CommonErrorCode_TIMEOUT.
func WaitForReply(msg *BusMessage, timeoutMS uint64) (*BusMessage, error) {
	mem, err := MarshalArg(msg)
	if err != nil {
		return nil, fmt.Errorf("marshalling: %w", err)
	}
	defer mem.Free()
	offs := waitForReply(mem.Offset(), timeoutMS)
	return UnmarshalReturn(offs)
}

// Marshaller represents a proto that can be marshalled, suitable for tinygo.
type Marshaller interface {
	MarshalVT() ([]byte, error)
}

// Unmarshaller represents a proto that can be unmarshalled, suitable for tinygo.
type Unmarshaller interface {
	UnmarshalVT([]byte) error
}

// UnmarshallerPTR represents a value type where a pointer to that value
// implements Unmarshaller
type UnmarshallerPTR[M any] interface {
	*M
	Unmarshaller
}

// WaitForReplyWrap is a convenience function that marshals req, sends it via
// WaitForReply using msg, unmarshals the response into a RESP, and invokes
// process on the result. An error in marshalling req, unmarshalling the RESP,
// or in calling WaitForReply is returned directly. If the reply message from
// WaitForReply is non-nil, that is passed to process. The process function will
// receive a non-nil resp or non-nil *bus.Error, but not both.
func WaitForReplyWrap[M any, REQ Marshaller, RESP UnmarshallerPTR[M]](
	msg *BusMessage, req REQ, process func(RESP, *Error), timeoutMS uint64,
) error {
	var err error
	msg.Message, err = req.MarshalVT()
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}
	reply, err := WaitForReply(msg, timeoutMS)
	if err != nil {
		return fmt.Errorf("waiting for reply: %w", err)
	}
	if reply.Error != nil {
		process(nil, reply.Error)
		return nil
	}
	var resp M
	if err := RESP(&resp).UnmarshalVT(reply.GetMessage()); err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	process(&resp, nil)
	return nil
}

// Subscribe to a topic
func Subscribe(topic string) error {
	b, err := (&SubscribeRequest{
		Topic: topic,
	}).MarshalVT()
	if err != nil {
		return err
	}
	msg := &BusMessage{
		Type:    int32(ExternalMessageType_SUBSCRIBE_REQ),
		Message: b,
	}
	return Send(msg)
}

// Unsubscribe from a topic. If no such subscription exists no error is returned.
func Unsubscribe(topic string) error {
	b, err := (&UnsubscribeRequest{
		Topic: topic,
	}).MarshalVT()
	if err != nil {
		return err
	}
	msg := &BusMessage{
		Type:    int32(ExternalMessageType_SUBSCRIBE_REQ),
		Message: b,
	}
	return Send(msg)
}

// MarshalArg marshals the provided argument for passing as an argument to an
// invocation of a host function. You probably want to use Send, SendReply, or
// WaitForReply instead.
func MarshalArg(msg *BusMessage) (pdk.Memory, error) {
	b, err := msg.MarshalVT()
	if err != nil {
		return pdk.Memory{}, err
	}
	mem := pdk.AllocateBytes(b)
	return mem, nil
}

// UnmarshalReturn unmarshals a message provided as the return value of the
// invocation of a host function. You probably want to use Send, SendReply, or
// WaitForReply instead.
func UnmarshalReturn(offs uint64) (*BusMessage, error) {
	mem := pdk.FindMemory(offs)
	defer mem.Free()
	msg := &BusMessage{}
	err := msg.UnmarshalVT(mem.ReadBytes())
	return msg, err
}

// Error implements the built in error interface
func (e *Error) Error() string {
	return e.GetDetail()
}
