package main

import (
	"log"

	"github.com/takaryo1010/OneTimeChat/server/controller"
	"github.com/takaryo1010/OneTimeChat/server/router"
	"github.com/takaryo1010/OneTimeChat/server/usecase"
	"github.com/takaryo1010/OneTimeChat/server/periodicTask"
)

func main() {


	
	// Usecase と Controller の初期化
	roomUsecase := usecase.NewRoomUsecase()
	mainController := &controller.MainController{
		RoomUsecase: roomUsecase,
	}
	
	// 定期タスクの開始
	go periodicTask.PeriodicTask(roomUsecase.RoomManager)
	
	// ルーターの設定
	e := router.NewRouter(mainController)

	// サーバーを起動
	log.Println("Server started at http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
