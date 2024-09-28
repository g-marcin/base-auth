
run db

```bash
	docker-compose up --build
```

open db container
```bash
	docker exec -it db psql -U user
```

show tables
```bash
	\dt
```

show records
```bash
	SELECT * from users
```

run server

```bash
	make run
```

open browser
```bash
	http://localhost:8080
```