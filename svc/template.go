package svc

import bus "github.com/autonomouskoi/core-tinygo"

func RenderTemplate(template string, json []byte) (string, *bus.Error) {
	msg := &bus.BusMessage{
		Type: int32(MessageType_TEMPLATE_RENDER_REQ),
	}
	bus.MarshalMessage(msg, &TemplateRenderRequest{
		Template: template,
		Json:     json,
	})
	if msg.Error != nil {
		return "", msg.Error
	}
	reply, err := bus.WaitForReply(msg, 1000)
	if err != nil {
		errStr := err.Error()
		return "", &bus.Error{
			Detail: &errStr,
		}
	}
	if reply.Error != nil {
		replyJS, _ := reply.MarshalJSON()
		bus.LogError("reply error", "reply", string(replyJS))
		return "", reply.Error
	}
	var resp TemplateRenderResponse
	if err := resp.UnmarshalVT(reply.GetMessage()); err != nil {
		bus.LogDebug("Unmarshalling")
		errStr := err.Error()
		return "", &bus.Error{
			Detail: &errStr,
		}
	}
	return resp.GetOutput(), nil
}
