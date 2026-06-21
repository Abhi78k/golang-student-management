package cache

import (
	"context"
	"encoding/json"
	"log"
	"project-1/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type StudentCache struct {
	client *redis.Client
}

func NewStudentCache(
	client *redis.Client,
) *StudentCache {

	return &StudentCache{
		client: client,
	}
}

func (c *StudentCache) GetStudent(
	ctx context.Context,
	id string,
) (*models.Student, error) {

	key := "student:" + id

	value, err := c.client.Get(
		ctx,
		key,
	).Result()

	if err == redis.Nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var student models.Student

	err = json.Unmarshal(
		[]byte(value),
		&student,
	)

	if err != nil {
		return nil, err
	}

	log.Println("CACHE HITTT")

	return &student, nil
}

func (c *StudentCache) SetStudent(
	ctx context.Context,
	student *models.Student,
) error {

	key := "student:" + student.ID

	data, err := json.Marshal(student)

	if err != nil {
		return err
	}

	return c.client.Set(
		ctx,
		key,
		data,
		5*time.Minute,
	).Err()
}

func (c *StudentCache) DeleteStudent(
	ctx context.Context,
	id string,
) error {

	key := "student:" + id

	return c.client.Del(
		ctx,
		key,
	).Err()
}
