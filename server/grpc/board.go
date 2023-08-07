package grpc

import (
	"context"
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Board struct {
	BoardServer
}

func (b *Board) CreateSubject(ctx context.Context, newSubject *NewSubject) (*Subject, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if len(newSubject.GetTitle()) == 0 {
		log.Errorf("CreateSubject: invalid input 'title'")
		return nil, errors.New("invalid 'title'")
	}

	if err := insertSubject(db, newSubject.Title); err != nil {
		log.Errorf("CreateSubject: %s", err)
		return nil, err
	}

	subject, err := selectSubjectByTitle(db, newSubject.Title)
	if err != nil {
		log.Errorf("CreateSubject: failed to select created subject. %s", err)
		return nil, err
	}

	return subject, nil
}

func (b *Board) DeleteSubject(ctx context.Context, subjectId *SubjectId) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if subjectId.GetId() == 0 {
		log.Errorf("DeleteSubject: invalid subject id;")
		return nil, errors.New("invalid 'id'")
	}

	if err := deleteSubject(db, subjectId.Id); err != nil {
		log.Errorf("DeleteSubject: %s", err)
		return nil, err
	}

	return nil, nil
}

func (b *Board) ListSubject(ctx context.Context, empty *emptypb.Empty) (*SubjectList, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	rows, err := db.Query("SELECT id, title FROM subject ORDER BY id;")
	if err != nil {
		log.Errorf("ListSubject: %s", err)
		return nil, err
	}
	defer rows.Close()

	var list []*Subject

	for rows.Next() {
		var id int64
		var title string

		if err := rows.Scan(&id, title); err != nil {
			log.Fatalf("ListSubject: %s", err)
		}

		list = append(list, &Subject{
			Id:    id,
			Title: title,
		})
	}

	return &SubjectList{
		SubjectList: list,
	}, nil
}

func (b *Board) CreateQuestion(ctx context.Context, newQuestion *NewQuestion) (*Question, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if len(newQuestion.GetQuestion()) == 0 {
		log.Errorf("CreateQuestion: empty input 'question'")
		return nil, errors.New("'question' is empty")
	}

	insertedId, err := insertQuestion(db, newQuestion.Question, newQuestion.SubjectId)
	if err != nil {
		log.Errorf("CreateQuestion: %s", err)
		return nil, err
	}

	question, err := selectQuestion(db, insertedId)
	if err != nil {
		log.Errorf("CreateQuestion: failed to select created question. %s", err)
		return nil, err
	}

	return question, nil
}

func (b *Board) DeleteQuestion(ctx context.Context, questionId *QuestionId) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if questionId.GetId() == 0 {
		log.Errorf("DeleteQuestion: invalid question id;")
		return nil, errors.New("invalid 'id'")
	}

	if err := deleteQuestion(db, questionId.Id); err != nil {
		log.Errorf("DeleteQuestion: %s", err)
		return nil, err
	}

	return nil, nil
}

func (b *Board) GetQuestion(ctx context.Context, questionId *QuestionId) (*Question, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	question, err := selectQuestion(db, questionId.Id)
	if err != nil {
		log.Errorf("GetQuestion: %s", err)
		return nil, err
	}

	return question, nil
}

func (b *Board) ListQuestion(ctx context.Context, subjectId *SubjectId) (*QuestionList, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	likesRows, err := db.Query("SELECT question_id, COUNT(*) AS count FROM likes WHERE subject_id = ? GROUP BY question_id ORDER BY count DESC;", subjectId)
	if err != nil {
		log.Errorf("ListQuestion: %s", err)
		return nil, err
	}
	defer likesRows.Close()

	var list []*Question

	for likesRows.Next() {
		var id int64
		var question string
		var likesCount int64

		if err := likesRows.Scan(&id, &likesCount); err != nil {
			log.Errorf("ListQuestion: %s", err)
			return nil, err
		}

		questionRows, err := db.Query("SELECT question FROM question WHERE id = ?;", id)
		if err != nil {
			log.Errorf("ListQuestion: %s", err)
			return nil, err
		}

		for questionRows.Next() {
			if err := questionRows.Scan(&question); err != nil {
				log.Errorf("ListQuestion: %s", err)
				return nil, err
			}
		}
		questionRows.Close()

		list = append(list, &Question{
			Id:         id,
			Question:   question,
			LikesCount: likesCount,
		})
	}

	return &QuestionList{
		QuestionList: list,
	}, nil
}

func (b *Board) Like(ctx context.Context, likes *Likes) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if err := insertLikes(db, likes.UserId, likes.QuestionId); err != nil {
		log.Errorf("Like: %s", err)
		return nil, err
	}

	return nil, nil
}

func (b *Board) Unlike(ctx context.Context, likes *Likes) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if err := deleteLikes(db, likes.UserId, likes.QuestionId); err != nil {
		log.Errorf("Unlike: %s", err)
		return nil, err
	}

	return nil, nil
}

func insertSubject(db *sql.DB, title string) error {
	stmt, err := db.Prepare("INSERT INTO subject(title) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(title)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func selectSubjectByTitle(db *sql.DB, title string) (*Subject, error) {
	rows, err := db.Query("SELECT id, title FROM subject WHERE title = '?'", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subject := &Subject{}

	for rows.Next() {
		if err := rows.Scan(&subject.Id, &subject.Title); err != nil {
			return nil, err
		}
	}

	return subject, nil
}

func deleteSubject(db *sql.DB, subjectId int64) error {
	stmt, err := db.Prepare("DELETE FROM subject WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(subjectId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func insertQuestion(db *sql.DB, question string, subjectId int64) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO question(question, subject_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	res, err := tx.Stmt(stmt).Exec(question, subjectId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	insertedId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	tx.Commit()
	return insertedId, nil
}

func selectQuestion(db *sql.DB, id int64) (*Question, error) {
	rows, err := db.Query("SELECT id, question FROM question WHERE id = '?'", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	question := &Question{}

	for rows.Next() {
		if err := rows.Scan(&question.Id, &question.Question); err != nil {
			return nil, err
		}
	}

	rows, err = db.Query("SELECT COUNT(*) FROM likes WHERE question_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&question.LikesCount); err != nil {
			return nil, err
		}
	}

	return question, nil
}

func deleteQuestion(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("DELETE FROM question WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func insertLikes(db *sql.DB, userId string, questionId int64) error {
	stmt, err := db.Prepare("INSERT INTO likes(user_id, question_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(userId, questionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func deleteLikes(db *sql.DB, userId string, questionId int64) error {
	stmt, err := db.Prepare("DELETE FROM likes WHERE user_id = ? AND question_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(userId, questionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
