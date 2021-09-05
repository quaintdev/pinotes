package main

import (
	"database/sql"
	"fmt"
	"log"
)

type DataStore struct {
	db *sql.DB
}

func (s *DataStore) Init(path string) error {
	var err error
	s.db, err = sql.Open("sqlite3", path)
	if err!=nil{
		return fmt.Errorf("error opening database %s", err)
	}
	mainStmt, _ := s.db.Prepare("CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY, topic TEXT, content TEXT, created_at DATETIME)")
	mainStmt.Exec()
	return nil
}

func (s *DataStore) NewTopic(n Topic) {
	stmt, _ := s.db.Prepare("INSERT INTO notes (topic, content, created_at) VALUES (?,?,datetime('now'))")
	stmt.Exec(n.Topic, n.Content)
}

func (s *DataStore) UpdateTopic(n Topic) {
	_, err:=s.db.Exec("UPDATE notes SET topic = ?, content = ? WHERE id = ?", n.Topic, n.Content, n.Id)
	if err!=nil{
		log.Println(err)
	}
}

func (s *DataStore) ViewTopic(n *Topic) *Topic {
	rows, err := s.db.Query("SELECT id, content FROM notes WHERE topic = ?", n.Topic)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&n.Id, &n.Content)
		if err != nil {
			panic(err)
		}
		break
	}
	return n
}

func (s *DataStore) ListTopics() []string {
	rows, err := s.db.Query("SELECT topic FROM notes")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var topicList []string
	for rows.Next() {
		var topic string
		err := rows.Scan(&topic)
		if err != nil {
			panic(err)
		}
		topicList = append(topicList, topic)
	}
	return topicList
}

func (s *DataStore) DeleteTopic(n Topic) error{
	_, err := s.db.Exec("DELETE FROM notes WHERE topic = ?", n.Topic)
	if err!=nil{
		fmt.Errorf("error deleting topic %s, err: %s", n.Topic, err)
	}
	return nil
}
