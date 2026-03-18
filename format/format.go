package format

import (
	"fmt"
	"json-client/models"
)

// Formatter - вывод
type Formatter struct{}

// NewFormatter - конструктор вывода
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Вывод поста
func (f *Formatter) PrintPost(p models.Post) {
	fmt.Printf("\nПост %d\n", p.ID)
	fmt.Printf("	Заголовок: %s\n", p.Title)
	fmt.Printf("	Содержание: %s\n", p.Body)
	fmt.Println("")
}

// Вывод комментария
func (f *Formatter) PrintComment(c models.Comment) {
	fmt.Printf("\nКомментарий %d\n", c.ID)
	fmt.Printf("	Название: %s\n", c.Name)
	fmt.Printf("	Почта: %s\n", c.Email)
	fmt.Printf("	Содержание: %s\n", c.Body)
	fmt.Println("")
}

// Вывод ошибки
func (f *Formatter) PrintError(err error) {
	fmt.Printf("Ошибка: %v\n", err)
}
