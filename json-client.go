package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type BaseResource struct {
	ID int `json:"id"`
}

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

type Album struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
}

type Photo struct {
	AlbumID      int    `json:"albumId"`
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}

type ResourceType string

const (
	ResourcePosts    ResourceType = "posts"
	ResourceComments ResourceType = "comments"
	ResourceAlbums   ResourceType = "albums"
	ResourcePhotos   ResourceType = "photos"
	ResourceTodos    ResourceType = "todos"
	ResourceUsers    ResourceType = "users"
)

type Resource interface {
	GetID() int
	GetType() ResourceType
	GetDisplayName() string
}

func (p Post) GetID() int             { return p.ID }
func (p Post) GetType() ResourceType  { return ResourcePosts }
func (p Post) GetDisplayName() string { return p.Title }

func (c Comment) GetID() int             { return c.ID }
func (c Comment) GetType() ResourceType  { return ResourceComments }
func (c Comment) GetDisplayName() string { return c.Name }

func (a Album) GetID() int             { return a.ID }
func (a Album) GetType() ResourceType  { return ResourceAlbums }
func (a Album) GetDisplayName() string { return a.Title }

func (p Photo) GetID() int             { return p.ID }
func (p Photo) GetType() ResourceType  { return ResourcePhotos }
func (p Photo) GetDisplayName() string { return p.Title }

func (t Todo) GetID() int             { return t.ID }
func (t Todo) GetType() ResourceType  { return ResourceTodos }
func (t Todo) GetDisplayName() string { return t.Title }

func (u User) GetID() int             { return u.ID }
func (u User) GetType() ResourceType  { return ResourceUsers }
func (u User) GetDisplayName() string { return u.Name }

type ProcessedResource struct {
	ID          int
	Type        ResourceType
	Original    interface{}
	Summary     string
	Details     map[string]interface{}
	ProcessedAt time.Time
}

type Config struct {
	Resource   ResourceType
	ID         int
	UserID     int
	Limit      int
	Expand     bool
	APIBaseURL string
	Timeout    time.Duration
}

type APIError struct {
	Message string
	Status  int
	Err     error
}

func (e *APIError) Error() string {
	if e.Status > 0 {
		return fmt.Sprintf("Ошибка API: %s (статус %d)", e.Message, e.Status)
	}
	return fmt.Sprintf("Ошибка API: %s", e.Message)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	client  HTTPClient
	baseURL string
}

func NewAPIClient(client HTTPClient, baseURL string) *APIClient {
	return &APIClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *APIClient) Fetch(resource ResourceType, id int, userID int) ([]interface{}, error) {
	var url string

	if id > 0 {
		url = fmt.Sprintf("%s/%s/%d", c.baseURL, resource, id)
	} else if userID > 0 {
		url = fmt.Sprintf("%s/%s?userId=%d", c.baseURL, resource, userID)
	} else {
		url = fmt.Sprintf("%s/%s", c.baseURL, resource)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, &APIError{Message: "не удалось создать запрос", Err: err}
	}

	req.Header.Set("User-Agent", "JSONPlaceholder-Client/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, &APIError{Message: "не удалось получить данные", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			Message: fmt.Sprintf("неожиданный статус код: %d", resp.StatusCode),
			Status:  resp.StatusCode,
		}
	}

	return c.decodeResponse(resp, resource, id)
}

func (c *APIClient) decodeResponse(resp *http.Response, resource ResourceType, id int) ([]interface{}, error) {
	var result []interface{}

	switch resource {
	case ResourcePosts:
		if id > 0 {
			var item Post
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				return nil, err
			}
			result = append(result, item)
		} else {
			var items []Post
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return nil, err
			}
			for _, item := range items {
				result = append(result, item)
			}
		}

	case ResourceComments:
		if id > 0 {
			var item Comment
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				return nil, err
			}
			result = append(result, item)
		} else {
			var items []Comment
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return nil, err
			}
			for _, item := range items {
				result = append(result, item)
			}
		}

	case ResourceAlbums:
		if id > 0 {
			var item Album
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				return nil, err
			}
			result = append(result, item)
		} else {
			var items []Album
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return nil, err
			}
			for _, item := range items {
				result = append(result, item)
			}
		}

	case ResourcePhotos:
		if id > 0 {
			var item Photo
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				return nil, err
			}
			result = append(result, item)
		} else {
			var items []Photo
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return nil, err
			}
			for _, item := range items {
				result = append(result, item)
			}
		}

	case ResourceTodos:
		if id > 0 {
			var item Todo
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				return nil, err
			}
			result = append(result, item)
		} else {
			var items []Todo
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return nil, err
			}
			for _, item := range items {
				result = append(result, item)
			}
		}

	case ResourceUsers:
		if id > 0 {
			var item User
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				return nil, err
			}
			result = append(result, item)
		} else {
			var items []User
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return nil, err
			}
			for _, item := range items {
				result = append(result, item)
			}
		}
	}

	return result, nil
}

func (c *APIClient) FetchRelated(resource ResourceType, id int, related ResourceType) ([]interface{}, error) {
	url := fmt.Sprintf("%s/%s/%d/%s", c.baseURL, resource, id, related)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, &APIError{Message: "не удалось создать запрос", Err: err}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, &APIError{Message: "не удалось получить связанные данные", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			Message: fmt.Sprintf("неожиданный статус код: %d", resp.StatusCode),
			Status:  resp.StatusCode,
		}
	}

	return c.decodeResponse(resp, related, 0)
}

type ResourceProcessor struct{}

func NewResourceProcessor() *ResourceProcessor {
	return &ResourceProcessor{}
}

func (p *ResourceProcessor) Process(resource interface{}) ProcessedResource {
	now := time.Now()

	switch v := resource.(type) {
	case Post:
		return ProcessedResource{
			ID:       v.ID,
			Type:     ResourcePosts,
			Original: v,
			Summary:  fmt.Sprintf("Пост от пользователя %d: %s", v.UserID, truncate(v.Title, 30)),
			Details: map[string]interface{}{
				"user_id":    v.UserID,
				"title":      v.Title,
				"body":       truncate(v.Body, 50),
				"word_count": len(strings.Fields(v.Body)),
			},
			ProcessedAt: now,
		}

	case Comment:
		return ProcessedResource{
			ID:       v.ID,
			Type:     ResourceComments,
			Original: v,
			Summary:  fmt.Sprintf("Комментарий к посту %d: %s", v.PostID, truncate(v.Name, 30)),
			Details: map[string]interface{}{
				"post_id": v.PostID,
				"name":    v.Name,
				"email":   v.Email,
				"body":    truncate(v.Body, 50),
			},
			ProcessedAt: now,
		}

	case Album:
		return ProcessedResource{
			ID:       v.ID,
			Type:     ResourceAlbums,
			Original: v,
			Summary:  fmt.Sprintf("Альбом от пользователя %d: %s", v.UserID, truncate(v.Title, 30)),
			Details: map[string]interface{}{
				"user_id": v.UserID,
				"title":   v.Title,
			},
			ProcessedAt: now,
		}

	case Photo:
		return ProcessedResource{
			ID:       v.ID,
			Type:     ResourcePhotos,
			Original: v,
			Summary:  fmt.Sprintf("Фото из альбома %d: %s", v.AlbumID, truncate(v.Title, 30)),
			Details: map[string]interface{}{
				"album_id":  v.AlbumID,
				"title":     v.Title,
				"url":       v.URL,
				"thumbnail": v.ThumbnailURL,
			},
			ProcessedAt: now,
		}

	case Todo:
		status := "ожидает"
		if v.Completed {
			status = "выполнено"
		}
		return ProcessedResource{
			ID:       v.ID,
			Type:     ResourceTodos,
			Original: v,
			Summary:  fmt.Sprintf("Задача для пользователя %d: %s [%s]", v.UserID, truncate(v.Title, 30), status),
			Details: map[string]interface{}{
				"user_id":   v.UserID,
				"title":     v.Title,
				"completed": v.Completed,
				"status":    status,
			},
			ProcessedAt: now,
		}

	case User:
		return ProcessedResource{
			ID:       v.ID,
			Type:     ResourceUsers,
			Original: v,
			Summary:  fmt.Sprintf("Пользователь: %s (%s)", v.Name, v.Username),
			Details: map[string]interface{}{
				"name":     v.Name,
				"username": v.Username,
				"email":    v.Email,
				"phone":    v.Phone,
				"website":  v.Website,
			},
			ProcessedAt: now,
		}

	default:
		return ProcessedResource{
			ID:          0,
			Type:        ResourcePosts,
			Summary:     "Неизвестный тип ресурса",
			Details:     map[string]interface{}{},
			ProcessedAt: now,
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

type OutputFormatter struct{}

func NewOutputFormatter() *OutputFormatter {
	return &OutputFormatter{}
}

func (f *OutputFormatter) FormatResource(res ProcessedResource) string {
	var emoji string
	var resourceName string

	switch res.Type {
	case ResourcePosts:
		emoji = "📝"
		resourceName = "Пост"
	case ResourceComments:
		emoji = "💬"
		resourceName = "Комментарий"
	case ResourceAlbums:
		emoji = "📀"
		resourceName = "Альбом"
	case ResourcePhotos:
		emoji = "🖼️"
		resourceName = "Фото"
	case ResourceTodos:
		emoji = "✅"
		resourceName = "Задача"
	case ResourceUsers:
		emoji = "👤"
		resourceName = "Пользователь"
	default:
		emoji = "📌"
		resourceName = "Ресурс"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n%s %s #%d\n", emoji, resourceName, res.ID))
	sb.WriteString(fmt.Sprintf("   📋 Кратко: %s\n", res.Summary))

	// Add details
	sb.WriteString("   📊 Детали:\n")
	for key, value := range res.Details {
		rusKey := translateKey(key)
		sb.WriteString(fmt.Sprintf("      • %s: %v\n", rusKey, value))
	}

	sb.WriteString(fmt.Sprintf("   🕐 Обработано: %s\n", res.ProcessedAt.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("   %s\n", strings.Repeat("─", 50)))

	return sb.String()
}

func (f *OutputFormatter) FormatHeader(config Config, count int) string {
	var filterDesc string
	resourceName := translateResource(string(config.Resource))

	switch {
	case config.ID > 0:
		filterDesc = fmt.Sprintf("%s #%d", resourceName, config.ID)
	case config.UserID > 0:
		filterDesc = fmt.Sprintf("%s от пользователя #%d", resourceName, config.UserID)
	default:
		filterDesc = fmt.Sprintf("Все %s", resourceName)
	}

	expandInfo := ""
	if config.Expand {
		expandInfo = " [со связанными данными]"
	}

	return fmt.Sprintf("\n📊 %s%s\n", filterDesc, expandInfo) +
		fmt.Sprintf("   Получено: %d элементов\n", count) +
		fmt.Sprintf("   %s\n", strings.Repeat("=", 60))
}

func translateResource(res string) string {
	switch res {
	case "posts":
		return "посты"
	case "comments":
		return "комментарии"
	case "albums":
		return "альбомы"
	case "photos":
		return "фото"
	case "todos":
		return "задачи"
	case "users":
		return "пользователи"
	default:
		return res
	}
}

func translateKey(key string) string {
	translations := map[string]string{
		"user_id":       "ID пользователя",
		"post_id":       "ID поста",
		"album_id":      "ID альбома",
		"title":         "заголовок",
		"body":          "содержание",
		"name":          "имя",
		"email":         "email",
		"phone":         "телефон",
		"website":       "сайт",
		"username":      "никнейм",
		"url":           "ссылка",
		"thumbnail":     "миниатюра",
		"completed":     "выполнено",
		"status":        "статус",
		"word_count":    "количество слов",
		"comment_count": "количество комментариев",
		"photo_count":   "количество фото",
		"post_count":    "количество постов",
		"todo_count":    "количество задач",
	}

	if translated, ok := translations[key]; ok {
		return translated
	}
	return key
}

// Этап 1
func fetchStage(config Config, client *APIClient) (<-chan interface{}, <-chan error) {
	out := make(chan interface{})
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		// Fetch main resources
		resources, err := client.Fetch(config.Resource, config.ID, config.UserID)
		if err != nil {
			errCh <- fmt.Errorf("этап загрузки: %w", err)
			return
		}

		// Apply limit if specified
		if config.Limit > 0 && len(resources) > config.Limit {
			resources = resources[:config.Limit]
		}

		log.Printf("Загружено %d %s", len(resources), translateResource(string(config.Resource)))

		for _, resource := range resources {
			out <- resource
		}
	}()

	return out, errCh
}

// Этап 2
func processStage(workerCount int, resources <-chan interface{}, processor *ResourceProcessor) <-chan ProcessedResource {
	out := make(chan ProcessedResource)

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for resource := range resources {
				log.Printf("Рабочий %d обрабатывает ресурс", workerID)
				time.Sleep(10 * time.Millisecond) // Simulate processing
				processed := processor.Process(resource)
				out <- processed
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Этап 3
func enrichStage(config Config, client *APIClient, processed <-chan ProcessedResource) (<-chan ProcessedResource, <-chan error) {
	out := make(chan ProcessedResource)
	errCh := make(chan error, 10)

	if !config.Expand {
		// Just pass through if no enrichment needed
		go func() {
			defer close(out)
			for res := range processed {
				out <- res
			}
		}()
		return out, errCh
	}

	go func() {
		defer close(out)

		for res := range processed {
			// Enrich based on resource type
			switch res.Type {
			case ResourcePosts:
				// Fetch comments for this post
				comments, err := client.FetchRelated(ResourcePosts, res.ID, ResourceComments)
				if err != nil {
					errCh <- fmt.Errorf("не удалось загрузить комментарии для поста %d: %w", res.ID, err)
				} else {
					res.Details["comment_count"] = len(comments)
				}

			case ResourceAlbums:
				// Fetch photos for this album
				photos, err := client.FetchRelated(ResourceAlbums, res.ID, ResourcePhotos)
				if err != nil {
					errCh <- fmt.Errorf("не удалось загрузить фото для альбома %d: %w", res.ID, err)
				} else {
					res.Details["photo_count"] = len(photos)
				}

			case ResourceUsers:
				// Fetch posts by this user
				posts, err := client.Fetch(ResourcePosts, 0, res.ID)
				if err != nil {
					errCh <- fmt.Errorf("не удалось загрузить посты для пользователя %d: %w", res.ID, err)
				} else {
					res.Details["post_count"] = len(posts)
				}

				// Fetch todos for this user
				todos, err := client.Fetch(ResourceTodos, 0, res.ID)
				if err != nil {
					errCh <- fmt.Errorf("не удалось загрузить задачи для пользователя %d: %w", res.ID, err)
				} else {
					res.Details["todo_count"] = len(todos)
				}
			}

			out <- res
		}
	}()

	return out, errCh
}

// Этап 4
func collectStage(processed <-chan ProcessedResource, limit int) <-chan ProcessedResource {
	out := make(chan ProcessedResource)

	go func() {
		defer close(out)

		resources := make([]ProcessedResource, 0)
		for res := range processed {
			resources = append(resources, res)
		}

		for i := 0; i < len(resources)-1; i++ {
			for j := i + 1; j < len(resources); j++ {
				if resources[i].ID > resources[j].ID {
					resources[i], resources[j] = resources[j], resources[i]
				}
			}
		}

		// Apply limit if specified
		if limit > 0 && len(resources) > limit {
			resources = resources[:limit]
		}

		// Output
		for _, res := range resources {
			out <- res
		}
	}()

	return out
}

// Stage 5: Format output
func formatStage(processed <-chan ProcessedResource) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		formatter := NewOutputFormatter()

		for res := range processed {
			out <- formatter.FormatResource(res)
		}
	}()

	return out
}

// Error aggregator
func errorAggregator(errChannels ...<-chan error) <-chan error {
	out := make(chan error)

	var wg sync.WaitGroup

	for _, errCh := range errChannels {
		wg.Add(1)
		go func(ch <-chan error) {
			defer wg.Done()
			for err := range ch {
				if err != nil {
					out <- err
				}
			}
		}(errCh)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	resource := flag.String("resource", "posts", "Ресурс: posts, comments, albums, photos, todos, users")
	id := flag.Int("id", 0, "Конкретный ID для загрузки (0 - все)")
	userID := flag.Int("user", 0, "Фильтр по ID пользователя")
	limit := flag.Int("limit", 0, "Лимит результатов (0 - без лимита)")
	workers := flag.Int("workers", 3, "Количество рабочих процессов")
	expand := flag.Bool("expand", false, "Загружать связанные ресурсы")
	timeout := flag.Duration("timeout", 10*time.Second, "Таймаут API")
	help := flag.Bool("help", false, "Показать справку")

	flag.Parse()

	if *help {
		printUsage()
		return
	}

	config := Config{
		Resource:   ResourceType(*resource),
		ID:         *id,
		UserID:     *userID,
		Limit:      *limit,
		Expand:     *expand,
		APIBaseURL: "https://jsonplaceholder.typicode.com",
		Timeout:    *timeout,
	}

	// Initialize components
	httpClient := &http.Client{Timeout: config.Timeout}
	apiClient := NewAPIClient(httpClient, config.APIBaseURL)
	processor := NewResourceProcessor()

	startTime := time.Now()

	// Этап 1
	resourcesCh, fetchErrCh := fetchStage(config, apiClient)
	// Этап 2
	processedCh := processStage(*workers, resourcesCh, processor)
	// Этап 3
	enrichedCh, enrichErrCh := enrichStage(config, apiClient, processedCh)
	// Этап 4
	collectedCh := collectStage(enrichedCh, config.Limit)
	// Этап 5
	outputCh := formatStage(collectedCh)
	// Ошибки
	errCh := errorAggregator(fetchErrCh, enrichErrCh)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("Веб-Клиент - %s", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("\n%s\n", strings.Repeat("=", 60))

	// Вывод ошибок
	go func() {
		for err := range errCh {
			log.Printf("⚠️ Ошибка: %v\n", err)
		}
	}()

	// Вывод результата
	count := 0
	for line := range outputCh {
		fmt.Print(line)
		count++
	}

	duration := time.Since(startTime)
	resourceName := translateResource(string(config.Resource))
	fmt.Printf("\n✅ Обработано %d %s за %v\n", count, resourceName, duration)
}

func printUsage() {
	fmt.Println("Веб-Клиент")
	fmt.Println("Автор: Давлетов Динис")
	fmt.Println("\nИспользование: json-client [опции]")
	fmt.Println("\nРесурсы:")
	fmt.Println("  posts     - 100 постов")
	fmt.Println("  comments  - 500 комментариев")
	fmt.Println("  albums    - 100 альбомов")
	fmt.Println("  photos    - 5000 фото")
	fmt.Println("  todos     - 200 задач")
	fmt.Println("  users     - 10 пользователей")
	fmt.Println("\nОпции:")
	fmt.Println("  -resource string  Тип ресурса (по умолчанию \"posts\")")
	fmt.Println("  -id int           Конкретный ID для загрузки (0 - все)")
	fmt.Println("  -user int         Фильтр по ID пользователя")
	fmt.Println("  -limit int        Лимит результатов (0 - без лимита)")
	fmt.Println("  -workers int      Количество рабочих горутин (по умолчанию 3)")
	fmt.Println("  -expand           Загружать связанные ресурсы (комментарии к постам и т.д.)")
	fmt.Println("  -timeout duration Таймаут API (по умолчанию 10s)")
	fmt.Println("  -help             Показать эту справку")
	fmt.Println("\nПримеры:")
	fmt.Println("  json-client -resource=posts -user=1 -limit=5")
	fmt.Println("  json-client -resource=posts -id=3 -expand")
	fmt.Println("  json-client -resource=users -expand")
	fmt.Println("  json-client -resource=albums -user=1 -limit=3")
	fmt.Println("  json-client -resource=todos -user=2 -limit=10")
	fmt.Println("\nAPI: https://jsonplaceholder.typicode.com")
}
