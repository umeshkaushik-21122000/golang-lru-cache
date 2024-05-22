package types

type KeyValue struct {
	Key        int `json:"key,omitempty"`
	Value      int `json:"value,omitempty"`
	Expiration int `json:"expiration,omitempty"`
}
