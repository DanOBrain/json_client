package main

import (
	"flag"
	"fmt"
	"json-client/client"
	"json-client/format"
	"json-client/models"
	"os"
)

func main() {
	// Флаги
	user := flag.Int("user", 1, "ID пользователя")
	limit := flag.Int("limit", 3, "Сколько постов показать")
	help := flag.Bool("help", false, "Помощь")
	flag.Parse()

	if *help {
		fmt.Println("Программа для постов")
		fmt.Println("\nФлаги:")
		fmt.Println("  -user  ID пользователя (1-10)")
		fmt.Println("  -limit сколько постов показать")
		fmt.Println("\nПример: go run main.go -user=2 -limit=5")
		return
	}

	// Создание клиента и обработчика вывода
	api := client.NewClient()
	out := format.NewFormatter()

	// Этап 1: Получение постов
	posts, err := api.GetPosts(*user)
	if err != nil {
		out.PrintError(err)
		os.Exit(1)
	}

	if len(posts) == 0 {
		fmt.Println("Нет постов")
		return
	}

	// ограничение
	if *limit > 0 && *limit < len(posts) {
		posts = posts[:*limit]
	}

	fmt.Printf("Загружено %d постов\n", len(posts))

	// Вывод постов
	fmt.Println("\n\n=== ПОСТЫ ===")
	for _, post := range posts {
		out.PrintPost(post)
	}

	// Этап 2 - Получение комментариев (Fan-Out)
	fmt.Println("\n\n=== КОММЕНТАРИИ ===")
	// Структура для результата
	type result struct {
		comments []models.Comment
		err      error
	}

	// Каналы многопоточности
	ch := make(chan result, len(posts))

	// Процесс получения комментариев каждого поста
	for _, p := range posts {
		go func(p models.Post) {
			// Получение комментариев
			comments, err := api.GetComments(p.ID)
			if err != nil {
				ch <- result{err: err}
				return
			}

			ch <- result{comments: comments}
		}(p)
	}

	// Получение результатов
	for _, p := range posts {
		res := <-ch
		if res.err != nil {
			out.PrintError(res.err)
			continue
		}
		// Вывод комментариев
		fmt.Printf("\nКомментарии к посту %d (%s):\n", p.ID, p.Title)
		for _, c := range res.comments {
			out.PrintComment(c)
		}
	}
	close(ch)
}
