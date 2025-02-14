package handlers

import (
	"context"
	"task-manager-backend/internal/database"

	"github.com/gofiber/fiber/v2"
)

type Task struct {
	ID     int    `json: "id"`
	Title  string `json: "title"`
	Status string `json: "status"`
}

func GetTasks(c *fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT id, title, status FROM tasks")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Title, &t.Status)
		tasks = append(tasks, t)
	}

	return c.JSON(tasks)
}
