go install github.com/pressly/goose/v3/cmd/goose@latest
git clone https://github.com/meehighlov/eventor.git eventor-tmp
./goose -dir=eventor-tmp/migrations sqlite3 eventor.db up
./goose -dir=eventor-tmp/migrations sqlite3 eventor.db status
rm -rf eventor-tmp