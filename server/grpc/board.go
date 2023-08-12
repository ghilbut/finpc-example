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

	rows, err := db.Query("SELECT id, title, enabled FROM subject ORDER BY id;")
	if err != nil {
		log.Errorf("ListSubject: %s", err)
		return nil, err
	}
	defer rows.Close()

	var list []*Subject

	for rows.Next() {
		var id int64
		var title string
		var enabled bool

		if err := rows.Scan(&id, &title, &enabled); err != nil {
			log.Errorf("ListSubject: %s", err)
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
		return nil, errors.New("this subject is disable")
	}

	if len(newQuestion.GetQuestion()) == 0 {
		log.Errorf("CreateQuestion: empty input 'question'")
		return nil, errors.New("'question' is empty")
	}

	err = insertQuestion(db, newQuestion.Question, newQuestion.SubjectId)
	if err != nil {
		log.Errorf("CreateQuestion: %s", err)
		return nil, err
	}

	return nil, nil
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

	rows, err := db.Query(
		"SELECT id, question, likes FROM question WHERE subject_id = $1 ORDER BY likes DESC;",
		subjectId.Id)
	if err != nil {
		log.Errorf("ListQuestion: %s", err)
		return nil, err
	}
	defer rows.Close()

	var list []*Question

	for rows.Next() {
		var id int64
		var question string
		var likesCount int64

		if err := rows.Scan(&id, &question, &likesCount); err != nil {
			log.Errorf("ListQuestion: %s", err)
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

	if err := subQuestionLikes(db, questionId.Id); err != nil {
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
	rows, err := db.Query("SELECT id, title, enabled FROM subject WHERE title = $1", title)
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

func deleteSubject(db *sql.DB, subjectId int64) error {
	stmt, err := db.Prepare("DELETE FROM subject WHERE id = $1")
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

func deleteQuestion(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("DELETE FROM question WHERE id = $1")
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

func addQuestionLikes(db *sql.DB, questionId int64) error {
	//stmt, err := db.Prepare("INSERT INTO likes(user_id, question_id) VALUES (?, ?)")
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
	//stmt, err := db.Prepare("DELETE FROM likes WHERE user_id = ? AND question_id = ?")
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
