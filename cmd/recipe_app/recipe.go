package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (app *application) getRecipe(c *gin.Context) {
	recipeID, err := app.readIDParam(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Рецепт не найден!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Отображение информации о рецепте %d", recipeID)})

	recipe := Recipe{
		ID:          int64(recipeID),
		Time:        time.Now(),
		Title:       "Паста с помидорами",
		Description: "Простой рецепт пасты с помидорами и базиликом.",
		Ingredients: []string{"200 г пасты", "2 помидора", "Свежий базилик"},
		Steps:       []string{"Сварите пасту по инструкции.", "Нарежьте помидоры и базилик.", "Смешайте готовую пасту с помидорами и базиликом."},
		Author:      1,
	}

	err = writeJSON(c.Writer, http.StatusOK, Envelope{"recipe": recipe}, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The server encountered a problem and could not process your request"})
	}

}

const jsonR = `[
    {
        "ID": 1,
        "Time": "2023-10-01T10:30:00Z",
        "Title": "Паста с помидорами",
        "Description": "Простой рецепт пасты с помидорами и базиликом.",
        "Ingredients": [
            "200 г пасты",
            "2 помидора",
            "Свежий базилик"
        ],
        "Steps": [
            "Сварите пасту по инструкции.",
            "Нарежьте помидоры и базилик.",
            "Смешайте готовую пасту с помидорами и базиликом."
        ],
        "Author": 12345
    },
    {
        "ID": 2,
        "Time": "2023-10-02T09:45:00Z",
        "Title": "Салат Цезарь",
        "Description": "Классический салат Цезарь с курицей.",
        "Ingredients": [
            "300 г куриного филе",
            "Салат Романо",
            "Гренки",
            "Пармезан",
            "Соус Цезарь"
        ],
        "Steps": [
            "Поджарьте куриное филе и нарежьте его кубиками.",
            "Приготовьте гренки и салат Романо.",
            "Смешайте все ингредиенты и добавьте соус Цезарь."
        ],
        "Author": 67890
    },
    {
        "ID": 3,
        "Time": "2023-10-03T14:15:00Z",
        "Title": "Суп с лапшой",
        "Description": "Традиционный суп с лапшой и овощами.",
        "Ingredients": [
            "200 г лапши",
            "2 моркови",
            "1 лук",
            "Сельдерей",
            "Бульон"
        ],
        "Steps": [
            "Обжарьте лук и морковь, добавьте сельдерей.",
            "Добавьте бульон и варите лапшу."
        ],
        "Author": 13579
    },
    {
        "ID": 4,
        "Time": "2023-10-04T18:00:00Z",
        "Title": "Омлет с овощами",
        "Description": "Вкусный омлет с помидорами и шпинатом.",
        "Ingredients": [
            "3 яйца",
            "Помидоры",
            "Шпинат",
            "Соль, перец"
        ],
        "Steps": [
            "Взбейте яйца, добавьте помидоры и шпинат.",
            "Готовьте на сковороде до золотистой корки."
        ],
        "Author": 24680
    },
    {
        "ID": 5,
        "Time": "2023-10-05T12:30:00Z",
        "Title": "Картофельный суп",
        "Description": "Сытный суп с картофелем и грибами.",
        "Ingredients": [
            "Картофель",
            "Грибы",
            "Лук",
            "Сливки",
            "Бульон"
        ],
        "Steps": [
            "Обжарьте лук и грибы, добавьте картофель и бульон.",
            "Подавайте с сливками."
        ],
        "Author": 98765
    },
    {
        "ID": 6,
        "Time": "2023-10-06T20:00:00Z",
        "Title": "Спагетти с морепродуктами",
        "Description": "Итальянские спагетти с морепродуктами и томатным соусом.",
        "Ingredients": [
            "Спагетти",
            "Креветки",
            "Мидии",
            "Томаты",
            "Чеснок"
        ],
        "Steps": [
            "Обжарьте креветки и мидии с чесноком, добавьте томаты.",
            "Подавайте со спагетти."
        ],
        "Author": 112233
    }
]
`

func addRecipe(c *gin.Context) {

	jsonData, err := json.Marshal(jsonR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при кодировании JSON"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)

}
