package main

import (
	"context"
	"log"
	"time"

	permissionpb "github.com/mrhumster/web-server-gin/proto/gen/go/permission"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := permissionpb.NewPermissionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Тестируем CheckPermission
	checkResp, err := client.CheckPermission(ctx, &permissionpb.CheckPermissionRequest{
		UserId:   "user123",
		Resource: "document",
		Action:   "read",
	})
	if err != nil {
		log.Fatalf("CheckPermission failed: %v", err)
	}
	log.Printf("✅ CheckPermission response: %v", checkResp.GetAllowed())

	// Тестируем AddPolicy
	addResp, err := client.AddPolicy(ctx, &permissionpb.AddPolicyRequest{
		Policy:     "user123",
		Resorce:    "document",
		Permission: "read",
	})
	if err != nil {
		log.Fatalf("AddPolicy failed: %v", err)
	}
	log.Printf("AddPolicy response: %v", addResp)
}
