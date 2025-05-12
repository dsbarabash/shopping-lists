package model

import "github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"

type State int

const (
	Archived = shopping_list_api.State_STATE_ARCHIVED
	Active   = shopping_list_api.State_STATE_ACTIVE
)

func (s State) String() string {
	return shopping_list_api.State_name[int32(s)]
}
