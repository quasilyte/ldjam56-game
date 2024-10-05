package gcombat

type Card struct {
	Kind CardKind
}

type CardKind int

const (
	CardUnknown CardKind = iota
)
