package helper

type ErrorHelper interface {
	CheckErr(err error)
	BehaviorOnFatal(f func(string, int))
}
