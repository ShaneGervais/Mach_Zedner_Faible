package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/unit/constant"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const I = complex(0, 1)

type Pointeur struct {
	Time     []float64
	Function []complex128
	Delay    float64
	//Intensity []float64
	//Polarized []float64
	//AngleRad  float64
}

type PolarizationState struct {
	Horizontal complex128
	Vertical   complex128
}

func linspace(start, stop float64, num int) []float64 {
	result := make([]float64, num)
	step := (stop - start) / float64(num-1)

	for i := 0; i < num; i++ {
		result[i] = start + float64(i)*step
	}

	return result
}

func generatePointerState(t []float64, pulseWidth float64, wavelength float64, z float64) []complex128 {
	amplitude := 1 / math.Sqrt(math.Sqrt(2*math.Pi)*pulseWidth)
	pointeur := make([]complex128, len(t))
	w := 2 * math.Pi * float64(constant.LightSpeedInVacuum) / wavelength
	k := (2 * math.Pi) / wavelength

	for j := 0; j < len(t); j++ {
		pointeur[j] = complex(amplitude*math.Exp(-math.Pow(t[j]/(2*pulseWidth), 2))*math.Cos((k*z-w*t[j])), amplitude*math.Exp(-math.Pow(t[j]/(2*pulseWidth), 2))*math.Sin((k*z-w*t[j])))
	}
	return pointeur
}

func coupleInitialState(state PolarizationState, pointeur_function []complex128, degree_of_freedom []float64) []PolarizationState {
	coupled_state := make([]PolarizationState, len(degree_of_freedom))

	for i := 0; i < len(degree_of_freedom); i++ {
		coupled_state[i].Horizontal = pointeur_function[i] * state.Horizontal
		coupled_state[i].Vertical = pointeur_function[i] * state.Vertical
	}

	return coupled_state
}

func intensity_HV_Profile(coupled_polarisation_state []PolarizationState, degree_of_freedom []float64) []float64 {
	intensity_profile := make([]float64, len(degree_of_freedom))
	for i := 0; i < len(degree_of_freedom); i++ {
		intensity_H := math.Pow(cmplx.Abs(coupled_polarisation_state[i].Horizontal), 2)
		intensity_V := math.Pow(cmplx.Abs(coupled_polarisation_state[i].Vertical), 2)
		intensity_profile[i] = intensity_H + intensity_V
	}

	return intensity_profile
}

func plotIntensity(t []float64, intensity []float64) {
	// Create a new plot
	p := plot.New()

	// Create a new scatter plotter
	points := make(plotter.XYs, len(t))
	for i := range points {
		points[i].X = t[i]
		points[i].Y = intensity[i]
	}
	s, err := plotter.NewScatter(points)
	if err != nil {
		fmt.Println("Error creating scatter plotter:", err)
		return
	}

	// Add the scatter plotter to the plot
	p.Add(s)

	// Set the plot title and labels
	p.Title.Text = "Intensity Plot"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Intensity"

	// Save the plot to a PNG file
	err = p.Save(6*vg.Inch, 4*vg.Inch, "intensity_plot.png")
	if err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

	fmt.Println("Intensity plot saved as intensity_plot.png")
}

func applyTemporalDelay(amplitude []complex128, t []float64, delay float64) []complex128 {
	delayedAmplitude := make([]complex128, len(amplitude))
	for i := 0; i < len(amplitude); i++ {
		delayedAmplitude[i] = amplitude[i] * complex(math.Exp(-2*math.Pi*delay*t[i]), 0) // Apply delay
	}
	return delayedAmplitude
}

func applyPolarizer(interference []float64, angleRad float64) []float64 {
	polarized := make([]float64, len(interference))
	for i := 0; i < len(interference); i++ {
		polarized[i] = interference[i] * math.Cos(angleRad) // Apply polarizer transformation
	}
	return polarized
}

func calculateIntensity(waveform []float64) []float64 {
	intensity := make([]float64, len(waveform))
	for i := 0; i < len(waveform); i++ {
		intensity[i] = math.Pow(waveform[i], 2) // Square the waveform
	}
	return intensity
}

func calculateAverage(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func main() {

	// Simulation parameters
	wavelength := 640e-9 // Wavelength of the pulsed laser diode in meters
	//freq := float64(constant.LightSpeedInVacuum) / wavelength
	//w := 2 * math.Pi * freq
	pulseWidth := 10e-9 // Width of the pulse in the temporal frame in seconds (sigma_t)
	delay := 1e-9       // Temporal delay in seconds for one of the polarization components
	//polarizerAngle := math.Pi / 4 // Orientation of the polarizer in degrees
	preselection_angle := 0
	// Generate time array
	timePoints := 1000
	startTime := -pulseWidth
	stopTime := pulseWidth
	t := linspace(startTime, stopTime, timePoints)

	delta := math.Pi / 2 //represente une lame quart d'onde (peut etre negative)
	phi_x := 0
	phi_y := delta
	// Create Polarization struct
	pointeur := Pointeur{
		Time:  t,
		Delay: delay,
	}

	polarisation := PolarizationState{
		Horizontal: complex(math.Cos(float64(preselection_angle))*math.Cos(float64(phi_x)), math.Cos(float64(preselection_angle))*math.Sin(float64(phi_x))),
		Vertical:   complex(math.Sin(float64(preselection_angle))*math.Cos(float64(phi_y)), math.Sin(float64(preselection_angle))*math.Sin(float64(phi_x))),
	}

	vect_jones_0 := []complex128{polarisation.Horizontal, polarisation.Vertical}
	fmt.Println("Vecteur de jones initial: ", vect_jones_0)

	// Generate Gaussian temporal profile
	pointeur.Function = generatePointerState(pointeur.Time, pulseWidth, wavelength, 0)
	fmt.Println("Pointeur amplitude:", pointeur.Function[0])

	initial_state := coupleInitialState(polarisation, pointeur.Function, pointeur.Time)
	fmt.Println("Initial state at index", 0, ":", initial_state[0])
	fmt.Println("Initial horizontal state at index", 0, ":", initial_state[0].Horizontal)
	fmt.Println("Initial vertical state at index", 0, ":", initial_state[0].Vertical)

	intensity := intensity_HV_Profile(initial_state, pointeur.Time)
	fmt.Println("Intensity profile:", intensity)

	plotIntensity(pointeur.Time, intensity)

	// Apply temporal delay to one polarization component
	//delayedAmplitude := applyTemporalDelay(pointeur.Amplitude, pointeur.Time, pointeur.Delay)

	/*
		// Combine delayed and non-delayed components for interference
		interference := make([]float64, len(polarization.Time))
		for i := 0; i < len(interference); i++ {
			interference[i] = polarization.Amplitude[i] + delayedAmplitude[i] // Superposition of waveforms
		}

		// Apply post-selection with polarizer
		polarization.Polarized = applyPolarizer(interference, polarization.AngleRad)

		// Calculate intensity
		polarization.Intensity = calculateIntensity(polarization.Polarized)

		// Calculate the real part of the weak value based on intensity
		realWeakValue := calculateAverage(polarization.Intensity)

		// Display results
		fmt.Println("Real part of the weak value:", realWeakValue)
	*/
}
