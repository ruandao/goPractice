package SafeSlice

type SafeSlice interface {
	Append(interface{})
	At(int) interface{}
	Close() []interface{}
	Delete(int)
	Len() int
	Update(int, UpdateFunc)
}

type UpdateFunc func(value interface{}, exist bool) interface{}

type command int

const (
	append command = iota
	index
	end
	remove
	length
	update
)

type safeSlice chan commandData

type commandData struct {
	action 		command
	key			int
	value		interface{}
	result		chan interface{}
	updateFunc  UpdateFunc
}

func (ss safeSlice)Append(v interface{})  {
	ss <- commandData{action:append, value:v}
}

func (ss safeSlice)At(idx int) interface{} {
	result := make(chan interface{})
	ss <- commandData{action:index, key:idx, result:result}
	v := <- result
	return v
}
func (ss safeSlice)Close() []interface{} {
	result := make(chan interface{})
	ss <- commandData{action:end, result:result}
	v := <- result
	return v
}
func (ss safeSlice)Delete(idx int) {
	ss <- commandData{action:remove, key:idx}
}
func (ss safeSlice)Update(idx int, f UpdateFunc) {
	ss <- commandData{action:update, updateFunc:f}
}

func (ss safeSlice)Len() int {
	result := make(chan interface{})
	ss <- commandData{action:length, result:result}
	v := ( <- result).(int)
	return v
}

func New() SafeSlice {
	ss := make(safeSlice)
	go ss.run()
	return ss
}
func (ss safeSlice)run() {
	s := []interface{}{}
	for {
		cd := <- ss
		switch cd.action {
		case append:
			val := cd.value
			s = append(s, val)
		case index:
			index := cd.key
			val := s[index]
			cd.result <- val
		case end:
			close(ss)
			cd.result <- s
			break
		case remove:
			index := cd.key
			s = append(s[:index], s[index+1:]...)
		case length:
			cd.result <- len(s)
		case update:
			index := cd.key
			s[index] = cd.updateFunc(s[index], len(s) > index)
		}
	}
}
