package main

type Car struct {
    ID      int     `json:"id"`
    Brand   string  `json:"brand"`   // марка автомобиля
    Model   string  `json:"model"`   // модель
    Year    int     `json:"year"`    // год выпуска
    Price   float64 `json:"price"`   // цена в тыс. рублей
}
