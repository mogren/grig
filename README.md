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
#### TODO
* Add more locales
* Make a simpler Roulette-randomizer for smaller datasets
* Optimise Vose
* Split Vose into a separate package
* Add web-server mode
* Correct weights for swedish data


