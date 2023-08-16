package grpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	// external packages
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Board struct {
	BoardServer
}

func (b *Board) ListSubjects(ctx context.Context, empty *emptypb.Empty) (*SubjectList, error) {
	tx := sentry.TransactionFromContext(ctx)
	span := tx.StartChild("/board.Board/ListSubjects")
	span.Status = sentry.SpanStatusOK
	defer span.Finish()

	db := ctx.Value(DBSession).(*sql.DB)

	rows, err := db.QueryContext(ctx, "SELECT id, title, enabled FROM subject ORDER BY id;")
	if err != nil {
		log.Errorf("ListSubjects: %s", err)
		return nil, err
	}
	defer rows.Close()

	var list []*Subject

	for rows.Next() {
		var id int64
		var title string
		var enabled bool

		if err := rows.Scan(&id, &title, &enabled); err != nil {
			log.Errorf("ListSubjects: %s", err)
			return nil, err
		}

		list = append(list, &Subject{
			Id:      id,
			Title:   title,
			Enabled: enabled,
		})
	}

	return &SubjectList{
		SubjectList: list,
	}, nil
}

func (b *Board) GetSubject(ctx context.Context, subjectId *SubjectId) (*Subject, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	rows, err := db.Query("SELECT id, title, enabled FROM subject WHERE id = $1", subjectId.Id)
	if err != nil {
		log.Errorf("GetSubject: %s", err)
		return nil, err
	}
	defer rows.Close()

	subject := &Subject{}

	for rows.Next() {
		if err := rows.Scan(&subject.Id, &subject.Title, &subject.Enabled); err != nil {
			log.Errorf("GetSubject: %s", err)
			return nil, err
		}
	}

	return subject, nil
}

func (b *Board) CreateQuestion(ctx context.Context, newQuestion *NewQuestion) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	subject, err := selectSubject(db, newQuestion.SubjectId)
	if err != nil {
		log.Errorf("CreateQuestion: failed to select subject. %s", err)
		return nil, err
	}

	if subject == nil {
		log.Errorf("CreateQuestion: subjectId '%d' is not exists", newQuestion.SubjectId)
		return nil, errors.New("this subject is not exists")
	}

	if subject.Enabled == false {
		log.Errorf("CreateQuestion: subjectId '%d' is disable", subject.Id)
		return nil, err
	}

	if len(newQuestion.GetQuestion()) == 0 {
		log.Errorf("CreateQuestion: empty input 'question'")
		return nil, err
	}

	err = insertQuestion(db, newQuestion.Question, newQuestion.SubjectId)
	if err != nil {
		log.Errorf("CreateQuestion: %s", err)

		return nil, err
	}

	return nil, nil
}

func (b *Board) ListQuestions(ctx context.Context, subjectId *SubjectId) (*QuestionList, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	rows, err := db.Query(
		"SELECT id, question, likes FROM question WHERE subject_id = $1 ORDER BY likes DESC, question ASC;",
		subjectId.Id)
	defer rows.Close()

	if err != nil {
		log.Errorf("ListQuestions: %s", err)
		return nil, err
	}

	var list []*Question

	for rows.Next() {
		var id int64
		var question string
		var likesCount int64

		if err := rows.Scan(&id, &question, &likesCount); err != nil {
			log.Errorf("ListQuestions: %s", err)
			return nil, err
		}

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

func (b *Board) Like(ctx context.Context, questionId *QuestionId) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	if err := addQuestionLikes(db, questionId.Id); err != nil {
		log.Errorf("Like: %s", err)
		return nil, err
	}

	return nil, nil
}

func (b *Board) Unlike(ctx context.Context, questionId *QuestionId) (*emptypb.Empty, error) {
	db := ctx.Value(DBSession).(*sql.DB)

	question, err := selectQuestion(db, questionId.Id)
	if err != nil {
		log.Errorf("Unlike: %s", err)
		return nil, err
	}
	if question.LikesCount <= 0 {
		err = fmt.Errorf("like count can not be negative")
		log.Errorf("Unlike: %s", err)
		return nil, err
	}

	if err := subQuestionLikes(db, questionId.Id); err != nil {
		log.Errorf("Unlike: %s", err)
		return nil, err
	}

	return nil, nil
}

func selectSubject(db *sql.DB, id int64) (*Subject, error) {
	rows, err := db.Query("SELECT id, title, enabled FROM subject WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subject := &Subject{}

	for rows.Next() {
		if err := rows.Scan(&subject.Id, &subject.Title, &subject.Enabled); err != nil {
			return nil, err
		}
	}

	return subject, nil
}

func insertQuestion(db *sql.DB, question string, subjectId int64) error {
	stmt, err := db.Prepare("INSERT INTO question(question, subject_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(question, subjectId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func selectQuestion(db *sql.DB, id int64) (*Question, error) {
	rows, err := db.Query("SELECT id, question, likes FROM question WHERE id = $1", id)
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

func addQuestionLikes(db *sql.DB, questionId int64) error {
	stmt, err := db.Prepare("UPDATE question SET likes = likes + 1 WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(questionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func subQuestionLikes(db *sql.DB, questionId int64) error {
	stmt, err := db.Prepare("UPDATE question SET likes = likes - 1 WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(questionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
