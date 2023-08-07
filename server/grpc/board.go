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

func (b *Board) PostSubject(ctx context.Context, newSubject *NewSubject) (*Subject, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if len(newSubject.GetTitle()) == 0 {
		log.Errorf("PostSubject: invalid input 'title'")
		return nil, errors.New("invalid 'title'")
	}

	if err := insertSubject(db, newSubject.Title); err != nil {
		log.Errorf("PostSubject: %s", err)
		return nil, err
	}

	subject, err := selectSubjectByTitle(db, newSubject.Title)
	if err != nil {
		log.Errorf("PostSubject: failed to select created subject. %s", err)
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

func (b *Board) GetSubjectList(ctx context.Context, empty *emptypb.Empty) (*SubjectList, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	rows, err := db.Query("SELECT id, title FROM subject ORDER BY id;")
	if err != nil {
		log.Errorf("GetSubjectList: %s", err)
		return nil, err
	}
	defer rows.Close()

	var list []*Subject

	for rows.Next() {
		var id int64
		var title string

		if err := rows.Scan(&id, title); err != nil {
			log.Fatalf("GetSubjectList: %s", err)
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

func (b *Board) PostQuestion(ctx context.Context, newQuestion *NewQuestion) (*Question, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if len(newQuestion.GetQuestion()) == 0 {
		log.Errorf("PostQuestion: empty input 'question'")
		return nil, errors.New("'question' is empty")
	}

	insertedId, err := insertQuestion(db, newQuestion.Question, newQuestion.SubjectId)
	if err != nil {
		log.Errorf("PostQuestion: %s", err)
		return nil, err
	}

	question, err := selectQuestion(db, insertedId)
	if err != nil {
		log.Errorf("PostQuestion: failed to select created question. %s", err)
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

func (b *Board) GetQuestionList(ctx context.Context, subjectId *SubjectId) (*QuestionList, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	rows, err := db.Query("SELECT id, question, likes FROM question WHERE subject_id = ? ORDER BY likes DESC;", subjectId)
	if err != nil {
		log.Errorf("GetQuestionList: %s", err)
		return nil, err
	}
	defer rows.Close()

	var list []*Question

	for rows.Next() {
		var id int64
		var question string
		var likes int64

		if err := rows.Scan(&id, &question, &likes); err != nil {
			log.Fatalf("GetQuestionList: %s", err)
		}

		list = append(list, &Question{
			Id:         id,
			Question:   question,
			LikesCount: likes,
		})
	}

	return &QuestionList{
		QuestionList: list,
	}, nil
}

func (b *Board) GetQuestion(ctx context.Context, questionId *QuestionId) (*Question, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	rows, err := db.Query("SELECT id, question, likes FROM question WHERE id = ?;", questionId)
	if err != nil {
		log.Errorf("GetQuestion: %s", err)
		return nil, err
	}
	defer rows.Close()

	question := &Question{}

	for rows.Next() {
		if err := rows.Scan(&question.Id, &question.Question, &question.LikesCount); err != nil {
			log.Fatalf("GetQuestion: %s", err)
		}
	}

	return question, nil
}

func (b *Board) LikeQuestion(ctx context.Context, questionId *QuestionId) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if err := addQuestionLikes(db, questionId.Id); err != nil {
		log.Errorf("LikeQuestion: %s", err)
		return nil, err
	}

	return nil, nil
}

func (b *Board) LikeQuestionCancel(ctx context.Context, questionId *QuestionId) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if err := subQuestionLikes(db, questionId.Id); err != nil {
		log.Errorf("LikeQuestionCancel: %s", err)
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
	rows, err := db.Query("SELECT id, question, likes FROM question WHERE id = '?'", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	question := &Question{}

	for rows.Next() {
		if err := rows.Scan(&question.Id, &question.Question, &question.LikesCount); err != nil {
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

func addQuestionLikes(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("UPDATE question SET likes = likes + 1 WHERE id = ?;")
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

func subQuestionLikes(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("UPDATE question SET likes = likes - 1 WHERE id = ?;")
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
