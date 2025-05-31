package converter

import "google.golang.org/protobuf/types/known/wrapperspb"

func ProtoStringToPtrString(s *wrapperspb.StringValue) *string {
	if s == nil {
		return nil
	}
	val := s.GetValue()
	return &val
}

func StringPtrToProtoString(ptr *string) *wrapperspb.StringValue {
	if ptr == nil {
		return nil
	}
	return wrapperspb.String(*ptr)
}
