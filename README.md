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

Parse like this, you need to specify the type parameter `[Foo]` for the return value (Go can't infer type based on the return value type you speicify on the left side):

```go
func main() {
  f, err := os.Open("./data.csv")
  if err != nil {
    panic(err)
  }
  defer f.Close()

  ret, err := gocsv.Read[Foo](f) // speicify type here, return []*T

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

If you want to ignore data imcompatibility with field type definition, you can tell gocsv to suppress those errors.
```go
ret, err := gocsv.Read[Foo](f, gocsv.WithSuppressError(true))
```

* Multiple overrides

If you want to map another header to the same field:    
```go
type Foo struct {
  Name   string
  Count  uint   
  Enable bool `csv:"is_enable,isEnabled"` // speicify multiple csv tag which used to mapping csv header  
}
```
