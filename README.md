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

Performance of different modes when generating 1,000,000 random addresses:
```
> time ./grig -n 1000000 >| /dev/null                                                                                                                                                                                                                          /Volumes/Unix/go/src/github.com/mogren/grig
./grig -n 1000000 >| /dev/null  1.79s user 0.38s system 89% cpu 2.418 total
> time ./grig -j -n 1000000 >| /dev/null                                                                                                                                                                                                                       /Volumes/Unix/go/src/github.com/mogren/grig
./grig -j -n 1000000 >| /dev/null  3.44s user 0.43s system 107% cpu 3.607 total
> time ./grig -x -n 1000000 >| /dev/null                                                                                                                                                                                                                       /Volumes/Unix/go/src/github.com/mogren/grig
./grig -x -n 1000000 >| /dev/null  7.01s user 0.80s system 122% cpu 6.379 total
```

#### TODO
* Add more locales
* Make a simpler Roulette-randomizer for smaller datasets
* Optimise Vose
* Split Vose into a separate package
* Add web-server mode
* Correct weights for Swedish data


