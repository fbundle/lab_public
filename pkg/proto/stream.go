package proto

import (
	"io"
)

const (
	separator = '\n'
)

func (p *proto) NewEncoder(w io.Writer) Encoder {
	return &encoder{w: w}
}

type Encoder interface {
	Encode(m interface{}) error
}
type encoder struct {
	p *proto
	w io.Writer
}

func (enc *encoder) Encode(m interface{}) error {
	b, err := enc.p.Marshal(m)
	if err != nil {
		return err
	}
	_, err = enc.w.Write(append(b, separator))
	return err
}

func (p *proto) NewDecoder(r io.Reader) Decoder {
	return &decoder{r: r}
}

type Decoder interface {
	Decode() (interface{}, error)
}

type decoder struct {
	p *proto
	r io.Reader
}

func (dec *decoder) Decode() (interface{}, error) {
	b, err := readUntil(dec.r, separator)
	if err != nil {
		return nil, err
	}
	return dec.p.Unmarshal(b)
}
