// Code generated by go generate; DO NOT EDIT.
package test

import (
	"io"

	"github.com/icon-project/goloop/module"
)

type BlockBase struct{}

func (_r *BlockBase) Version() int {
	panic("not implemented")
}

func (_r *BlockBase) ID() []byte {
	panic("not implemented")
}

func (_r *BlockBase) Height() int64 {
	panic("not implemented")
}

func (_r *BlockBase) PrevID() []byte {
	panic("not implemented")
}

func (_r *BlockBase) NextValidators() module.ValidatorList {
	panic("not implemented")
}

func (_r *BlockBase) Votes() module.CommitVoteSet {
	panic("not implemented")
}

func (_r *BlockBase) NormalTransactions() module.TransactionList {
	panic("not implemented")
}

func (_r *BlockBase) PatchTransactions() module.TransactionList {
	panic("not implemented")
}

func (_r *BlockBase) Timestamp() int64 {
	panic("not implemented")
}

func (_r *BlockBase) Proposer() module.Address {
	panic("not implemented")
}

func (_r *BlockBase) LogsBloom() module.LogsBloom {
	panic("not implemented")
}

func (_r *BlockBase) Result() []byte {
	panic("not implemented")
}

func (_r *BlockBase) MarshalHeader(w io.Writer) error {
	panic("not implemented")
}

func (_r *BlockBase) MarshalBody(w io.Writer) error {
	panic("not implemented")
}

func (_r *BlockBase) Marshal(w io.Writer) error {
	panic("not implemented")
}

func (_r *BlockBase) ToJSON(rcpVersion int) (interface{}, error) {
	panic("not implemented")
}
