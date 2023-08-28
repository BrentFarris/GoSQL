package db

type Join struct {
	leftTable  string
	leftField  string
	rightTable string
	rightField string
}

func (j *Join) On(leftField, rightField string) {
	j.leftField = leftField
	j.rightField = rightField
}
