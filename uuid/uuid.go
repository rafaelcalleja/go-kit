package uuid

type UUID [16]byte

type Uuid interface {
	Parse(s string) (UUID, error)
	String(uuid UUID) string
}
