# book-alert

RSS feed for new publications of your favorite authors.

This is currently not running in any public place, so if you want to use it you'll have to self-host it.

## Building and running

You only need to install Go, or use the Go provided by the Nix flake development environment.

```
go build
PORT=8080 ./book-alert
```

Now you can query the server for the latest books by your favorite authors, like this:

```
curl 'localhost:8080?author=Q65969383&author=Q6774606'
```

Every author gets their own `author=<id>` query parameter. The IDs can be found by searching for the author by name on wikidata.org. The id is printed behind the author name in the page title ([example](https://www.wikidata.org/wiki/Q6774606)).

Book Alert returns the last 10 published books, regardless of who wrote them.

## Thanks

Many thanks to the amazing wikidata.org project! I tried to get book information from a couple of different places, and this was the easiest one to integrate with I could find, thanks to some great documentation.
