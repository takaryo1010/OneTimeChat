package periodicTask

import (
	"fmt"
	"time"

	"github.com/takaryo1010/OneTimeChat/server/model"
)

func PeriodicTask(rm *model.RoomManager) {
	ticker := time.NewTicker(1 *time.Second) // 1分ごとに実行
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// fmt.Println("定期タスク実行:", time.Now())
			// 期限切れのルームを削除
			deleteExpireSortRooms( rm)
			// ここでDB更新やログ処理などを行う
		}
	}
}

func deleteExpireSortRooms( rm *model.RoomManager){
	// 期限切れのルームを削除

	for i, room := range rm.ExpireSortRooms {
		if room.Expires.Before(time.Now()) {
			rm.ExpireSortRooms = append(rm.ExpireSortRooms[:i], rm.ExpireSortRooms[i+1:]...)
			delete(rm.Rooms, room.ID)
			fmt.Println("Delete room:", room.ID)
		}else if(room.Expires.After(time.Now())){
			break
		}
	}
}