package rest

type Parser interface {
	Parse() ContentType
}

type EncoderRegistry map[ContentType]MarshallFunc

func (r EncoderRegistry) Encode(object any, content Parser) ([]byte, error) {
	if object == nil {
		return nil, nil
	}
	f, ok := r[content.Parse()]
	if !ok {
		return []byte{}, ErrMarshallerFuncNotFound
	}
	return f(object)
}

func (r EncoderRegistry) Set(contentType ContentType, f MarshallFunc) {
	r[contentType] = f
}

func (r EncoderRegistry) Clone() EncoderRegistry {
	result := make(EncoderRegistry, len(r))
	for k, v := range r {
		result[k] = v
	}
	return result
}

type DecoderRegistry map[ContentType]UnmarshallFunc

func (r DecoderRegistry) Decode(data []byte, object any, content Parser) error {
	f, ok := r[content.Parse()]
	if !ok {
		return ErrUnmarshalerFuncNotFound
	}
	return f(data, object)
}

func (r DecoderRegistry) Set(contentType ContentType, f UnmarshallFunc) {
	r[contentType] = f
}

func (r DecoderRegistry) Clone() DecoderRegistry {
	result := make(DecoderRegistry, len(r))
	for k, v := range r {
		result[k] = v
	}
	return result
}

type Encoder struct {
	Registry interface {
		Encode(object any, content Parser) ([]byte, error)
		Set(contentType ContentType, f MarshallFunc)
		Clone() EncoderRegistry
	}
}

func (e Encoder) Encode(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
	registry := e.Registry
	if marshalFunc != nil {
		registry = registry.Clone()
		registry.Set(contentType, marshalFunc)
	}
	return registry.Encode(value, contentType)
}

type Decoder[T any] struct {
	Registry interface {
		Decode(data []byte, object any, content Parser) error
		Set(contentType ContentType, f UnmarshallFunc)
		Clone() DecoderRegistry
	}
}

func (d Decoder[T]) Decode(data []byte, val *T, contentType ContentType, unmarshallFunc UnmarshallFunc) error {
	plain := string(data)
	switch any(val).(type) {
	case *string:
		*val = *any(&plain).(*T)
	default:
		registry := d.Registry
		if unmarshallFunc != nil {
			registry = registry.Clone()
			registry.Set(contentType, unmarshallFunc)
		}
		return d.Registry.Decode(data, val, contentType)
	}
	return nil
}
