package rest

type Parser interface {
	Parse() ContentType
}

type EncodeRegistry map[ContentType]MarshallFunc

func (r EncodeRegistry) Encode(object any, content Parser) ([]byte, error) {
	f, ok := r[content.Parse()]
	if !ok {
		return []byte{}, ErrMarshallerFuncNotFound
	}
	return f(object)
}

func (r EncodeRegistry) Set(contentType ContentType, f MarshallFunc) {
	r[contentType] = f
}

func (r EncodeRegistry) Clone() EncodeRegistry {
	result := make(EncodeRegistry, len(r))
	for k, v := range r {
		result[k] = v
	}
	return result
}

type DecodeRegistry map[ContentType]UnmarshallFunc

func (r DecodeRegistry) Decode(data []byte, object any, content Parser) error {
	f, ok := r[content.Parse()]
	if !ok {
		return ErrUnmarshalerFuncNotFound
	}
	return f(data, object)
}

func (r DecodeRegistry) Set(contentType ContentType, f UnmarshallFunc) {
	r[contentType] = f
}

func (r DecodeRegistry) Clone() DecodeRegistry {
	result := make(DecodeRegistry, len(r))
	for k, v := range r {
		result[k] = v
	}
	return result
}
