package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    "day2/grpc/gen"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil { log.Fatal(err) }
    defer conn.Close()

    c := gen.NewUserServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Пример вызова: GetUser
    resp, err := c.GetUser(ctx, &gen.GetUserRequest{Id: 1})
    if err != nil {
        log.Println("GetUser error:", err)
        return
    }
    log.Printf("user: %+v\n", resp.GetUser())
}


