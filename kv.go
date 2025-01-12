package core

import (
	"github.com/autonomouskoi/akcore"
)

// KVGet retrieves a value from the KV store. If no value with that key is
// present, akcore.ErrNotFound will be returned
func KVGet(key []byte) ([]byte, error) {
	msg := &BusMessage{
		Type: int32(ExternalMessageType_KV_GET_REQ),
	}
	req := &KVGetRequest{Key: key}

	var value []byte
	var innerErr error
	err := WaitForReplyWrap(msg, req, func(resp *KVGetResponse, busErr *Error) {
		if busErr != nil {
			if busErr.GetCode() == int32(CommonErrorCode_NOT_FOUND) {
				innerErr = akcore.ErrNotFound
				return
			}
			innerErr = busErr
			return
		}
		value = resp.GetValue()
	}, 1000)
	if err != nil {
		return nil, err
	}
	return value, innerErr
}

// KVGetProto retrieves the value associated with key from the KV store and
// unmarshals it into p.
func KVGetProto(key []byte, p Unmarshaller) error {
	value, err := KVGet(key)
	if err != nil {
		return err
	}
	return p.UnmarshalVT(value)
}

// KVSet sets a value in the KV store with the specified key. If there's an
// existing value with that key it is overwritten
func KVSet(key, value []byte) error {
	msg := &BusMessage{
		Type: int32(ExternalMessageType_KV_SET_REQ),
	}
	req := &KVSetRequest{Key: key, Value: value}
	var err error
	err = WaitForReplyWrap(msg, req, func(resp *KVSetResponse, busErr *Error) {
		err = busErr
	}, 1000)
	return err
}

// KVSetProto marshals p and sets key to that value in the KV store.
func KVSetProto(key []byte, p Marshaller) error {
	value, err := p.MarshalVT()
	if err != nil {
		return err
	}
	return KVSet(key, value)
}

// KVList lists keys matching a given prefix or all values if prefix is nil.
// The limit parameter limits the number of matching keys returned. If limit
// is 0, all matches are returned. The offset parameter can be used to skip
// matches. If offset is greater than or equal to the total number of matches,
// no matches are returned. Pagination can be implemented by using the same
// limit in successive calls and specifying the offset to be the total number
// of matches retrieved until it reaches the total number of matches.
func KVList(prefix []byte, limit, offset int) (*KVListResponse, error) {
	msg := &BusMessage{
		Type: int32(ExternalMessageType_KV_LIST_REQ),
	}
	req := &KVListRequest{
		Prefix: prefix,
		Limit:  uint32(limit),
		Offset: uint32(offset),
	}
	var resp *KVListResponse
	var err error
	err = WaitForReplyWrap(msg, req, func(gotResp *KVListResponse, busErr *Error) {
		resp, err = gotResp, busErr
	}, 1000)
	return resp, err
}

// KVDelete deletes the value associated with the provided key. If there's no
// value with that key no error is returned.
func KVDelete(key []byte) error {
	msg := &BusMessage{
		Type: int32(ExternalMessageType_KV_DELETE_REQ),
	}
	req := &KVDeleteRequest{Key: key}
	var err error
	err = WaitForReplyWrap(msg, req, func(_ *KVDeleteResponse, busErr *Error) {
		err = busErr
	}, 1000)
	return err
}
