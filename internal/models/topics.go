package models

import (
	"database/sql"
	"fmt"
	queryHelpers "vulh/soundcommunity/internal/utils"
)

type Topic struct {
	ID        int            `json:"id"`
	Slug      string         `json:"slug"`
	TopicName string         `json:"topic_name"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt sql.NullString `json:"-"`
}

type TopicWithTotalPost struct {
	ID         int    `json:"id"`
	Slug       string `json:"slug"`
	TopicName  string `json:"topic_name"`
	TotalPosts int    `json:"total_posts"`
	CreatedAt  string `json:"created_at"`
}

type CreateTopicInput struct {
	TopicName string `json:"topic_name"`
}

type TopicModel struct {
	DB *sql.DB
}

func (m *TopicModel) GetAllTopics() ([]Topic, error) {
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	var topics []Topic
	rows, err := queryBuilder.
		Select("*").
		FromTable("topics").
		OrderBy("created_at", queryHelpers.QueryDirection.DESC).
		GetMany()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var topic Topic
		err = rows.Scan(&topic.ID, &topic.Slug, &topic.TopicName, &topic.CreatedAt, &topic.UpdatedAt)
		if err == nil {
			topics = append(topics, topic)
		}
	}
	return topics, nil
}

func (m *TopicModel) FindTopicsBy(queryConfig *queryHelpers.QueryConfig) ([]Topic, error) {
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	rows, err := queryBuilder.
		Select("*").
		FromTable("topics").
		WhereColumn(queryConfig.WhereColumnName).
		Search(queryConfig.Operator, queryConfig.SearchValue).
		GetMany()
	if err != nil {
		return nil, err
	}
	var topics []Topic
	for rows.Next() {
		var topic Topic
		err = rows.Scan(&topic.ID, &topic.Slug, &topic.TopicName, &topic.CreatedAt, &topic.UpdatedAt)
		if err == nil {
			topics = append(topics, topic)
		}
	}
	return topics, nil
}

func (m *TopicModel) SearchTopicsByName(q string) ([]TopicWithTotalPost, error) {
	rows, err := m.DB.Query(`SELECT topics.id, topics.slug, topics.topic_name , COUNT(posts_topics.topic_id) AS total_posts, topics.created_at FROM topics LEFT JOIN posts_topics
							ON topics.id = posts_topics.topic_id
							WHERE topics.topic_name ILIKE $1 OR REPLACE(topics.slug, '-', ' ') ILIKE $1
							group by topics.id 
							ORDER BY total_posts DESC`, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	var topics []TopicWithTotalPost
	for rows.Next() {
		var topic TopicWithTotalPost
		err = rows.Scan(&topic.ID, &topic.Slug, &topic.TopicName, &topic.TotalPosts, &topic.CreatedAt)
		if err == nil {
			topics = append(topics, topic)
		}
	}
	return topics, nil
}

func (m *TopicModel) GetRelativeTopics(slug string) ([]Topic, error) {
	var topicId int
	row := m.DB.QueryRow("SELECT id FROM topics WHERE slug=$1", slug)
	err := row.Scan(&topicId)
	if err != nil {
		return nil, err
	}
	var topics []Topic
	rows, err := m.DB.Query(`SELECT post_id FROM posts_topics WHERE posts_topics.topic_id = $1`, topicId)
	if err != nil {
		return nil, err
	}
	var postIds []int
	for rows.Next() {
		var postId int
		err = rows.Scan(&postId)
		if err == nil {
			postIds = append(postIds, postId)
		}
	}
	if len(postIds) > 0 {
		inPostIdsQry := "IN ("
		for i, postId := range postIds {
			if i < len(postIds)-1 {
				inPostIdsQry += fmt.Sprintf(`'%v',`, postId)
			} else {
				inPostIdsQry += fmt.Sprintf(`'%v'`, postId) + ")"
			}
		}
		query := fmt.Sprintf(`SELECT topic_id FROM posts_topics 
							WHERE post_id %v 
							group by topic_id 
							order by COUNT(topic_id) DESC 
							LIMIT 10`, inPostIdsQry)
		rows, err := m.DB.Query(query)
		if err != nil {
			return nil, err
		}

		var topicIds []int
		for rows.Next() {
			var topicId int
			err = rows.Scan(&topicId)
			if err == nil {
				topicIds = append(topicIds, topicId)
			}
		}
		if len(topicIds) > 0 {
			for _, topicId := range topicIds {
				rows, err := m.DB.Query("SELECT * FROM topics WHERE id = $1", topicId)
				if err != nil {
					return make([]Topic, 0), nil
				}
				for rows.Next() {
					var topic Topic
					err = rows.Scan(&topic.ID, &topic.Slug, &topic.TopicName, &topic.CreatedAt, &topic.UpdatedAt)
					if err == nil {
						topics = append(topics, topic)
					}
				}
			}
		}
	}
	return topics, nil
}

func (m *TopicModel) GetRecommendedTopics() ([]Topic, error) {
	rows, err := m.DB.Query(`SELECT topic_id 
							FROM posts_topics
							group by topic_id
							order by COUNT(topic_id) desc
							limit 8`)
	if err != nil {
		return nil, err
	}
	var topics []Topic
	for rows.Next() {
		var topicId int
		err = rows.Scan(&topicId)
		if err == nil {
			var topic Topic
			row := m.DB.QueryRow("SELECT id, slug, topic_name FROM topics WHERE id=$1", topicId)
			err = row.Scan(&topic.ID, &topic.Slug, &topic.TopicName)
			if err == nil {
				topics = append(topics, topic)
			}
		}
	}
	return topics, nil
}
