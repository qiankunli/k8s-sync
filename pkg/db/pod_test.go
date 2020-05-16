package db

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	pod := PersistentPod{
		AppName:     sql.NullString{"demo3", true},
		Isolation:   sql.NullString{"stable", true},
		Name:        "xfm-demo1-xxx",
		ContainerId: sql.NullString{"abcdefg", true},
		Node:        sql.NullString{"192.168.1.1", true},
		Ip:          sql.NullString{"172.1.1.1", true},
		Status:       sql.NullString{"running", true},
	}
	err := Save(pod)
	fmt.Println(err)
}

func TestUpdateByName(t *testing.T) {
	pod := PersistentPod{
		AppName:     sql.NullString{"demo3", true},
		Isolation:   sql.NullString{"stable", true},
		Name:        "xfm-demo1-xxx",
		ContainerId: sql.NullString{"abcdefg", true},
		Node:        sql.NullString{"192.168.1.1", true},
		Ip:          sql.NullString{"172.1.1.1", true},
		Status:       sql.NullString{"running", true},
	}
	err := UpdateByName(pod)
	fmt.Println(err)
}

func TestList(t *testing.T) {
	if pods, err := List(); err == nil {
		for _, pod := range pods {
			fmt.Println(pod.Name)
			fmt.Println(pod.Id)
			// 必须是这个时间点, 据说是go诞生之日
			fmt.Println(pod.CreateAt.Time.Format("2006-01-02 15:04:05"))
			fmt.Println(pod.UpdateAt)
		}
	}else {
		fmt.Println(err)
	}
}

func TestTraverse(t *testing.T) {
	handler := func(pod PersistentPod) error {
		fmt.Printf("list pod:%s\n",pod.Name)
		return nil
	}
	if err := Traverse(handler); err != nil{
		fmt.Println(err)
	}
}

func TestExist(t *testing.T) {
	exist,err := Exist("inter-enter-live-category-business-stable-66bd575dcc-s8fh7")
	fmt.Println(err)
	fmt.Println(exist)
}
