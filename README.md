# Wikigraph

Tool for parsing Wikipedia articles graph.

It takes the provided initial wikipedia links, parses them, gets the child links 
and adds them to parsing queue. It also saves parsed article name.

Your can export collected links as a CSV table (see [usage](#usage)) to be used 
in external software. For example you can directly import that file to 
[Cosmograph](https://cosmograph.app/run/) to visualize the graph:

<p float="left">
  <img src="https://github.com/user-attachments/assets/d398b28a-5028-44bd-8049-0aa1d98df3c5" width="400" /> 
  <img src="https://github.com/user-attachments/assets/0db5a20f-9c5f-4f6c-add5-4dc61814a152" width="400" />
</p>

## Install

```
go install github.com/tymbaca/wikigraph@latest
```

## Usage

### Parse

Launch parse process by doing this:

```
wikigraph parse my_graph.db https://en.wikipedia.org/wiki/Kingdom_of_Greece
```

Or:

```
wikigraph parse my_graph.db https://en.wikipedia.org/wiki/Kingdom_of_Greece https://en.wikipedia.org/wiki/Christmas
```

Program will begin parsing the wikipedia, starting from provided URLs. You can exit 
the program at any moment by pressing `<Ctrl-c>` (see [graceful shutdown](#good-thingies)). 

Your can continue by launching program with already existing database file:

```
wikigraph parse my_graph.db
```

It will continue parsing as expected (without loosing the progress). 
Also it will retry all links that it failed to parse in previous attempt.

### Export

Now you can export the graph to CSV file:

```
wikigraph export my_graph.db my_graph.csv
```

Notice that you can run this at any time. You don't need to parse all articles in Wikipedia :)

The exported graph will look something like this:

```csv
from,to
"Heil unserm König, Heil!",Kingdom of Greece
"Heil unserm König, Heil!",Hymn to Liberty
"Heil unserm König, Heil!",Greece
Constitutional monarchy,Absolute monarchy
Constitutional monarchy,State religion
Constitutional monarchy,Unitary state
...
```

### Help

Run `wikigraph help` for more info.

## Good thingies

**Worker Pool**. Program uses a pool of workers parallelize parsing workload.

**Job Queueing**. Every fetched child link is pushed to the job queue with 
PENDING status, so later on another worker can grab it, parse it and produce 
more child articles. Job queue is basically an `article` table in SQLite DB, 
you can explore it with any suitable DB client. If program exits, you still 
have all the queue in database, so you can restart easely.

**Graceful Shutdown**. Any time while program is executing you can press 
`<Crlt-c>`. Program will wait until all workers will parse and save their 
results and only then exits.

**Rate Limiter**. Program uses internal http client with rate limiting by 
default set to 20 RPS (it's hardcoded). I found this rate most optimal. If you 
get `429 Too Many Requests` then just wait a bit and try again. Or you can can 
change the rate (in `cmd/wikigraph/main.go`) in code and recompile the program. 
I'm too lazy to add RPS flag (just look at how I handle cli arguments in main.go 
lol).


## Language support
Program was tested with English, Russian and [Wolof](https://en.wikipedia.org/wiki/Wolof_language) Wikipedia.
So any other non-ascii language articles can be supported, as long as they match similar HTML layout.

## TODO
- [ ] Use official [MediaWiki API](https://www.mediawiki.org/wiki/API:Main_page) instead of parsing the whole HTML of every article.
