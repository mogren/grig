grig - Like rig, but in go
====
The program will generate a random address complete with street names and zip codes.

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
