# OOWLISH (geyslan-golang)

## Golang Engineer

### How to run this Solution

In the project folder

Open the local [consume.html](./consume.html) file.

Run *docker* making sure you're member of the *docker* group.

```bash
❯ docker build -t oowlish . && docker run -it -p4000:4000 --name oowlish oowlish
```

When the *docker run* start, as the service logs show up in the terminal, click on _Get Logs_ button in the already open _web page_ to consume the REST API and see the mutable state of the database.

#### Testing the Watch logic

Right after *Watch* sleeps (terminal: `Watch sleeping for 5m0s`), you can insert new log lines; for that just connect to the running container *oowlish* using *docker exec*.

```bash
❯ docker exec -it oowlish bash
```

And append to the *test.access.log* file using the `watch_test.sh` script.

```bash
root@7bd2108b5067:/oowlish# ./watch_test.sh
```

Wait until *Watch* wake up again for new entry logs be readily available through the API.

---

### Solving the Challenge

For conciseness, _FLAT_ was the project structure type used, defining all source files as the package main in the base tree. Please don't think I'm disorganised. =)

> 1. Create micro-service that reads the attached log file (test.access.log) every five minutes;

The solution needs to know what log file to load. So the first argument passed to the executable by the OS environment makes it possible. Eg.

```bash
❯ oowlish ~/projs/go/src/oowlish/test.access.log
```

Calling `watchLog()` in the background was the method used to continuously read the log input file even when it reaches `EOF`.

For each line read `watchLog()` calls `parse()`.

After sleeping for 5 minutes, if the file size is changed, it attempts to read new data again. If during the sleep the file is truncated its offset is reset and the whole file is parsed again.

> 2. Parse each line to the application. The application should be able to filter it by GET or POST http methods. We are assuming that you will create a parser to transform each register in a model to be handled by the service;

The parser was implemented using `fmt.Sscanf()` which fills the `clf` (Common Log Format) structure. Filtering GET and POST was done by testing the Method field of the `clf` in the `watchLog()` loop.

> 3. Store these filtered entries on a database;

_Postgresql_ was the database of choice.

With security in mind no connection data was hardcoded. `dbConnect()` gets the variables from the environment. Eg.

```bash
❯ DBHOST=localhost DBUSER=oowlish DBPASS=oowlish DBDATABASE=oowlish go run oowlish ~/projs/go/src/oowlish/test.access.log
```

During the `watchLog()` loop `dbInsertLog()` is called for each (GET or POST) log line parsed.

> 4. Each log entry must appear once in a table. We consider it unique if these criteria match (client host, received request date-time, request line, size of the object returned);

To accomplish this requirement the `logs` table was created with the criterion `UNIQUE (client_host, date_time, resource, size)`. In insertion `ON CONFLICT DO NOTHING` is used making possible to ignore conflicts.

> 5. Implement a REST API that exposes one endpoint to filter the entries from the saved log lines, we must filter it by an specific HTTP method GET or POST;

The REST API was done using gorilla/mux lib due its closer api to go.

> 6. We'd love to see Structs and Channels on this test;

- `clf` Common Log Format structure
- goroutine `go watchLog(os.Args[1], os.Stdout)`
- Channel `go dbInsertLog(insChan, db, cl)`

> 7. Create a simple web page to consume the REST API (you can use auto-documentation);

The _consume.html_ file shows how the REST API can be consumed.

> 8. Write README.md instructions on how to get your code up-and-running;

This README.md does that.

> 9. Send us your code using a private Gitlab ([@oowlish-career](https://gitlab.com/oowlish-career)) or Github ([@oowlishcareer](https://github.com/oowlishcareer)) repository. Add to the project the respective Oowlish user as a "Maintainer" if applicable.

**Mission accomplished!** Please keep me informed. And again thank you for this opportunity.
