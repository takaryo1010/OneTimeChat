package usecase

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/takaryo1010/OneTimeChat/server/model"
)

func Test_generateRoomID(t *testing.T) {
	type args struct {
		rm *model.RoomManager
	}
	tests := []*struct {
		name string
		args args
		want *regexp.Regexp
	}{
		{
			name: "Test generateRoomID",
			args: args{
				rm: &model.RoomManager{
					Rooms: map[string]*model.Room{},
				},
			},
			want: regexp.MustCompile(`^[A-Z0-9]{5}$`),
		},
		{
			name: "Test generateRoomID with existing room",
			args: args{
				rm: &model.RoomManager{
					Rooms: map[string]*model.Room{
						"12345": nil,
						"ABCDE": nil,
					},
				},
			},
			want: regexp.MustCompile(`^[A-Z0-9]{5}$`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRoomID(tt.args.rm)
			fmt.Printf("t: %v\n", got)
			if !tt.want.MatchString(got) {
				t.Errorf("generateRoomID() = %v, does not match %v", got, tt.want.String())
			}
		})
	}
}

func Test_appendExpireBinarySearch(t *testing.T) {
	type args struct {
		rm   *model.RoomManager
		room *model.Room // ポインタ型に修正
	}
	tests := []*struct {
		name string
		args args
		want []*model.Room
	}{
		{
			name: "Test appendExpireBinarySearch with empty list",
			args: args{
				rm: &model.RoomManager{
					ExpireSortRooms: []*model.Room{}, // ポインタ型で初期化
				},
				room: &model.Room{
					ID:      "12345",
					Name:    "test",
					Owner:   "test",
					Expires: time.Unix(0, 0),
				},
			},
			want: []*model.Room{
				{
					ID:      "12345",
					Name:    "test",
					Owner:   "test",
					Expires: time.Unix(0, 0),
				},
			},
		},
		{
			name: "Test appendExpireBinarySearch with existing room",
			args: args{
				rm: &model.RoomManager{
					ExpireSortRooms: []*model.Room{
						{
							ID:      "ABCDE",
							Name:    "test",
							Owner:   "test",
							Expires: time.Unix(1, 0),
						},
					},
				},
				room: &model.Room{
					ID:      "12345",
					Name:    "test",
					Owner:   "test",
					Expires: time.Unix(0, 0),
				},
			},
			want: []*model.Room{
				{
					ID:      "12345",
					Name:    "test",
					Owner:   "test",
					Expires: time.Unix(0, 0),
				},
				{
					ID:      "ABCDE",
					Name:    "test",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
			},
		},
		{
			name: "Test appendExpireBinarySearch with multiple rooms",
			args: args{
				rm: &model.RoomManager{
					ExpireSortRooms: []*model.Room{
						{
							ID:      "ABCDE",
							Name:    "test1",
							Owner:   "test",
							Expires: time.Unix(1, 0),
						},
						{
							ID:      "XYZ",
							Name:    "test2",
							Owner:   "test",
							Expires: time.Unix(2, 0),
						},
					},
				},
				room: &model.Room{
					ID:      "12345",
					Name:    "test3",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
			},
			want: []*model.Room{
				{
					ID:      "ABCDE",
					Name:    "test1",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
				{
					ID:      "12345",
					Name:    "test3",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
				{
					ID:      "XYZ",
					Name:    "test2",
					Owner:   "test",
					Expires: time.Unix(2, 0),
				},
			},
		},
		{
			name: "Test appendExpireBinarySearch with room after all",
			args: args{
				rm: &model.RoomManager{
					ExpireSortRooms: []*model.Room{
						{
							ID:      "ABCDE",
							Name:    "test1",
							Owner:   "test",
							Expires: time.Unix(1, 0),
						},
					},
				},
				room: &model.Room{
					ID:      "XYZ",
					Name:    "test2",
					Owner:   "test",
					Expires: time.Unix(3, 0),
				},
			},
			want: []*model.Room{
				{
					ID:      "ABCDE",
					Name:    "test1",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
				{
					ID:      "XYZ",
					Name:    "test2",
					Owner:   "test",
					Expires: time.Unix(3, 0),
				},
			},
		},
		{
			name: "Test appendExpireBinarySearch with room before all",
			args: args{
				rm: &model.RoomManager{
					ExpireSortRooms: []*model.Room{
						{
							ID:      "ABCDE",
							Name:    "test1",
							Owner:   "test",
							Expires: time.Unix(1, 0),
						},
					},
				},
				room: &model.Room{
					ID:      "XYZ",
					Name:    "test2",
					Owner:   "test",
					Expires: time.Unix(0, 0),
				},
			},
			want: []*model.Room{
				{
					ID:      "XYZ",
					Name:    "test2",
					Owner:   "test",
					Expires: time.Unix(0, 0),
				},
				{
					ID:      "ABCDE",
					Name:    "test1",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
			},
		},
		{
			name: "Test appendExpireBinarySearch with identical expiration times",
			args: args{
				rm: &model.RoomManager{
					ExpireSortRooms: []*model.Room{
						{
							ID:      "ABCDE",
							Name:    "test1",
							Owner:   "test",
							Expires: time.Unix(1, 0),
						},
					},
				},
				room: &model.Room{
					ID:      "XYZ",
					Name:    "test2",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
			},
			want: []*model.Room{
				{
					ID:      "ABCDE",
					Name:    "test1",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
				{
					ID:      "XYZ",
					Name:    "test2",
					Owner:   "test",
					Expires: time.Unix(1, 0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 引数をポインタ型で渡す
			appendExpireBinarySearch(tt.args.rm, tt.args.room)
			for i := range tt.args.rm.ExpireSortRooms {
				fmt.Printf("t: %v %v\n", tt.args.rm.ExpireSortRooms[i].ID, tt.args.rm.ExpireSortRooms[i].Expires)
			}
			if len(tt.args.rm.ExpireSortRooms) != len(tt.want) {
				t.Errorf("appendExpireBinarySearch() = %v, want %v", tt.args.rm.ExpireSortRooms, tt.want)
			}
			for i := range tt.args.rm.ExpireSortRooms {
				if tt.args.rm.ExpireSortRooms[i].Expires != tt.want[i].Expires {
					t.Errorf("appendExpireBinarySearch() = %v, want %v", tt.args.rm.ExpireSortRooms, tt.want)
				}
			}
		})
	}
}

