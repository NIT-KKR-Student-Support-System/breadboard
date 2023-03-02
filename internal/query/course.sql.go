// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: course.sql

package query

import (
	"context"
	"encoding/json"

	"github.com/lib/pq"
)

const createCourse = `-- name: CreateCourse :exec
INSERT INTO course (
    code,
    title,
    prereq,
    kind,
    objectives,
    content,
    book_names,
    outcomes
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (code)
DO
    UPDATE SET
        title = $2,
        prereq = $3,
        kind = $4,
        objectives = $5,
        content = $6,
        book_names = $7,
        outcomes = $8
`

type CreateCourseParams struct {
	Code       string   `json:"code"`
	Title      string   `json:"title"`
	Prereq     []string `json:"prereq"`
	Kind       string   `json:"kind"`
	Objectives []string `json:"objectives"`
	Content    []string `json:"content"`
	BookNames  []string `json:"book_names"`
	Outcomes   []string `json:"outcomes"`
}

func (q *Queries) CreateCourse(ctx context.Context, arg CreateCourseParams) error {
	_, err := q.db.ExecContext(ctx, createCourse,
		arg.Code,
		arg.Title,
		pq.Array(arg.Prereq),
		arg.Kind,
		pq.Array(arg.Objectives),
		pq.Array(arg.Content),
		pq.Array(arg.BookNames),
		pq.Array(arg.Outcomes),
	)
	return err
}

const createSpecifics = `-- name: CreateSpecifics :exec
INSERT INTO branch_specifics (
    code,
    branch,
    semester,
    credits
)
VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (code, branch)
DO
    UPDATE SET
        semester = $3,
        credits = $4
`

type CreateSpecificsParams struct {
	Code     string  `json:"code"`
	Branch   string  `json:"branch"`
	Semester int16   `json:"semester"`
	Credits  []int16 `json:"credits"`
}

func (q *Queries) CreateSpecifics(ctx context.Context, arg CreateSpecificsParams) error {
	_, err := q.db.ExecContext(ctx, createSpecifics,
		arg.Code,
		arg.Branch,
		arg.Semester,
		pq.Array(arg.Credits),
	)
	return err
}

const getCourse = `-- name: GetCourse :one
SELECT
    c.code, c.title, c.prereq, c.kind, c.objectives, c.content, c.book_names, c.outcomes, (
        SELECT
            COALESCE(JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            )), '[]')::JSON
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code = $1
`

type GetCourseRow struct {
	Code       string          `json:"code"`
	Title      string          `json:"title"`
	Prereq     []string        `json:"prereq"`
	Kind       string          `json:"kind"`
	Objectives []string        `json:"objectives"`
	Content    []string        `json:"content"`
	BookNames  []string        `json:"book_names"`
	Outcomes   []string        `json:"outcomes"`
	Specifics  json.RawMessage `json:"specifics"`
}

func (q *Queries) GetCourse(ctx context.Context, code string) (GetCourseRow, error) {
	row := q.db.QueryRowContext(ctx, getCourse, code)
	var i GetCourseRow
	err := row.Scan(
		&i.Code,
		&i.Title,
		pq.Array(&i.Prereq),
		&i.Kind,
		pq.Array(&i.Objectives),
		pq.Array(&i.Content),
		pq.Array(&i.BookNames),
		pq.Array(&i.Outcomes),
		&i.Specifics,
	)
	return i, err
}

const getCourses = `-- name: GetCourses :many
SELECT
    c.code, c.title, c.prereq, c.kind, c.objectives, c.content, c.book_names, c.outcomes, (
        SELECT
            COALESCE(JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            )), '[]')::JSON
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
`

type GetCoursesRow struct {
	Code       string          `json:"code"`
	Title      string          `json:"title"`
	Prereq     []string        `json:"prereq"`
	Kind       string          `json:"kind"`
	Objectives []string        `json:"objectives"`
	Content    []string        `json:"content"`
	BookNames  []string        `json:"book_names"`
	Outcomes   []string        `json:"outcomes"`
	Specifics  json.RawMessage `json:"specifics"`
}

func (q *Queries) GetCourses(ctx context.Context) ([]GetCoursesRow, error) {
	rows, err := q.db.QueryContext(ctx, getCourses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCoursesRow
	for rows.Next() {
		var i GetCoursesRow
		if err := rows.Scan(
			&i.Code,
			&i.Title,
			pq.Array(&i.Prereq),
			&i.Kind,
			pq.Array(&i.Objectives),
			pq.Array(&i.Content),
			pq.Array(&i.BookNames),
			pq.Array(&i.Outcomes),
			&i.Specifics,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCoursesByBranch = `-- name: GetCoursesByBranch :many
SELECT
    c.code, c.title, c.prereq, c.kind, c.objectives, c.content, c.book_names, c.outcomes, (
        SELECT
            JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            ))
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code IN (SELECT DISTINCT bs.code FROM branch_specifics AS bs WHERE bs.branch = $1)
`

type GetCoursesByBranchRow struct {
	Code       string          `json:"code"`
	Title      string          `json:"title"`
	Prereq     []string        `json:"prereq"`
	Kind       string          `json:"kind"`
	Objectives []string        `json:"objectives"`
	Content    []string        `json:"content"`
	BookNames  []string        `json:"book_names"`
	Outcomes   []string        `json:"outcomes"`
	Specifics  json.RawMessage `json:"specifics"`
}

func (q *Queries) GetCoursesByBranch(ctx context.Context, branch string) ([]GetCoursesByBranchRow, error) {
	rows, err := q.db.QueryContext(ctx, getCoursesByBranch, branch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCoursesByBranchRow
	for rows.Next() {
		var i GetCoursesByBranchRow
		if err := rows.Scan(
			&i.Code,
			&i.Title,
			pq.Array(&i.Prereq),
			&i.Kind,
			pq.Array(&i.Objectives),
			pq.Array(&i.Content),
			pq.Array(&i.BookNames),
			pq.Array(&i.Outcomes),
			&i.Specifics,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCoursesByBranchAndSemester = `-- name: GetCoursesByBranchAndSemester :many
SELECT
    c.code, c.title, c.prereq, c.kind, c.objectives, c.content, c.book_names, c.outcomes, (
        SELECT
            JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            ))
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code IN (SELECT DISTINCT bs.code FROM branch_specifics AS bs WHERE bs.branch = $1 AND bs.semester = $2)
`

type GetCoursesByBranchAndSemesterParams struct {
	Branch   string `json:"branch"`
	Semester int16  `json:"semester"`
}

type GetCoursesByBranchAndSemesterRow struct {
	Code       string          `json:"code"`
	Title      string          `json:"title"`
	Prereq     []string        `json:"prereq"`
	Kind       string          `json:"kind"`
	Objectives []string        `json:"objectives"`
	Content    []string        `json:"content"`
	BookNames  []string        `json:"book_names"`
	Outcomes   []string        `json:"outcomes"`
	Specifics  json.RawMessage `json:"specifics"`
}

func (q *Queries) GetCoursesByBranchAndSemester(ctx context.Context, arg GetCoursesByBranchAndSemesterParams) ([]GetCoursesByBranchAndSemesterRow, error) {
	rows, err := q.db.QueryContext(ctx, getCoursesByBranchAndSemester, arg.Branch, arg.Semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCoursesByBranchAndSemesterRow
	for rows.Next() {
		var i GetCoursesByBranchAndSemesterRow
		if err := rows.Scan(
			&i.Code,
			&i.Title,
			pq.Array(&i.Prereq),
			&i.Kind,
			pq.Array(&i.Objectives),
			pq.Array(&i.Content),
			pq.Array(&i.BookNames),
			pq.Array(&i.Outcomes),
			&i.Specifics,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCoursesBySemester = `-- name: GetCoursesBySemester :many
SELECT
    c.code, c.title, c.prereq, c.kind, c.objectives, c.content, c.book_names, c.outcomes, (
        SELECT
            JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            ))
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code IN (SELECT DISTINCT bs.code FROM branch_specifics AS bs WHERE bs.semester = $1)
`

type GetCoursesBySemesterRow struct {
	Code       string          `json:"code"`
	Title      string          `json:"title"`
	Prereq     []string        `json:"prereq"`
	Kind       string          `json:"kind"`
	Objectives []string        `json:"objectives"`
	Content    []string        `json:"content"`
	BookNames  []string        `json:"book_names"`
	Outcomes   []string        `json:"outcomes"`
	Specifics  json.RawMessage `json:"specifics"`
}

func (q *Queries) GetCoursesBySemester(ctx context.Context, semester int16) ([]GetCoursesBySemesterRow, error) {
	rows, err := q.db.QueryContext(ctx, getCoursesBySemester, semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCoursesBySemesterRow
	for rows.Next() {
		var i GetCoursesBySemesterRow
		if err := rows.Scan(
			&i.Code,
			&i.Title,
			pq.Array(&i.Prereq),
			&i.Kind,
			pq.Array(&i.Objectives),
			pq.Array(&i.Content),
			pq.Array(&i.BookNames),
			pq.Array(&i.Outcomes),
			&i.Specifics,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
