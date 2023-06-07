package nebula_amqp_lib

type Driver int32

const (
	Rabbit Driver = 0
)

func (driver Driver) String() string {
	switch driver {
	case Rabbit:
		return "rabbit"
	default:
		return "rabbit"
	}
}
