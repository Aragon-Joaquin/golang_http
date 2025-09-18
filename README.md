i used go-blueprint to give me a hand on the project structure

## USING/doing to do:

- chi for a more easy and lightweight routing
- [x] psql connection
- [ ] kafka
- [ ] elastic search
- [x] swagger **kinda done??**
- [ ] ratelimiter **<- priotize this**
- [ ] docker
- [ ] websockets **<- priotize this**
- [ ] redis cache
- [ ] migrations
- [ ] logger
- [ ] make public file server
- [x] jwt/authentication

**db name**: testgodb (yes, thats the name)

## steps to work with psql:

1. next commands will be executed inside this shell + initialize psql (if freshly installed, else, can be skipped) and enable it on startup

```sh
sudo -iu postgres


#commands if psql was freshly installed
initdb --locale $LANG -E UTF8 -D '/var/lib/postgres/data'
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

2. create user

```sh
createuser --interactive
```

3. create db.

```sh
createdb YOUR-DB-NAME!!
```

4. access another shell inside this shell idk what the fuck is this

```sh
psql YOUR-DB-NAME!!
```

5. useful commands when using the shell on a shell on a shell i really dont have any idea why is like that

- \dt: show tables

## or just use dbeaver
