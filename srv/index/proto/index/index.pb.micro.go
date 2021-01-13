// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/index/index.proto

package index

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Index service

func NewIndexEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Index service

type IndexService interface {
	GetCodeName(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type indexService struct {
	c    client.Client
	name string
}

func NewIndexService(name string, c client.Client) IndexService {
	return &indexService{
		c:    c,
		name: name,
	}
}

func (c *indexService) GetCodeName(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Index.GetCodeName", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Index service

type IndexHandler interface {
	GetCodeName(context.Context, *Request, *Response) error
}

func RegisterIndexHandler(s server.Server, hdlr IndexHandler, opts ...server.HandlerOption) error {
	type index interface {
		GetCodeName(ctx context.Context, in *Request, out *Response) error
	}
	type Index struct {
		index
	}
	h := &indexHandler{hdlr}
	return s.Handle(s.NewHandler(&Index{h}, opts...))
}

type indexHandler struct {
	IndexHandler
}

func (h *indexHandler) GetCodeName(ctx context.Context, in *Request, out *Response) error {
	return h.IndexHandler.GetCodeName(ctx, in, out)
}
