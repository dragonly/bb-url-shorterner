# URL Shortener

This is a url shortener similar to bitly.com.

It converts a long url to a short url in 7 characters, and converts the short url back to the original ones.

Posting the same long url multiple times will produce different short links.

# Instructions

## start server

Run instructions below to start the server. It will also create a sqlite3 database file if not present.

```bash
make run
```

## test and coverage

This demo is 100% test covered. The trickest part is the db mocking.

```bash
# run test only
make test
# run test with coverage
make coverage
```

# Mocking

The test cases mock database using `gomock`.
