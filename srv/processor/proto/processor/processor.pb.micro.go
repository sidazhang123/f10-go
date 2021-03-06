// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/processor/processor.proto

package processor

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

// Api Endpoints for Processor service

func NewProcessorEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Processor service

type ProcessorService interface {
	Process(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	RegexTest(ctx context.Context, in *RegexReq, opts ...client.CallOption) (*RegexRsp, error)
	BuildSo(ctx context.Context, in *BuildSoReq, opts ...client.CallOption) (*BuildSoRsp, error)
	GetSourceCode(ctx context.Context, in *GetSourceCodeReq, opts ...client.CallOption) (*GetSourceCodeRsp, error)
	GetPluginPath(ctx context.Context, in *GetPluginPathReq, opts ...client.CallOption) (*GetPluginPathRsp, error)
}

type processorService struct {
	c    client.Client
	name string
}

func NewProcessorService(name string, c client.Client) ProcessorService {
	return &processorService{
		c:    c,
		name: name,
	}
}

func (c *processorService) Process(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Processor.Process", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processorService) RegexTest(ctx context.Context, in *RegexReq, opts ...client.CallOption) (*RegexRsp, error) {
	req := c.c.NewRequest(c.name, "Processor.RegexTest", in)
	out := new(RegexRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processorService) BuildSo(ctx context.Context, in *BuildSoReq, opts ...client.CallOption) (*BuildSoRsp, error) {
	req := c.c.NewRequest(c.name, "Processor.BuildSo", in)
	out := new(BuildSoRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processorService) GetSourceCode(ctx context.Context, in *GetSourceCodeReq, opts ...client.CallOption) (*GetSourceCodeRsp, error) {
	req := c.c.NewRequest(c.name, "Processor.GetSourceCode", in)
	out := new(GetSourceCodeRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processorService) GetPluginPath(ctx context.Context, in *GetPluginPathReq, opts ...client.CallOption) (*GetPluginPathRsp, error) {
	req := c.c.NewRequest(c.name, "Processor.GetPluginPath", in)
	out := new(GetPluginPathRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Processor service

type ProcessorHandler interface {
	Process(context.Context, *Request, *Response) error
	RegexTest(context.Context, *RegexReq, *RegexRsp) error
	BuildSo(context.Context, *BuildSoReq, *BuildSoRsp) error
	GetSourceCode(context.Context, *GetSourceCodeReq, *GetSourceCodeRsp) error
	GetPluginPath(context.Context, *GetPluginPathReq, *GetPluginPathRsp) error
}

func RegisterProcessorHandler(s server.Server, hdlr ProcessorHandler, opts ...server.HandlerOption) error {
	type processor interface {
		Process(ctx context.Context, in *Request, out *Response) error
		RegexTest(ctx context.Context, in *RegexReq, out *RegexRsp) error
		BuildSo(ctx context.Context, in *BuildSoReq, out *BuildSoRsp) error
		GetSourceCode(ctx context.Context, in *GetSourceCodeReq, out *GetSourceCodeRsp) error
		GetPluginPath(ctx context.Context, in *GetPluginPathReq, out *GetPluginPathRsp) error
	}
	type Processor struct {
		processor
	}
	h := &processorHandler{hdlr}
	return s.Handle(s.NewHandler(&Processor{h}, opts...))
}

type processorHandler struct {
	ProcessorHandler
}

func (h *processorHandler) Process(ctx context.Context, in *Request, out *Response) error {
	return h.ProcessorHandler.Process(ctx, in, out)
}

func (h *processorHandler) RegexTest(ctx context.Context, in *RegexReq, out *RegexRsp) error {
	return h.ProcessorHandler.RegexTest(ctx, in, out)
}

func (h *processorHandler) BuildSo(ctx context.Context, in *BuildSoReq, out *BuildSoRsp) error {
	return h.ProcessorHandler.BuildSo(ctx, in, out)
}

func (h *processorHandler) GetSourceCode(ctx context.Context, in *GetSourceCodeReq, out *GetSourceCodeRsp) error {
	return h.ProcessorHandler.GetSourceCode(ctx, in, out)
}

func (h *processorHandler) GetPluginPath(ctx context.Context, in *GetPluginPathReq, out *GetPluginPathRsp) error {
	return h.ProcessorHandler.GetPluginPath(ctx, in, out)
}
