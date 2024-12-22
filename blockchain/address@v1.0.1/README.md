# address

## Приступая к работе над сервисом

1. Подготовка окружения
    Необходимо получить [api или read_api токен](https://code.emcdtech.com/-/user_settings/personal_access_tokens) от gitlab.
    Создаем файл `touch ~/.netrc` и выполняем `cat ~/.netrc machine code.emcdtech.com login readonly password YOUR_TOKEN`
    Установить переменную `go env -w GOPRIVATE=code.emcdtech.com`


2. Устанавливаем все необходимые зависимости

```bash
make install-tools
```

3. Кодогенерация (запускает protoc, go generate ./..., и т.д.)

```bash
make generate
```

P.S. Для того чтоб docker смог получить образ,

```bash
docker login --username readonly --password YOUR_TOKEN registry.emcdtech.com
```

4. Форматирование кода

```bash
make fmt
```

5. Запуск линтеров

```bash
make linter
```

6. Запуск и остановка зависимостей для локальной разработки/тестирования

запуск

```bash
make deps-up
```

остановка

```bash
make deps-down
```

## Переносим практику на свой сервис

При переносе практики на свой сервис необходимо внести коррективы в Makefile. Провести по файлу замену:

1. ``code.emcdtech.com/emcd/blockchain/address`` -> ``code.emcdtech.com/emcd/service/newservice``
2. ``/protocol/address`` -> ``/protocol/newservice``
3. В наборе команд generate, найти инструкцию mockgen и в ней задать правильное имя файла сформированного protoc grpc.  

## Стиль кода и концепция
### service
содержит всю бизнес логику основанную на моделях model. Сервису плевать кто его вызывает  http/grpc/worker он работает с моделью и репозиториями repository
### controller/worker
верхний уровень сервиса, он отвечает за внешнее апи, может работать с одинм или несколькими сервисами. Как правило на этом уровне происходит валидация и конвертация модели в request/response
#### repository
нижний уровень хранения данных. Данные могут хранится в базе,в редисе, в GRPC иного сервиса. Репозиторий точно так же работает с моделями, но не содержит никакой логики.

## Работа с ошибками
есть два варианта работы с ошибками.
1. Ошибки которые не возвращаются их нужно залогировать. Такое может быть в отдельных горутинах. В этом случае мы просто вызываем пакет log из SDK и логируем
2. Ошибки которые нужно вернуть. Тут ошибки не логируем а врапаем через `fmt.errorf` или `errors.wrap` первый вариант гибче, второй  современнее
3. Ошибки на верхнем уровне. Логируются все ошибки без исключениея на уровне controller и worker

## Именование
```
type Task struct {
	db *pgxpool.Pool
}

func NewTask(db *pgxpool.Pool) *Task {
	return &Task{
		db: db,
	}
}

func (r *Task) Create(ctx context.Context, task model.Task) error {
	const createTaskSQL = `
insert into task(id, name, status, created_at, deleted_at)
values($1, $2, $3, $4, $5)
returning id, "name", status, created_at, deleted_at
`
	_, err := r.db.Exec(ctx, createTaskSQL, task.ID, task.Name, task.Status, task.CreatedAt, task.DeletedAt)
	if err != nil {
		return fmt.Errorf("exec createTaskSQL: %w", err)
	}
	return nil
}
```
у нас уже есть пакет repository. Поэтому внутри мы называем `Task`
с кода идет вызов `repository.Task`
`func (r *Task) Create(ctx context.Context, task model.Task) error`
в ресивере используем `r` как первую букву репозитория

такой же концепт для всех остальных структур
