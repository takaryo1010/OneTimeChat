package usecase

import (
	"fmt"
	"time"

	"crypto/sha256"
	"encoding/hex"
	"github.com/takaryo1010/OneTimeChat/server/model"
	"math/rand/v2"
)

func generateRoomID(rm *model.RoomManager) string {
	//合計5文字、アルファベット大文字と数字のランダムな文字列を生成
	const chars = "ABCDEFGHJKLMNPQRSTUVWXY0123456789"
	for {
		roomID := ""
		for i := 0; i < 5; i++ {
			roomID += string(chars[rand.IntN(len(chars))])
		}
		if _, exists := rm.Rooms[roomID]; !exists {
			return roomID
		}
	}

}

func appendExpireBinarySearch(rm *model.RoomManager, room *model.Room) {
	// 二分探索で期限切れのルームを探し、適切な位置に挿入
	// 期限切れのルームがない場合は、最後尾に追加
	if len(rm.ExpireSortRooms) == 0 {
		rm.ExpireSortRooms = append(rm.ExpireSortRooms, room)
		return
	}

	low := 0
	high := len(rm.ExpireSortRooms) - 1
	for low <= high {
		mid := (low + high) / 2
		if rm.ExpireSortRooms[mid].Expires.Before(room.Expires) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if low == len(rm.ExpireSortRooms) {
		rm.ExpireSortRooms = append(rm.ExpireSortRooms, room)
	} else {
		rm.ExpireSortRooms = append(rm.ExpireSortRooms[:low], append([]*model.Room{room}, rm.ExpireSortRooms[low:]...)...)
	}
}

func GenerateSessionID() (string, error) {

	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}
