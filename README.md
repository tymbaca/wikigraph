# Wikigraph

Tool for parsing wikipedia articles graph.

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
the program at any moment by pressing `<Ctrl-c>` (see [graceful shutdown](#usagedev-features)). 

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

Notice that you can run this at any time. You don't need to parse all wikipedia :)

The exported graph will loop something like this:

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

**Graceful Shutdown**. Any time while program is executing you can press 
`<Crlt-c>`. Program will wait until all workers will parse and save their 
results and only then exits.

**Rate limiter**. Program uses internal http client with rate limiting by 
default set to 20 RPS (it's hardcoded). I found this rate most optimal. If you 
get `429 Too Many Requests` then just wait a bit and try again. Or you can can 
change the rate (in `cmd/wikigraph/main.go`) in code and recompile the program. 
I'm too lazy to add RPS flag (just look at how I handle cli arguments in main.go 
lol).
