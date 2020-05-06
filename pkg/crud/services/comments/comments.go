package comments

import (
	"commentsService/pkg/crud/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
)
type CommentsSvc struct {
	pool *pgxpool.Pool
}

func NewCommentsSvc(pool *pgxpool.Pool) *CommentsSvc {
	if pool == nil {
		panic(errors.New("pool can't be nil")) // <- be accurate
	}
	return &CommentsSvc{pool: pool}
}

func (service *CommentsSvc) Comments(ctx context.Context) (list []models.Comments, err error) {
	list = make([]models.Comments, 0)
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(ctx, "SELECT id, comment, book_id, commentator_name FROM comments WHERE removed = FALSE")
	if err != nil {
		return nil, err // TODO: wrap to specific error
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Comments{}
		err := rows.Scan(&item.Id, &item.Comment, &item.BookId, &item.CommentatorName)
		if err != nil {
			return nil, err // TODO: wrap to specific error
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}



func (service *CommentsSvc) SaveComments( model models.Comments, ctx context.Context) (err error) {
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		log.Print(err)
		return errors.New("can't execute pool: ")
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, "INSERT INTO comments(comment, book_id, commentator_name) VALUES ($1, $2, $3);", model.Comment, model.BookId, "Vasya")
	if err != nil {
		log.Print(err)
		return errors.New("can't save comment: ")
	}
	return nil
}

func (service *CommentsSvc) RemoveCommentsById(ctx context.Context, id int) (err error) {
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		return errors.New("can't execute pool: ")
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, "UPDATE comments SET removed = true where id = $1;", id)
	if err != nil {
		return errors.New("can't remove burger: ")
	}
	return nil
}


func (service *CommentsSvc) ShowCommentsById(id string, ctx context.Context) (list []models.Comments, err error) {
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	bookId, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	rows, err := conn.Query(context.Background(), "SELECT id, comment, commentator_name FROM comments WHERE removed = FALSE AND book_id = $1", int64(bookId))
	if err != nil {
		return nil, err // TODO: wrap to specific error
	}
	defer rows.Close()


	for rows.Next() {
		item := models.Comments{
			Id:              0,
			Comment:         "",
			BookId:          int64(bookId),
			CommentatorName: "",
			Removed:         false,
		}
		err := rows.Scan(&item.Id, &item.Comment, &item.CommentatorName)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

