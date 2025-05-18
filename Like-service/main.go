package main

import (
	"fmt"

	"github.com/shrey209/Like-Service/repository"
	"github.com/shrey209/Like-Service/service"
)

func main() {
	session := NewSession("127.0.0.1", "like_service")

	repo := repository.NewLikeRepository(session)
	svc := service.NewLikeService(repo)

	err := svc.LikePost("post-123")
	if err != nil {
		fmt.Println("Error liking post:", err)
	}

	count, err := svc.GetPostLikes("post-123")
	if err != nil {
		fmt.Println("Error getting likes:", err)
	}

	fmt.Println("Post has likes:", count)
}
