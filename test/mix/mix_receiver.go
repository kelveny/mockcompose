package mix

type mixReceiver struct {
	value string
}

func (v mixReceiver) getValue() string {
	return v.value
}

func (p *mixReceiver) setValue(val string) {
	p.value = val
}

// only callees with the same receiver type (either by-value or by-reference type)
// will be considered as peer callee method
//
// in checkAndSet method, since receiver type of checkAndSet is by-reference,
// only by-reference setValue will be considered as its peer callee method automatically.
//
// in checkAndSetOnTarget, since receiver type of checkAndSetOnTarget is by-value
// only by-value getValue will be considered as its peer callee method automatically.

//go:generate mockcompose -n mix_checkAndSet -c mixReceiver -real checkAndSet,this
func (p *mixReceiver) checkAndSet(s string, val string) {
	if p.getValue() != s {
		p.setValue(val)
	}
}

//go:generate mockcompose -n mix_checkAndSetOnTarget -c mixReceiver -real checkAndSetOnTarget,this
func (v mixReceiver) checkAndSetOnTarget(p *mixReceiver, s string, val string) {
	if v.getValue() != s {

		// use it for code generation validation, has no effect on v
		v.setValue(val)
		p.setValue(val)
	}
}
