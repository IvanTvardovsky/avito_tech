## Старцев Иван  |  Авито  |  Тестовое задание  |  Apr 2024

### Инструкция по запуску
1. Прогон интеграционных тестов:
* Заменить в *config.yml*:
  > app_mode: test: true

* > docker-compose up

2. Запуск сервиса:
* Заменить в *config.yml*:
  > app_mode: test: false
  >
* > docker-compose up go_app my_db 

3. Если меняется *config.yml*, то надо пересобрать image для rest-api:

* > docker rm rest-server db
* > docker image rm avito_tech-go_app

4. Прогон линтера из корневой директории:
   
   > golangci-lint -c .golangci.yml run ./...

### Вопросы, с которыми столкнулся

1. Дата проследнего обновления на сервере или дата обновления самого баннера с помощью PUT?
   >
   
          updated_at:
         
          type: string
      
          format: date-time
         
          description: Дата обновления баннера
   >
   Был выбран второй вариант, так как я думаю, что его полезность больше

2. Если у баннера !isActive и запрос от пользователя, то возвращать, что баннер не найден или пустой блок? Код ответа?
  
   Код 200, пустой блок
  
