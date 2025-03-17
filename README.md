## grig - Like rig, but in go

The program will generate any number random addresses, complete with names, street names, cities and zip codes.

To compile:
```
go build grig.go
```

To run:

```
> ./grig --help 
Usage of ./grig:
  -j=false: Print as JSON
  -l=false: List available ISO language codes
  -lang="en_us": Select ISO 639-1 language code, defaults to USA
  -n=1: Number of identities to output
  -v=false: Verbose output
  -x=false: Print as XML
```    
Example output:
```
> ./grig
Carroll Olsen
31 Glenwood Ave
33730 St. Petersburg

> ./grig -j -lang no
{
  "firstname": "Kjellaug",
  "lastname": "Sviland",
  "street": "Øvre Gjelrustvegen",
  "nr": 55,
  "zip": 1101,
  "city": "Eigersund"
}

> ./grig -x -lang sv 
<Rig>
  <firstname>Britta</firstname>
  <lastname>Fransson</lastname>
  <street>Karlsviksgatan</street>
  <nr>45</nr>
  <zip>39363</zip>
  <city>Kalmar</city>
</Rig>
```

Performance of different modes when generating 1,000,000 random addresses on an M1 Pro CPU:
```
> time ./grig -n 1000000 >| /dev/null
./grig -n 1000000 >| /dev/null  0.68s user 0.30s system 87% cpu 1.109 total
> time ./grig -j -n 1000000 >| /dev/null
./grig -j -n 1000000 >| /dev/null  0.73s user 0.31s system 110% cpu 0.940 total
> time ./grig -x -n 1000000 >| /dev/null
./grig -x -n 1000000 >| /dev/null  3.42s user 0.34s system 178% cpu 2.105 total
```

To run the benchmark tests:

```
> go test -bench=.
goos: darwin
goarch: arm64
pkg: github.com/mogren/grig
cpu: Apple M1 Pro
BenchmarkGenerateIdentities-10                   9449667               111.7 ns/op
BenchmarkGenerateMultipleIdentities-10             84114             14260 ns/op
BenchmarkAsJSON-10                               4673596               246.3 ns/op
BenchmarkAsXML-10                                 813054              1342 ns/op
BenchmarkLoadData-10                                3876            323048 ns/op
PASS
ok      github.com/mogren/grig  10.016s
```

#### TODO
* Add more locales
* Make a simpler Roulette-randomizer for smaller datasets
* Optimise Vose
* Split Vose into a separate package
* Add web-server mode
* Correct weights for Swedish data


