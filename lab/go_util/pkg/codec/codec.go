package codec

type Codec interface {
	Marshal(o interface{}) (b []byte, err error)
	Unmarshal(b []byte, o interface{}) (err error)
}
