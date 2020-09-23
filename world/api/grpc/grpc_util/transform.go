package grpc_util

import (
	"time"

	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

func ToProtoDuration(d time.Duration) *duration.Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return &duration.Duration{Seconds: int64(secs), Nanos: int32(nanos)}
}

func ToProtoTimestamp(t time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{Seconds: int64(t.Unix()), Nanos: int32(t.Nanosecond())}
}

func FromProtoTimestamp(x *timestamp.Timestamp) time.Time {
	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanos()))
}
