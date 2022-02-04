package transaction

import (
	"bytes"
	"context"
	"encoding/gob"
	"unsafe"
)

type AnyPointer unsafe.Pointer

type TransactionalKey string

func (s TransactionalKey) String() string {
	return string(s)
}

type InMemTransaction struct {
	ctx     context.Context
	structs []*interface{}
	clones  []interface{}
}

func NewInitializerInMem(ctx context.Context) InMemTransaction {
	return InMemTransaction{
		ctx:    ctx,
		clones: make([]interface{}, 0),
	}
}

func (i *InMemTransaction) Add(structs ...*interface{}) error {
	i.structs = append(i.structs, structs...)
	return nil
}

type Cloneable interface {
	Clone(out interface{})
}

type Recovery interface {
	Recovery()
}

type CloneableFunc struct {
	out     interface{}
	in      interface{}
	cloneFn func() interface{}
}

func NewCloneableFunc(fn func() interface{}) CloneableFunc {
	return CloneableFunc{
		cloneFn: fn,
	}
}

func (c CloneableFunc) Clone() interface{} {
	fn := c.cloneFn()

	return fn
}

func (i *InMemTransaction) Begin(_ context.Context) (Transaction, error) {

	cloneable := i.ctx.Value("repository").(CloneableFunc)

	//cloneable := i.ctx.Value("repository").(Cloneable)

	i.clones = append(i.clones, cloneable.out)

	//fmt.Printf("clone: %T %v %s\n", cloneable.Clone(), cloneable.Clone(), (*cloneable.Clone()).(string))

	return i, nil
}

func (i InMemTransaction) Rollback() error {
	cloneable := i.ctx.Value("repository").(CloneableFunc)
	cloneable.in = i.clones[0]
	return nil
}

func (i InMemTransaction) Commit() error {
	cloneable := i.ctx.Value("repository").(CloneableFunc)
	cloneable.Clone()
	return nil
}

func NewCloneable(in, out interface{}) CloneableFunc {
	return CloneableFunc{
		out: out,
		in:  in,
		cloneFn: func() interface{} {
			CloneMap(in, out)
			return nil
		},
	}
}

func CloneMap(in, out interface{}) {
	buf := new(bytes.Buffer)
	_ = gob.NewEncoder(buf).Encode(in)
	_ = gob.NewDecoder(buf).Decode(out)
}
