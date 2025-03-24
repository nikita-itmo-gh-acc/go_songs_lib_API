package storage

import (
	"database/sql"
	"fmt"
	"songsapi/logger"
	"songsapi/query"
)

type Song struct {
	Id			int		`json:"id"`
	Name        string 	`json:"song"`
	Group       string 	`json:"group,omitempty"`
	ReleaseDate string 	`json:"releaseDate,omitempty"`
	Text        string 	`json:"text,omitempty"`
	Link        string 	`json:"link,omitempty"`
	GroupId		int		`json:"groupId,omitempty"`
}

type SongStorage struct {
	DB *sql.DB
}

func (s *SongStorage) Get(id int) (*Song, error) {
	song := Song{}
	err := s.DB.QueryRow("SELECT * FROM songs WHERE id = $1", id).Scan(
		&song.Id, &song.GroupId, &song.Name, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		logger.Err.Println("can't find song with id = ", id)
		return nil, err
	}
	
	return &song, err
}

func (s *SongStorage) Create(song *Song) error {
	_, err := s.DB.Exec(`INSERT INTO songs ("groupId", "name", "releaseDate", "text", "link") VALUES ($1, $2, $3, $4, $5)`, 
						song.GroupId, song.Name, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		logger.Err.Println("can't insert into songs table - ", err)
		return err
	}

	return nil
}

func (s *SongStorage) Delete(song *Song) error {
	_, err := s.DB.Exec(`DELETE FROM songs WHERE id = $1`, song.Id)
	if err != nil {
		logger.Err.Println("can't delete from songs table - ", err)
		return err
	}

	return nil
}

func (s *SongStorage) Update(song *Song) error {
	_, err := s.DB.Exec(`UPDATE songs SET "groupId" = $1, "name" = $2, "releaseDate" = $3, "text" = $4, "link" = $5
						WHERE songs.id = $6`, song.GroupId, song.Name, song.ReleaseDate, song.Text, song.Link, song.Id)
	
	if err != nil {
		logger.Err.Println("can't update songs table - ", err)
		return err
	}
	return nil			
}

func (s *SongStorage) Find(q query.Query) ([]*Song, error) {
	songQuery, ok := q.(*query.SongQuery)
	if !ok {
		return nil, fmt.Errorf("can't convert search query into songQuery")
	}

	reserved_capacity := 100
	songs := make([]*Song, 0, reserved_capacity)

	query := songQuery.GenerateSQL()

	// fmt.Println("Generated SQL:", query)

	rows, err := s.DB.Query(query)

	if err != nil {
		logger.Err.Println("error during songs search - ", err)
		return nil, err
	}

	noRowsFound := true

	defer rows.Close()
	for rows.Next() {
		noRowsFound = false
		song := Song{}
		if err := rows.Scan(&song.Id, &song.Name, &song.ReleaseDate, &song.Text, &song.Link, &song.Group); err != nil {
			logger.Err.Println("can't scan songs table row:", err)
            continue
		}
		songs = append(songs, &song)
	}

	if noRowsFound {
		return nil, sql.ErrNoRows
	}

	return songs, nil
}
