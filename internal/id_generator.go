package internal

type Generator struct {
	currId int
}

func (g *Generator) generateId() int {
	g.currId += 1
	return g.currId
}