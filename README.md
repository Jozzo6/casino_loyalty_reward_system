# Casino loyalty reward system

### Requirements

- Docker
- Go (Golang)

### Description

App will run in 3 separated services `user`, `promotions`, and `notifications`.

`user` service handles authentication of user and CRUD operations.
`promotions` service handles CRUD operations for promotions and assigning promotions to the user.
`notifications` service handles sending notifications to the user.

On user registration event will be sent to the `promotions` service from `users` service trough redis. Which will add `Welcome promotion` to the user.
When any promotion is added to user it will once again send event trough redis from `promotions` service to `notifications` service which will than send notification that user has recived that promotion.

### How to run the app debug mode

Running following command it will start database and redis

```
	docker-compose -f docker-compose.debug.yaml up --build
```

Use this commands it will clear everything and start fresh

```
	docker-compose -f docker-compose.debug.yaml stop
    docker-compose -f docker-compose.debug.yaml rm -vf
```

In `./vscode/launch.json` configurations for debugging of each of the services is defined and can be ran by itself.

| Service Name          | URL                   |
| --------------------- | --------------------- |
| User service          | localhost:3003/api/v1 |
| Promotions service    | localhost:3002/api/v1 |
| Notifications service | localhost:3003/api/v1 |

### How to run the app

Running following command will start database. It will also start 2 copies of each service.

```
	docker-compose -f docker-compose.yaml up --build
```

Use this commands it will clear everything and start fresh

```
	docker-compose -f docker-compose.yaml stop
    docker-compose -f docker-compose.yaml rm -vf
```

nginx is used to distribute load through the services that are running.

| Service Name          | URL                               |
| --------------------- | --------------------------------- |
| User service          | localhost:80/users/api/v1         |
| Promotions service    | localhost:80/promotions/api/v1    |
| Notifications service | localhost:80/notifications/api/v1 |

---

The database will be pre-filled with some sample data for testing purposes. This ensures that you can immediately test the functionality. In `promotions` table it is required to have one promotion of type `welcome_bonus` which is then used on user registration.

Test staff acccount

```
	email - john@example.com
	password - password
```

Test user account

```
	email - marc@example.com
	password - password
```

### Documantation

Swagger documenation is available in `/docs/api` folder
