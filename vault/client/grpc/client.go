package grpc

import (
	"google.golang.org/grpc"
	"github.com/ruandao/goPractice/vault"
	"github.com/ruandao/goPractice/vault/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func New(conn *grpc.ClientConn) vault.Service {
	var hashEndpoint = grpctransport.NewClient(
		conn, "Vault", "Hash",
		vault.EncodeGRPCHashRequest,
		vault.DecodeGRPCHashResponse,
		pb.HashResponse{},
	).Endpoint()
	var validateEndpoint = grpctransport.NewClient(
		conn, "Vault", "Validate",
		vault.EncodeGRPCValidateRequest,
		vault.DecodeGRPCValidateResponse,
		pb.ValidateResponse{},
	).Endpoint()
	return vault.Endpoints{
		HashEndpoint:hashEndpoint,
		ValidateEndPoint:validateEndpoint,
	}
}