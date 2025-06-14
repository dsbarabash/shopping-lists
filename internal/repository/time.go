package repository

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockgen -source=time.go -destination=time_mock.go -package=repository

type TimeNower interface {
	Now() *timestamppb.Timestamp
}

type timeNowerImpl struct{}

func NewTimeNower() TimeNower {
	return &timeNowerImpl{}
}

func (i timeNowerImpl) Now() *timestamppb.Timestamp {
	return timestamppb.Now()
}
