# URL Shortener

This is a url shortener similar to bitly.com.

It converts a long url to a short url in 7 characters, and converts the short url back to the original ones.

Posting the same long url multiple times will produce different short links.

# Start

Run instructions below to start the server.

Visit http://localhost:8080 to play.

```bash
make run
```

# Test

This demo is 100% test covered. The trickiest part is the db mocking.

```bash
# run test only
make test
# run test with coverage
make coverage
```

## Mock

The test cases mock database using [gomock](https://github.com/golang/mock).
