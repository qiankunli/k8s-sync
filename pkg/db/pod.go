package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"                   // mysql db driver
	_ "github.com/golang-migrate/migrate/database/mysql" // mysql driver for migration
	"github.com/jmoiron/sqlx"
	v1 "k8s.io/api/core/v1"
	"time"

	"github.com/rs/zerolog/log"
)

type PersistentPod struct {
	Id        int64          `db:"id"`
	Env       sql.NullString `db:"env"`
	AppName   sql.NullString `db:"app_name"`
	Isolation sql.NullString `db:"isolation"`
	// k8s apply 后才可能拿到的 属性
	Name        string         `db:"name"`
	ContainerId sql.NullString `db:"container_id"`
	Node        sql.NullString `db:"node"`
	Ip          sql.NullString `db:"ip"`
	Status      sql.NullString `db:"status"`
	CreateAt    sql.NullTime   `db:"create_at"`
	UpdateAt    sql.NullTime   `db:"update_at"` // 可自动赋值当前时间
}

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Open("mysql", "username:password@tcp(ip:port)/k8s?parseTime=true")
	db.SetMaxOpenConns(10)
	if err != nil {
		log.Err(err).Msg("init db error")
	}

}

func NewPersistentPod(pod *v1.Pod, env string) PersistentPod {

	ppod := PersistentPod{
		Env:         sql.NullString{env, true},
		AppName:     sql.NullString{pod.Labels["appName"], true},
		Isolation:   sql.NullString{pod.Labels["isolation"], true},
		Name:        pod.ObjectMeta.Name,
		ContainerId: sql.NullString{},
		Node:        sql.NullString{pod.Status.HostIP, true},
		Ip:          sql.NullString{pod.Status.PodIP, true},
		Status:      sql.NullString{(string)(pod.Status.Phase), true},
	}
	return ppod
}

const saveSql = `insert into tb_pod (env,app_name,isolation,name,container_id,node,ip,status,create_at,update_at) 
					values(?,?,?,?,?,?,?,?,?,?)`

func SaveOrUpdate(pod PersistentPod) error {
	if len(pod.Name) == 0 {
		return errors.New("pod name can not be null")
	}
	if exist, err := Exist(pod.Name); err != nil {
		log.Err(err)
		return err
	} else {
		if exist {
			log.Debug().Msgf("pod name %s is exist,update it", pod.Name)
			UpdateByName(pod)
		} else {
			log.Debug().Msgf("pod name %s is not exist,update it", pod.Name)
			Save(pod)
		}
	}
	return nil
}

func Exist(name string) (bool, error) {
	var count int32
	if err := db.Get(&count, "select count(*) from tb_pod where name = ? limit 0,1", name); err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func Save(pod PersistentPod) error {
	if len(pod.Name) == 0 {
		return errors.New("pod name can not be null")
	}
	if !pod.CreateAt.Valid {
		pod.CreateAt.Scan(time.Now())
	}
	if !pod.UpdateAt.Valid {
		pod.UpdateAt.Scan(time.Now())
	}
	if _, err := db.Exec(saveSql,
		pod.Env,
		pod.AppName,
		pod.Isolation,
		pod.Name,
		pod.ContainerId,
		pod.Node,
		pod.Ip,
		pod.Status,
		pod.CreateAt,
		pod.UpdateAt); err != nil {
		return err
	}
	return nil

}

const updateSql = `update tb_pod set app_name = ?,isolation = ? ,container_id = ? ,node = ? ,ip = ? ,status = ?, update_at = ? where name = ?`

func UpdateByName(pod PersistentPod) error {
	if len(pod.Name) == 0 {
		return errors.New("pod name can not be null")
	}
	pod.UpdateAt.Scan(time.Now())
	if _, err := db.Exec(updateSql,
		pod.AppName,
		pod.Isolation,
		pod.ContainerId,
		pod.Node,
		pod.Ip,
		pod.Status,
		pod.UpdateAt,
		pod.Name); err != nil {
		return err
	}
	return nil
}

const deleteSql = "delete from tb_pod where name=?"

func DeleteByName(name string) error {
	if len(name) == 0 {
		return errors.New("pod name can not be null")
	}
	if result, err := db.Exec(deleteSql, name); err != nil {
		return err
	} else {
		rowsAffected, _ := result.RowsAffected()
		log.Debug().Msgf("delete pod,name = %s,affected %d", name, rowsAffected)
	}
	return nil
}

type PodHandler func(pod PersistentPod) error

const step = 10

func Traverse(podHandler PodHandler) error {
	var minId int64
	if err := db.Get(&minId, "select min(id) from tb_pod"); err != nil {
		log.Err(err)
		return err
	}
	var maxId int64
	if err := db.Get(&maxId, "select max(id) from tb_pod"); err != nil {
		log.Err(err)
		return err
	}
	log.Debug().Msgf("traverse db,minId %d,maxId %d, step 10", minId, maxId)
	successCount := 0
	failCount := 0
	for i := minId; i <= maxId; i += step {
		sc, fc, _ := queryAndHandle(i, i+step-1, podHandler)
		successCount += sc
		failCount += fc
	}
	log.Debug().Msgf("handle %d rows, success %d, fail %d", successCount+failCount, successCount, failCount)
	return nil
}

func queryAndHandle(start, end int64, podHandler PodHandler) (int, int, error) {
	log.Debug().Msgf("query rows, traverse id from %d to %d", start, end)
	rows, err := db.Queryx("SELECT * FROM tb_pod where id >= ? and id <= ?", start, end)
	if err != nil {
		log.Err(err)
		return 0, 0, err
	}
	successCount := 0
	failCount := 0
	for rows.Next() {
		var pod PersistentPod
		err = rows.StructScan(&pod)
		if err != nil {
			log.Err(err)
			continue
		}
		if err := podHandler(pod); err != nil {
			log.Err(err)
			failCount++
			continue
		}
		successCount++
	}
	return successCount, failCount, nil
}

func List() ([]PersistentPod, error) {
	pods := []PersistentPod{}
	if err := db.Select(&pods, "select * from tb_pod"); err != nil {
		return nil, err
	}
	return pods, nil
}
