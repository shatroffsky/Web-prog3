package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Record - структура для збереження одного показника
type Record struct {
	DeviceName string
	Voltage    int
	Date       string
}

// PageData - структура для передачі даних у HTML-шаблон
type PageData struct {
	Error    string
	FormData Record   // Щоб зберігати введені дані при помилці
	Records  []Record // Список всіх збережених записів
}

// Тимчасова "база даних" у пам'яті (зникне при перезапуску сервера)
var recordsDB = []Record{}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Дані для шаблону за замовчуванням
	data := PageData{
		Records: recordsDB,
	}

	// Якщо метод запиту POST — це означає, що користувач відправив форму
	if r.Method == http.MethodPost {
		// Парсимо форму
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Помилка обробки форми", http.StatusBadRequest)
			return
		}

		// Отримуємо значення з полів форми
		deviceName := r.FormValue("deviceName")
		voltageStr := r.FormValue("voltage")
		date := r.FormValue("date")

		// Зберігаємо введені дані, щоб повернути їх у форму в разі помилки
		data.FormData = Record{DeviceName: deviceName, Date: date}

		// 1. ВАЛІДАЦІЯ: Перевірка на порожні поля
		if deviceName == "" || voltageStr == "" || date == "" {
			data.Error = "Всі поля обов'язкові для заповнення!"
		} else {
			// 2. ВАЛІДАЦІЯ: Перевірка, чи напруга є числом
			voltage, err := strconv.Atoi(voltageStr)
			if err != nil || voltage <= 0 {
				data.Error = "Напруга повинна бути додатнім числом!"
			} else {
				// Якщо валідація успішна — зберігаємо дані
				newRecord := Record{
					DeviceName: deviceName,
					Voltage:    voltage,
					Date:       date,
				}
				recordsDB = append(recordsDB, newRecord)

				// РОБИМО ПЕРЕНАПРАВЛЕННЯ (Redirect), щоб уникнути повторної відправки форми при оновленні сторінки
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return // Завершуємо виконання функції, щоб шаблон не рендерився двічі
			}
		}
	}

	// Рендеримо HTML шаблон (виконується для GET запитів і для POST з помилками)
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Помилка завантаження шаблону", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", formHandler)

	log.Println("Сервер запущено: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Помилка:", err)
	}
}
