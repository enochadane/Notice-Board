package repository

import (
	"database/sql"
	"errors"

	"github.com/amthesonofGod/Notice-Board/entity"
)

//PostRepositoryImpl ...
type PostRepositoryImpl struct {
	conn *sql.DB
}

//NewPostRepositoryImpl ...
func NewPostRepositoryImpl(Conn *sql.DB) *PostRepositoryImpl {
	return &PostRepositoryImpl{conn: Conn}
}

//Posts ...
func (pri *PostRepositoryImpl) Posts() ([]entity.Post, error) {

	rows, err := pri.conn.Query("SELECT * FROM posts;")
	if err != nil {
		return nil, errors.New("could not query the database")
	}
	defer rows.Close()

	posts := []entity.Post{}

	for rows.Next() {
		post := entity.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.Category)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil

}

//Post ...
func (pri *PostRepositoryImpl) Post(id int) (entity.Post, error) {

	row := pri.conn.QueryRow("SELECT * FROM posts WHERE id = $1", id)

	post := entity.Post{}

<<<<<<< HEAD:post/repository/psql_post.go
	err := row.Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.Category)
=======
	err := row.Scan(&post.ID, &post.CompanyID, &post.Title, &post.Description, &post.Image, &post.Category)
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f:model/repository/psql_post.go
	if err != nil {
		return post, err
	}

	return post, nil

}

//UpdatePost ...
func (pri *PostRepositoryImpl) UpdatePost(post entity.Post) error {

	// _, err := pri.conn.Exec("UPDATE posts SET title=$1, description=$2, image=$3, category=$4 WHERE id=$5", post.Title, post.Description, post.Image, post.Category, post.Id)
	// if err != nil {
	// 	return errors.New("Update has failed")
	// }

	return nil
}

//DeletePost ...
func (pri *PostRepositoryImpl) DeletePost(id int) error {

	_, err := pri.conn.Exec("DELETE FROM posts WHERE id=$1", id)
	if err != nil {
		return errors.New("Delete has failed")
	}

	return nil
}

//StorePost ...
func (pri *PostRepositoryImpl) StorePost(post entity.Post) error {
<<<<<<< HEAD:post/repository/psql_post.go
	
=======

>>>>>>> 56480e1450127de4cec062eea25b723b5216035f:model/repository/psql_post.go
	_, err := pri.conn.Exec("INSERT INTO posts (title,description,image,category) values($1, $2, $3,$4)", post.Title, post.Description, post.Image, post.Category)
	if err != nil {
		return errors.New("Insertion has failed")
	}

	return nil
}
