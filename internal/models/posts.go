package models

import (
	"cmp"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	queryHelpers "vulh/soundcommunity/internal/utils"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Post struct {
	ID           int            `json:"id"`
	UserId       int            `json:"user_id"`
	Slug         string         `json:"slug"`
	Title        string         `json:"title"`
	ShortContent string         `json:"short_content"`
	ThumbnailUrl string         `json:"thumbnail_url"`
	Content      string         `json:"content"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    sql.NullString `json:"-"`
}

type PostDetails struct {
	ID           int            `json:"id"`
	UserId       int            `json:"user_id"`
	Slug         string         `json:"slug"`
	Title        string         `json:"title"`
	ShortContent string         `json:"short_content"`
	ThumbnailUrl string         `json:"thumbnail_url"`
	Content      string         `json:"content"`
	Topics       []Topic        `json:"topics"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    sql.NullString `json:"-"`
}

type PostWithUserInfo struct {
	ID           int            `json:"id"`
	UserId       int            `json:"user_id"`
	Username     string         `json:"username"`
	Slug         string         `json:"slug"`
	Title        string         `json:"title"`
	ShortContent string         `json:"short_content"`
	ThumbnailUrl string         `json:"thumbnail_url"`
	Content      string         `json:"content"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    sql.NullString `json:"-"`
}

type CreatePostTopic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreatePostInput struct {
	Title        string            `json:"title"`
	ShortContent string            `json:"short_content"`
	ThumbnailUrl string            `json:"thumbnail_url"`
	Content      string            `json:"content"`
	Topics       []CreatePostTopic `json:"topics"`
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(createPostInput *CreatePostInput, userId int) error {
	stmt := `INSERT INTO posts(user_id, slug, title, short_content, thumbnail_url, content) VALUES($1, $2, $3, $4, $5, $6)`
	uuid, _ := gonanoid.New()
	slug := removePunctuation(strings.ToLower(fmt.Sprintf("%v-%v", strings.ReplaceAll(createPostInput.Title, " ", "-"),
		strings.ReplaceAll(
			strings.ReplaceAll(uuid, "-", ""), "_", ""))))
	_, err := m.DB.Exec(stmt,
		userId,
		slug,
		createPostInput.Title,
		createPostInput.ShortContent,
		createPostInput.ThumbnailUrl,
		createPostInput.Content)

	if err != nil {
		return err
	}

	row := m.DB.QueryRow("SELECT id from posts WHERE slug = $1", slug)

	var postId int

	err = row.Scan(&postId)

	if err != nil {
		return err
	}

	for _, topic := range createPostInput.Topics {
		if topic.ID != -1 {
			stmt := `INSERT INTO posts_topics(post_id, topic_id) VALUES($1, $2)`
			m.DB.Exec(stmt, postId, topic.ID)
		} else {
			topicSlug := removePunctuation(strings.ToLower(strings.ReplaceAll(topic.Name, " ", "-")))
			stmt := `INSERT INTO topics(slug, topic_name) VALUES($1, $2)`
			_, err = m.DB.Exec(stmt, topicSlug, topic.Name)
			if err == nil {
				row := m.DB.QueryRow("SELECT id from topics WHERE slug = $1", topicSlug)

				var topicId int
				err = row.Scan(&topicId)

				if err == nil {
					stmt := `INSERT INTO posts_topics(post_id, topic_id) VALUES($1, $2)`
					m.DB.Exec(stmt, postId, topicId)
				}
			}
		}
	}

	return nil
}

func (m *PostModel) GetPostBySlug(slug string) (*PostDetails, error) {
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	post := &PostDetails{}
	row := queryBuilder.
		Select("id", "user_id", "slug", "title", "content", "created_at", "updated_at").
		FromTable("posts").
		WhereColumn("slug").
		Equal(slug).
		GetOne()
	err := row.Scan(&post.ID, &post.UserId, &post.Slug, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	var topics []Topic
	rows, err := m.DB.Query(`SELECT 
							topics.id, topics.slug, topics.topic_name, topics.created_at, topics.updated_at
							FROM posts_topics LEFT JOIN topics 
							ON posts_topics.topic_id = topics.id WHERE post_id = $1`, post.ID)

	if err != nil {
		post.Topics = make([]Topic, 0)
	}

	for rows.Next() {
		var topic Topic
		err = rows.Scan(&topic.ID, &topic.Slug, &topic.TopicName, &topic.CreatedAt, &topic.UpdatedAt)
		if err == nil {
			topics = append(topics, topic)
		}
	}

	post.Topics = topics
	return post, nil
}

func (m *PostModel) GetPosts(q string, page int, take int, start int, end int) ([]PostWithUserInfo, bool, bool, error) {
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	rows := &sql.Rows{}
	var err error
	limit := take
	if q != "*" {
		queryString := ""
		searchParams := "'%" + q + "%'"
		if end == -1 && start == -1 {
			queryString = fmt.Sprintf(`SELECT id, user_id, slug, title, short_content, thumbnail_url, created_at 
										FROM posts 
										WHERE title ILIKE %v OR REPLACE(slug, '-', ' ') ILIKE %v 
										ORDER BY created_at DESC 
										LIMIT %v`, searchParams, searchParams, limit)
		} else {
			if start != -1 {
				queryString = fmt.Sprintf(`SELECT id, user_id, slug, title, short_content, thumbnail_url, created_at 
											FROM posts 
											WHERE (title ILIKE %v OR REPLACE(slug, '-', ' ') ILIKE %v) AND created_at > '%v'
											ORDER BY created_at ASC 
											LIMIT %v`, searchParams, searchParams, start, limit)
			}
			if end != -1 {
				queryString = fmt.Sprintf(`SELECT id, user_id, slug, title, short_content, thumbnail_url, created_at 
											FROM posts 
											WHERE (title ILIKE %v OR REPLACE(slug, '-', ' ') ILIKE %v) AND created_at < '%v'
											ORDER BY created_at DESC 
											LIMIT %v`, searchParams, searchParams, end, limit)
			}
		}
		rows, err = m.DB.Query(queryString)
		if err != nil {
			return nil, false, false, err
		}
	} else {
		queryString := ""
		if end == -1 && start == -1 {
			queryString = fmt.Sprintf(`SELECT id, user_id, slug, title, short_content, thumbnail_url, created_at 
										FROM posts ORDER BY created_at DESC LIMIT %v`, limit)
		} else {
			if start != -1 {
				queryString = fmt.Sprintf(`SELECT id, user_id, slug, title, short_content, thumbnail_url, created_at 
											FROM posts 
											WHERE id > %v
											ORDER BY id ASC 
											LIMIT %v
										`, start, limit)
			}

			if end != -1 {
				queryString = fmt.Sprintf(`SELECT id, user_id, slug, title, short_content, thumbnail_url, created_at 
											FROM posts 
											WHERE id < %v
											ORDER BY id DESC 
											LIMIT %v 
										`, end, limit)
			}
		}

		rows, err = m.DB.Query(queryString)
		if err != nil {
			return nil, false, false, err
		}
	}

	var posts []PostWithUserInfo

	for rows.Next() {
		var post PostWithUserInfo
		err := rows.Scan(&post.ID,
			&post.UserId,
			&post.Slug,
			&post.Title,
			&post.ShortContent,
			&post.ThumbnailUrl,
			&post.CreatedAt)
		if err == nil {
			var username string
			row := queryBuilder.
				Select("username").
				FromTable("users").
				WhereColumn("id").
				Equal(fmt.Sprintf("%v", post.UserId)).
				GetOne()
			err = row.Scan(&username)
			if err == nil {
				post.Username = username
			}
			posts = append(posts, post)
		}
	}
	slices.SortFunc(posts, func(a, b PostWithUserInfo) int {
		return cmp.Compare(b.ID, a.ID)
	})
	hasNext := false
	hasPrev := false
	if len(posts) > 0 && len(posts) <= limit {
		if len(posts) == limit {
			queryStr := fmt.Sprintf(`SELECT id FROM posts
								WHERE id < %v
								ORDER BY id DESC LIMIT 1`, posts[len(posts)-1].ID)
			row := m.DB.QueryRow(queryStr)
			var nextId int
			err = row.Scan(&nextId)
			if err == nil {
				hasNext = true
			}
		}
		queryStr := fmt.Sprintf(`SELECT id FROM posts
								WHERE id > %v
								ORDER BY id DESC LIMIT 1`, posts[0].ID)
		row := m.DB.QueryRow(queryStr)
		var prevId int
		err = row.Scan(&prevId)
		if err == nil {
			hasPrev = true
		}
	}
	return posts, hasNext, hasPrev, nil
}

func (m *PostModel) GetPostsByTopic(topicSlug string, page int, take int) ([]PostWithUserInfo, int, bool, error) {
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	row := queryBuilder.
		Select("id").
		FromTable("topics").
		WhereColumn("slug").
		Equal(topicSlug).
		GetOne()
	var topicId int
	err := row.Scan(&topicId)
	if err != nil {
		return nil, 0, false, err
	}
	limit := take
	skip := page * limit
	rows, err := m.DB.Query(`SELECT posts.id, posts.user_id, posts.slug, posts.title, posts.short_content, posts.thumbnail_url, posts.created_at 
							FROM posts LEFT JOIN posts_topics
							ON posts.id = posts_topics.post_id
							WHERE posts_topics.topic_id = $1 OFFSET $2 LIMIT $3`, topicId, skip, limit)
	if err != nil {
		return nil, 0, false, err
	}
	var totalRows int
	row = m.DB.QueryRow(`SELECT COUNT(id) FROM posts_topics WHERE posts_topics.topic_id = $1`, topicId)
	errTotalRows := row.Scan(&totalRows)
	if errTotalRows != nil {
		totalRows = 0
	}
	var posts []PostWithUserInfo
	for rows.Next() {
		var post PostWithUserInfo
		err := rows.Scan(&post.ID,
			&post.UserId,
			&post.Slug,
			&post.Title,
			&post.ShortContent,
			&post.ThumbnailUrl,
			&post.CreatedAt)
		if err == nil {
			var username string
			row := queryBuilder.
				Select("username").
				FromTable("users").
				WhereColumn("id").
				Equal(fmt.Sprintf("%v", post.UserId)).
				GetOne()
			err = row.Scan(&username)
			if err == nil {
				post.Username = username
			}
			posts = append(posts, post)
		}
	}
	hasNext := (skip + take) < totalRows
	return posts, totalRows, hasNext, nil
}
