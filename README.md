# Golang пакет, реализующий очередь воркеров

Пакет `poolk` реализует worker pool для параллельного выполнения задач. Пулл управляет фиксированным количеством горутин-воркеров и очередью задач с ограниченной емкостью.

Проект выложен на [Github](https://github.com/Libod3/poolk).

## Описание

**WorkerPool** - центральная структура, реализующая интерфейс пула:

* Управляет жизненным циклом воркеров
* Управляет очередь задач
* Предоставляет статистику выполнения

**Ключевые особенности:**

* Контролируемый параллелизм через фиксированное количество воркеров
* Буферизованная очередь для ограничения её емкости
* Graceful shutdown - корректное завершение выполнения задач
* Поддержка callback-хуков для пост-обработки задач
* Корректная обработка и логирование panic в задачах

## Пример использования

**Создание worker pool:**
```go
wp, err := pool.NewWorkerPool(3, 5)
if err != nil {
	return fmt.Errorf("create worker pool: %w", err)
}
```

**Добавление callback hook:**
```go
err := wp.SetDoneCallback(func() {
	fmt.Println("task is over")
})
if err != nil {
	return fmt.Errorf("set done callback: %w", err)
}
```

**Добавление задачи в очередь:**
```go
for i := 1; i <= 10; i++ {
	err := wp.Submit(func() {
		fmt.Printf("task %d started...\n", i)
		time.Sleep(time.Second)
	})
	if err != nil {
		return fmt.Errorf("submit task %d: %w", i, err)
	}
}
```

**Завершение выполнения:**
```go
err := wp.Stop()
if err != nil {
	return fmt.Errorf("stop worker pool: %w", err)
}
```

## Тесты

Запуск тестов:

```bash 
go test -v
```



