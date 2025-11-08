package opt

func Has(target byte, has byte) bool {
	return (target & has) != 0
}

func SetOpt(target byte, opt byte) byte {
	target |= opt
	return target
}
func Unset(target byte, unset byte) byte {
	target &= ^unset
	return target
}

func SetOpts(target byte, opts ...byte) byte {
	for i := 0; i < len(opts); i++ {
		target |= opts[i]
	}
	return target
}

func Opts(opts ...byte) byte {
	opt := byte(0)
	for i := 0; i < len(opts); i++ {
		opt |= opts[i]
	}
	return opt
}
