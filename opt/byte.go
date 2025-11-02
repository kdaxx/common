package opt

func Has(target byte, has byte) bool {
	return (target & has) != 0
}

func Set(target byte, set byte) byte {
	target |= set
	return target
}
func Unset(target byte, unset byte) byte {
	target &= ^unset
	return target
}
