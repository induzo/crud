# A little opinionated CRUD interface

[![Documentation](https://godoc.org/github.com/induzo/crud?status.svg)](http://godoc.org/github.com/induzo/crud) [![Go Report Card](https://goreportcard.com/badge/github.com/induzo/crud)](https://goreportcard.com/report/github.com/induzo/crud) [![Maintainability](https://api.codeclimate.com/v1/badges/3e4ee9ac6a7a39a18c36/maintainability)](https://codeclimate.com/github/induzo/crud/maintainability) [![Coverage Status](https://coveralls.io/repos/github/induzo/crud/badge.svg?branch=master)](https://coveralls.io/github/induzo/crud?branch=master) [![CircleCI](https://circleci.com/gh/induzo/crud.svg?style=svg)](https://circleci.com/gh/induzo/crud)

## Use case

Let's take a simple use case, you have an entity A, and you want to easily enable a CRUD API.

Simply create a "manager" fulfilling the induzo/crud MgrI interface for the entity A you want to enable CRUD for.
This interface (described in mgr.go) is as follows:

```golang
type MgrI interface {
    NewEmptyEntity() interface{}
    Create(context.Context, interface{}, io.Reader) (interface{}, error)
    Delete(context.Context, xid.ID) error
    Get(context.Context, xid.ID) (interface{}, error)
    GetList(context.Context, ListModifiers) (interface{}, error)
    Update(context.Context, xid.ID, interface{}, io.Reader) (interface{}, error)
    PartialUpdate(context.Context, xid.ID, PartialUpdateData, io.Reader) error
    MapErrorToHTTPError(error) *gohttperror.ErrResponse
}
```

You can see an example of implementation in the [mock](./mock) folder.

Once this is done, you can just use this newly created manager and wrap it to enable the API.

You want to spawn a REST API, following the std library http handler?
Easy, just wrap you manager with the [REST implementation in the rest folder](./rest).

And there you go, you have a complete REST API in 5 lines

## Example

A very short example with the rest wrapper in the [example folder](./example).

## Coming soon

If this interface is successful, we are planning to simply add more wrappers on top of the REST one:

- GRPCWeb
- CQRS

## Opinions

- Your entity ids should be using github.com/rs/xid
- Your crud errors should be handlable by an httpresponse
