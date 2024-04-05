package integral

//
// IACC - Integration of ACCeleration time history
//

func Iacc(ddy []float64, dt float64) ([]float64, []float64) {
	n := len(ddy)
	dy := make([]float64, n)
	y := make([]float64, n)

	dy[0] = 0.0
	y[0] = 0.0
	for m := 1; m < n; m++ {
		dy[m] = dy[m-1] + (ddy[m-1]+ddy[m])*dt/2.0
		y[m] = y[m-1] + dy[m-1]*dt + (2.0*ddy[m-1]+ddy[m])*dt*dt/6.0
	}

	return dy, y
}
