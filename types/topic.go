package types

type Topic struct {
	Uuid      string      `json:"uuid"`
	Content   *Content    `json:"content,omitempty"`
	Viewport  *Viewport   `json:"viewport,omitempty"`
	Selection *Selections `json:"selection,omitempty"`
}
