package models

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type Comment struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	ProjectID  uuid.UUID  `db:"project_id" json:"project_id"`
	IdentityID uuid.UUID  `db:"identity_id" json:"identity_id"`
	MediaID    *uuid.UUID `db:"media_id" json:"media_id"`
	ParentID   *uuid.UUID `db:"parent_id" json:"parent_id"`
	Content    string     `db:"content" json:"content"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`

	Identity *Identity `db:"-" json:"identity"`
	Media    *Media    `db:"-" json:"media"`

	IdentityJson types.JSONText `db:"identity" json:"-"`
	MediaJson    types.JSONText `db:"media" json:"-"`

	ChildrenCount int64     `db:"children_count" json:"children_count"`
	Children      []Comment `db:"-" json:"children"`

	Likes         int64 `db:"likes" json:"likes"`
	IdentityLiked bool  `db:"identity_liked" json:"identity_liked"`

	Reactions        types.JSONText `db:"reactions" json:"reactions"`
	IdentityReaction *string        `db:"identity_reaction" json:"identity_reaction"`
}

type Like struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CommentID  uuid.UUID `db:"comment_id" json:"comment_id"`
	IdentityID uuid.UUID `db:"identity_id" json:"identity_id"`
	LikedAt    time.Time `db:"created_at" json:"created_at"`

	Comment  *Comment  `db:"-" json:"comment"`
	Identity *Identity `db:"-" json:"identity"`
}

type Reaction struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CommentID  uuid.UUID `db:"comment_id" json:"comment_id"`
	IdentityID uuid.UUID `db:"identity_id" json:"identity_id"`
	Reaction   string    `db:"reaction" json:"reaction"`
	ReactedAt  time.Time `db:"created_at" json:"created_at"`

	Comment  *Comment  `db:"-" json:"comment"`
	Identity *Identity `db:"-" json:"identity"`
}

func (Comment) TableName() string {
	return "comments"
}

func (Comment) FetchQuery() string {
	return "comments/fetch"
}

func (Like) TableName() string {
	return "comment_likes"
}

func (Like) FetchQuery() string {
	return "comments/fetch_likes"
}

func (Reaction) TableName() string {
	return "comment_reactions"
}

func (Reaction) FetchQuery() string {
	return "comments/fetch_reactions"
}

func (c *Comment) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"comments/create",
		c.ProjectID,
		c.IdentityID,
		c.MediaID,
		c.ParentID,
		c.Content,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(c); err != nil {
			return err
		}
	}
	c, err = GetComment(c.ID, c.IdentityID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "comments/delete", c.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (c *Comment) Update(ctx context.Context) error {

	rows, err := database.Query(
		ctx,
		"comments/update",
		c.ID,
		c.Content,
		c.MediaID,
	)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(c); err != nil {
			return err
		}
	}
	c, err = GetComment(c.ID, c.IdentityID)
	if err != nil {
		return err
	}
	return nil
}

func (l *Like) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"comments/like",
		l.CommentID,
		l.IdentityID,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(l); err != nil {
			return err
		}
	}
	return nil
}

func (l *Like) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "comments/unlike", l.CommentID, l.IdentityID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *Reaction) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"comments/reaction",
		r.CommentID,
		r.IdentityID,
		r.Reaction,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(r); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reaction) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "comments/remove_reaction", r.CommentID, r.IdentityID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetComments(projectID, identityID uuid.UUID, p database.Paginate) ([]Comment, int, error) {
	var (
		comments  = []Comment{}
		fetchList []database.FetchList
	)

	if err := database.QuerySelect("comments/get_list", &fetchList, projectID, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	rows, err := database.Queryx("comments/get", projectID, identityID, p.Limit, p.Offet)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		c := new(Comment)
		if err := rows.StructScan(c); err != nil {
			return nil, 0, err
		}
		if err := database.UnmarshalJSONTextFields(c); err != nil {
			log.Println(err)
		}
		comments = append(comments, *c)
	}

	return comments, fetchList[0].TotalCount, nil
}

func GetCommentChildren(parentID, identityID uuid.UUID) ([]Comment, error) {
	comments := []Comment{}

	rows, err := database.Queryx("comments/get_children", parentID, identityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := new(Comment)
		if err := rows.StructScan(c); err != nil {
			return nil, err
		}
		if err := database.UnmarshalJSONTextFields(c); err != nil {
			log.Println(err)
		}
		comments = append(comments, *c)
	}

	return comments, nil
}

func GetComment(id, identityID uuid.UUID) (*Comment, error) {
	c := new(Comment)
	if err := database.Get(c, "comments/get_by_id", id, identityID); err != nil {
		return nil, err
	}
	if c.ChildrenCount > 0 {
		if children, err := GetCommentChildren(c.ID, identityID); err == nil {
			c.Children = children
		}
	}
	return c, nil
}

func GetLike(commentID, identityID uuid.UUID) (*Like, error) {
	l := new(Like)
	if err := database.Get(l, "comments/get_like", commentID, identityID); err != nil {
		return nil, err
	}
	return l, nil
}

func GetReaction(commentID, identityID uuid.UUID) (*Reaction, error) {
	r := new(Reaction)
	if err := database.Get(r, "comments/get_reaction", commentID, identityID); err != nil {
		return nil, err
	}
	return r, nil
}
