package periodicTask

import (
	"fmt"
	"time"

	"github.com/takaryo1010/OneTimeChat/server/model"
)

func PeriodicTask(rm *model.RoomManager) {
	ticker := time.NewTicker(5 * time.Minute) // 5分ごとに実行
	defer ticker.Stop()

	for range ticker.C {
		// fmt.Println("定期タスク実行:", time.Now())
		// 期限切れのルームを削除
		deleteExpireSortRooms(rm)
		// ここでDB更新やログ処理などを行う

	}
}

func deleteExpireSortRooms(rm *model.RoomManager) {
	now := time.Now()
	newRooms := make([]*model.Room, 0, len(rm.ExpireSortRooms)) // ポインタのスライスに変更

	for _, room := range rm.ExpireSortRooms {
		if room.Expires.After(now) {
			newRooms = append(newRooms, room) // ポインタ型のまま追加
		} else {
			delete(rm.Rooms, room.ID) // 期限切れなら削除
			fmt.Println("Delete room:", room.ID, "by PeriodicTask")
		}
	}

	rm.ExpireSortRooms = newRooms // ポインタのスライスを代入
}
