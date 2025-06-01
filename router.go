package core

// A Handler handles bus messages, optionally returning a reply
type Handler func(*BusMessage) *BusMessage

// TypeRouter dispatches a message to its respective handler by message type
type TypeRouter map[int32]Handler

// Handle a message based on type. If there's no handler for the type no action
// is taken
func (r TypeRouter) Handle(msg *BusMessage) *BusMessage {
	handler, present := r[msg.GetType()]
	if !present {
		return nil
	}
	return handler(msg)
}

// A TopicRouter dispatches a message to the appropriate TypeHandler by topic.
type TopicRouter map[string]TypeRouter

// Handle a message using the TypeHandler for the type. If there's no handler
// for the topic no action is taken.
func (r TopicRouter) Handle(msg *BusMessage) {
	tr, present := r[msg.GetTopic()]
	if !present {
		return
	}
	reply := tr.Handle(msg)
	if reply == nil {
		return
	}
	reply.ReplyTo = msg.ReplyTo
	SendReply(reply)
}

// DefaultReply creates a template reply by copying msg's topic and incrementing
// the message's type
func DefaultReply(msg *BusMessage) *BusMessage {
	return &BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
}

// MarshalMessage marshals v to the Message field of msg. If marshalling fails
// an error is logged and msg.Error is set.
func MarshalMessage(msg *BusMessage, v Marshaller) {
	var err error
	msg.Message, err = v.MarshalVT()
	if err != nil {
		errStr := err.Error()
		LogError("marshalling", "error", errStr)
		msg.Error = &Error{
			Code:        int32(CommonErrorCode_INVALID_TYPE),
			UserMessage: &errStr,
		}
	}
}

// UnmarshalMessage unmarshals msg.Message to v. If unmarshalling fails an error
// is logged and an Error is returned.
func UnmarshalMessage(msg *BusMessage, v Unmarshaller) *Error {
	if err := v.UnmarshalVT(msg.GetMessage()); err != nil {
		errStr := err.Error()
		LogError("unmarshalling", "error", errStr)
		return &Error{
			Code:        int32(CommonErrorCode_INVALID_TYPE),
			UserMessage: &errStr,
		}
	}
	return nil
}
