package typ

import "fmt"

func Switch[TVal, A, B any](v *TVal, a func(*A), b func(*B)) {
	switch v := any(v).(type) {
	case *A:
		a(v)
	case *B:
		b(v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch3[TVal, A, B, C any](v *TVal, a func(*A), b func(*B), c func(*C)) {
	switch v := any(v).(type) {
	case *A:
		a(v)
	case *B:
		b(v)
	case *C:
		c(v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch3UV[TVal, A, B, C any](u, v *TVal, a func(*A, *A), b func(*B, *B), c func(*C, *C)) {
	switch v := any(v).(type) {
	case *A:
		a(any(u).(*A), v)
	case *B:
		b(any(u).(*B), v)
	case *C:
		c(any(u).(*C), v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch4UV[TVal, A, B, C, D any](u, v *TVal, a func(*A, *A), b func(*B, *B), c func(*C, *C), d func(*D, *D)) {
	switch v := any(v).(type) {
	case *A:
		a(any(u).(*A), v)
	case *B:
		b(any(u).(*B), v)
	case *C:
		c(any(u).(*C), v)
	case *D:
		d(any(u).(*D), v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch4[TVal, A, B, C, D any](v *TVal, a func(*A), b func(*B), c func(*C), d func(*D)) {
	switch v := any(v).(type) {
	case *A:
		a(v)
	case *B:
		b(v)
	case *C:
		c(v)
	case *D:
		d(v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch5[TVal, A, B, C, D, E any](v *TVal, a func(*A), b func(*B), c func(*C), d func(*D), e func(*E)) {
	switch v := any(v).(type) {
	case *A:
		a(v)
	case *B:
		b(v)
	case *C:
		c(v)
	case *D:
		d(v)
	case *E:
		e(v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch6[TVal, A, B, C, D, E, F any](v *TVal, a func(*A), b func(*B), c func(*C), d func(*D), e func(*E), f func(*F)) {
	switch v := any(v).(type) {
	case *A:
		a(v)
	case *B:
		b(v)
	case *C:
		c(v)
	case *D:
		d(v)
	case *E:
		e(v)
	case *F:
		f(v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch7[TVal, A, B, C, D, E, F, G any](v *TVal, a func(*A), b func(*B), c func(*C), d func(*D), e func(*E), f func(*F), g func(*G)) {
	switch v := any(v).(type) {
	case *A:
		a(v)
	case *B:
		b(v)
	case *C:
		c(v)
	case *D:
		d(v)
	case *E:
		e(v)
	case *F:
		f(v)
	case *G:
		g(v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}

func Switch8[TVal, A, B, C, D, E, F, G, H any](v *TVal, a func(*A), b func(*B), c func(*C), d func(*D), e func(*E), f func(*F), g func(*G), h func(*H)) {
	switch v := any(v).(type) {
	case *A:
		a(v)
	case *B:
		b(v)
	case *C:
		c(v)
	case *D:
		d(v)
	case *E:
		e(v)
	case *F:
		f(v)
	case *G:
		g(v)
	case *H:
		h(v)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}
