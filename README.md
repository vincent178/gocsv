gocsv
=====

> parse csv file to slice of structs with one line code

## Requirement

require go 1.18 which support generics

## Usage

* Basic

Define a struct:

```go
type Foo struct {
  Name   string
  Count  uint                   // speicify type other than string
  Enable bool `csv:"is_enable"` // speicify csv tag which used to mapping csv header  
}
```

Read a csv file:

```csv
name,count,is_enable,other key
bar,1,false,kk
```

Parse like this, need to specify the return type with go generics way `[Foo]`(Go can't infer type based on the return value type you speicify on the left side):

```go
func main() {
  f, err := os.Open("./data.csv")
  if err != nil {
    panic(err)
  }
  defer f.Close()

  ret, err := gocsv.Read[Foo](f) // speicify type here

  if err != nil {
    panic(err)
  }

  for _, r := range ret {
    fmt.Printf("%+v\n", r)
  }
}
```

You can check the full example in [examples](./examples/) folder.

* Suppress Error

In the case where do you want to deal with data imcompatible with field type definition, you can tell gocsv to suppress those errors.
```go
ret, err := gocsv.Read[Foo](f, gocsv.WithSuppressError(true))
```
