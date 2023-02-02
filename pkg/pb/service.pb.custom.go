package pb

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

func OutToPB(in *[]interface{}) (*ScriptResponse, error) {
	inCP := *in
	out := make([]*structpb.Value, len(inCP))
	for i, v := range inCP {
		tmp, err := structpb.NewValue(fmt.Sprintf("%v", v))
		if err != nil {
			return nil, err
		}
		out[i] = tmp
	}

	return &ScriptResponse{Result: out}, nil
}
