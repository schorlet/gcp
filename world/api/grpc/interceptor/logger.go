package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func UnaryClientLogger(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	log.Printf("RPC client message: %T, value: %v\n", req, req)
	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		log.Printf("RPC client method: %q, failed: %v\n", method, err)
	}

	log.Printf("RPC client message: %T, value: %v\n", reply, reply)
	log.Printf("RPC client method: %q, took: %s\n", method, time.Since(start))

	return err
}

func UnaryServerLogger(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log.Printf("RPC server message: %T, value: %v\n", req, req)
	start := time.Now()

	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC server method: %q, failed: %v\n", info.FullMethod, err)
	}

	log.Printf("RPC server method: %q, took: %s\n", info.FullMethod, time.Since(start))

	return m, err
}

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	log.Printf("RPC server message: %T, value: %v\n", m, m)
	return err
}

// func (w *wrappedStream) SendMsg(m interface{}) error {
// log.Printf("RPC send message: %T\n", m)
// return w.ServerStream.SendMsg(m)
// }

func StreamServerLogger(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	start := time.Now()

	wrapped := wrappedStream{ss}

	err := handler(srv, &wrapped)
	if err != nil {
		log.Printf("RPC server method: %q, failed: %v\n", info.FullMethod, err)
	}

	log.Printf("RPC server method: %q, took: %s\n", info.FullMethod, time.Since(start))

	return err
}
