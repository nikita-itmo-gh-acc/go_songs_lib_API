package storage

import (
	"database/sql"
	"fmt"
	"songsapi/logger"
	"songsapi/query"
)

type Group struct {
	Id			int
	Name 		string
}

type GroupStorage struct {
	DB *sql.DB
}

func (s *GroupStorage) Get(id int) (*Group, error) {
	group := Group{}

	if err := s.DB.QueryRow("SELECT * FROM groups WHERE id = $1", id).Scan(&group.Id, &group.Name); err != nil {
		logger.Err.Println("can't find group with id = ", id)
		return nil, err
	}

	return &group, nil
}

func (s *GroupStorage) Create(group *Group) error {
	_, err := s.DB.Exec(`INSERT INTO groups (name) VALUES ($1)`, group.Name)
	if err != nil {
		logger.Err.Println("can't insert into groups table - ", err)
		return err
	}

	return nil
}

func (s *GroupStorage) Delete(group *Group) error {
	_, err := s.DB.Exec(`DELETE FROM groups WHERE id = $1`, group.Id)
	if err != nil {
		logger.Err.Println("can't delete from groups table - ", err)
		return err
	}

	return nil
}

func (s *GroupStorage) Update(group *Group) error {
	_, err := s.DB.Exec(`UPDATE groups SET name = $1 WHERE id = $2`, group.Name, group.Id)
	if err != nil {
		logger.Err.Println("can't update groups table - ", err)
		return err
	}

	return nil
}

func (s *GroupStorage) Find(q query.Query) ([]*Group, error) {
	groupQuery, ok := q.(*query.GroupQuery)
	if !ok {
		return nil, fmt.Errorf("can't convert search query into groupQuery")
	}

	reserved_capacity := 100
	groups := make([]*Group, 0, reserved_capacity)

	rows, err := s.DB.Query(`SELECT * FROM groups WHERE name LIKE '%' || $1 || '%'`, groupQuery.Name)

	if err != nil {
		logger.Err.Println("groups search failed - ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		group := Group{}
		if err := rows.Scan(&group.Id, &group.Name); err != nil {
			logger.Err.Println("can't scan groups row:", err)
			continue
		}
		groups = append(groups, &group)
	}

	return groups, nil
}