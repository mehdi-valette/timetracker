package entity

type Task struct {
	id   DbId
	name Name
}

func CreateTask(id DbId, name string) Task {
	return Task{
		id:   id,
		name: CreateName(name),
	}
}

func (t Task) Rename(newName string) Task {
	t.name = CreateName(newName)

	return t
}
