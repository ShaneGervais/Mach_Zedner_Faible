package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"os"
	"strconv"

	"gonum.org/v1/gonum/unit/constant"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const C = float64(constant.LightSpeedInVacuum)

type Pointeur struct {
	Time     []float64
	Function []complex128
	Delay    float64
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

/*
	func plotIntensity(t []float64, intensity []float64, save_as string) {
		p := plot.New()

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

		p.Add(s)

		p.Title.Text = "Intensity Plot"
		p.X.Label.Text = "Time"
		p.Y.Label.Text = "Intensity"

		err = p.Save(6*vg.Inch, 4*vg.Inch, save_as+".png")
		if err != nil {
			fmt.Println("Error saving plot:", err)
			return
		}

		fmt.Println("Intensity plot saved as intensity_plot.png")
	}
*/
func generatePointerState(t []float64, pulseWidth float64, wavelength float64, z float64) []complex128 {
	amplitude := 1 / math.Sqrt(math.Sqrt(2*math.Pi)*pulseWidth)
	pointeur := make([]complex128, len(t))
	w := 2 * math.Pi * C / wavelength
	k := (2 * math.Pi) / wavelength

	for j := 0; j < len(t); j++ {
		pointeur[j] = complex(amplitude*math.Exp(-math.Pow((t[j]-z/C)/(2*pulseWidth), 2))*math.Cos((k*z-w*t[j])), amplitude*math.Exp(-math.Pow((t[j]-z/C)/(2*pulseWidth), 2))*math.Sin((k*z-w*t[j])))
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

func mach_zedner_interference(coupled_polarisation_state []PolarizationState, degree_of_freedom []float64, delay float64, delayed_part string) []PolarizationState {
	partie_H := make([]complex128, len(degree_of_freedom))
	partie_V := make([]complex128, len(degree_of_freedom))
	for i := 0; i < len(degree_of_freedom); i++ {
		partie_H[i] = coupled_polarisation_state[i].Horizontal
		partie_V[i] = coupled_polarisation_state[i].Vertical
	}

	if delayed_part == "H" || delayed_part == "0" || delayed_part == "h" {
		H_delta := interaction_operator(partie_H, degree_of_freedom, delay) //weakly interacted function
		//fmt.Println("h delta", H_delta[0])
		interference := make([]complex128, len(degree_of_freedom))
		for i := 0; i < len(interference); i++ {
			//interference[i] = H_delta[i] + partie_V[i] // Superposition of waveforms
			coupled_polarisation_state[i].Horizontal = H_delta[i]
			coupled_polarisation_state[i].Vertical = partie_V[i]
		}

	} else if delayed_part == "V" || delayed_part == "0" || delayed_part == "v" {
		V_delta := interaction_operator(partie_H, degree_of_freedom, delay) //weakly interacted function
		interference := make([]complex128, len(degree_of_freedom))
		for i := 0; i < len(interference); i++ {
			//interference[i] = partie_H[i] + V_delta[i] // Superposition of waveforms
			coupled_polarisation_state[i].Horizontal = partie_H[i]
			coupled_polarisation_state[i].Vertical = V_delta[i]
		}
	} else {
		panic("Invalid value for delayed_part. Expected 'H', '0', 'V' for horizontal polarisation or '1', 'h', or 'v' for vertical polarisation.")
	}
	return coupled_polarisation_state
}

func intensity_HV_Profile(coupled_polarisation_state []PolarizationState, degree_of_freedom []float64) []float64 {
	intensity_profile := make([]float64, len(degree_of_freedom))
	for i := 0; i < len(degree_of_freedom); i++ {
		intensity_H := math.Pow(cmplx.Abs(coupled_polarisation_state[i].Horizontal), 2)
		intensity_V := math.Pow(cmplx.Abs(coupled_polarisation_state[i].Vertical), 2)
		intensity_profile[i] = (1 / math.Sqrt(2)) * (intensity_H + intensity_V)
	}

	return intensity_profile
}

func interaction_operator(amplitude []complex128, t []float64, delay float64) []complex128 {
	delayedAmplitude := make([]complex128, len(amplitude))
	for i := 0; i < len(amplitude); i++ {
		//fmt.Println("weakly by", complex(math.Exp(-2*math.Pi*delay*t[i]), 0))
		delayedAmplitude[i] = amplitude[i] * complex(math.Exp(-2*math.Pi*delay*t[i]), 0) // Apply delay
	}
	return delayedAmplitude
}

func apply_postselection(intermediate_state []PolarizationState, postselected_angle float64, phi_x float64, phi_y float64, degree_of_freedom []float64) []PolarizationState {
	polarisation := PolarizationState{
		Horizontal: complex(math.Cos(float64(postselected_angle))*math.Cos(float64(phi_x)), math.Cos(float64(postselected_angle))*math.Sin(float64(phi_x))),
		Vertical:   complex(math.Sin(float64(postselected_angle))*math.Cos(float64(phi_y)), math.Sin(float64(postselected_angle))*math.Sin(float64(phi_y))),
	}

	fmt.Println("horizontal polarisation:", polarisation.Horizontal)
	fmt.Println("vertical polarisation:", polarisation.Vertical)
	fmt.Println("horizontal conj polarisation:", cmplx.Conj(polarisation.Horizontal))
	fmt.Println("vertical conj polarisation:", cmplx.Conj(polarisation.Vertical))

	postselected_state := make([]PolarizationState, len(degree_of_freedom))
	for i := 0; i < len(degree_of_freedom); i++ {
		postselected_state[i].Horizontal = intermediate_state[i].Horizontal * cmplx.Conj(polarisation.Horizontal)
		postselected_state[i].Vertical = intermediate_state[i].Vertical * cmplx.Conj(polarisation.Vertical)
	}
	return postselected_state
}

/*
func degree_of_coherence(postselected_state []PolarizationState) []complex128 {
	g_1 := make([]complex128, len(postselected_state))

	for i := 0; i < len(postselected_state); i++ {
		E_1 := postselected_state[i].Horizontal
		E_2 := postselected_state[i].Vertical

		//fmt.Println("E_1", E_1)
		//fmt.Println("E_2", E_2)

		meanProduct := E_2 * cmplx.Conj(E_1)

		magnitude := cmplx.Abs(meanProduct)

		angle := cmplx.Phase(meanProduct)

		degree := cmplx.Rect(magnitude, angle)

		g_1[i] = degree
	}

	return g_1
}*/

func degree_of_coherence(postselected_state []PolarizationState) []complex128 {
	g_1 := make([]complex128, len(postselected_state))

	maxMagnitude := 0.0

	// Calculate the complex product and find the maximum magnitude
	for i := 0; i < len(postselected_state); i++ {
		E_1 := postselected_state[i].Horizontal
		E_2 := postselected_state[i].Vertical

		meanProduct := E_2 * cmplx.Conj(E_1)

		magnitude := cmplx.Abs(meanProduct)
		if magnitude > maxMagnitude {
			maxMagnitude = magnitude
		}

		g_1[i] = meanProduct
	}

	// Normalize the values in g_1
	for i := 0; i < len(g_1); i++ {
		g_1[i] /= complex(maxMagnitude, 0)
	}

	return g_1
}

//---------------------------------------------------------------------------------------------------------------------------------

func main() {

	// Simulation parameters
	wavelength := 640e-9 // Wavelength of the pulsed laser diode in meters
	//freq := float64(constant.LightSpeedInVacuum) / wavelength
	//w := 2 * math.Pi * freq
	pulseWidth := 10e-9 // Width of the pulse in the temporal frame in seconds (sigma_t)
	delay := 1e-3       // Temporal delay in seconds for one of the polarization components
	preselection_angle := math.Pi / 3
	postselection_angle := math.Pi / 4
	// Generate time array
	timePoints := 1000
	startTime := -pulseWidth
	stopTime := pulseWidth
	t := linspace(startTime, stopTime, timePoints)

	delta := math.Pi / 2 //represente une lame quart d'onde (peut etre negative)
	phi_x := 0
	phi_y := delta

	weak_part := "h"

	//---------------------------------------------------------------------------------------------------------------------------------
	pointeur := Pointeur{
		Time:  t,
		Delay: delay,
	}

	polarisation := PolarizationState{
		Horizontal: complex(math.Cos(float64(preselection_angle))*math.Cos(float64(phi_x)), math.Cos(float64(preselection_angle))*math.Sin(float64(phi_x))),
		Vertical:   complex(math.Sin(float64(preselection_angle))*math.Cos(float64(phi_y)), math.Sin(float64(preselection_angle))*math.Sin(float64(phi_y))),
	}

	//---------------------------------------------------------------------------------------------------------------------------------

	vect_jones_0 := []complex128{polarisation.Horizontal, polarisation.Vertical}
	fmt.Println("Vecteur de jones initial: ", vect_jones_0)

	pointeur.Function = generatePointerState(pointeur.Time, pulseWidth, wavelength, 0)
	fmt.Println("Pointeur amplitude:", pointeur.Function[0])

	initial_state := coupleInitialState(polarisation, pointeur.Function, pointeur.Time)
	fmt.Println("Initial state at index", 0, ":", initial_state[0])
	fmt.Println("Initial horizontal state at index", 0, ":", initial_state[0].Horizontal)
	fmt.Println("Initial vertical state at index", 0, ":", initial_state[0].Vertical)

	intensity := intensity_HV_Profile(initial_state, pointeur.Time)
	//fmt.Println("Intensity profile:", intensity)

	//plotIntensity(pointeur.Time, intensity, "initial_intensity")
	write_to_csv(t, intensity, "time", "intensity", "initial_intensity")

	interfered_state := mach_zedner_interference(initial_state, pointeur.Time, delay, weak_part)
	fmt.Println("Intermediate state at index", 0, ":", interfered_state[0])
	fmt.Println("Intermediate horizontal state at index", 0, ":", interfered_state[0].Horizontal)
	fmt.Println("Intermediate vertical state at index", 0, ":", interfered_state[0].Vertical)

	postselected_state := apply_postselection(interfered_state, postselection_angle, 0, 0, pointeur.Time)
	fmt.Println("Postselected state at index", 0, ":", postselected_state[0])
	fmt.Println("Postselected horizontal state at index", 0, ":", postselected_state[0].Horizontal)
	fmt.Println("Postselected vertical state at index", 0, ":", postselected_state[0].Vertical)

	post_intensity_profile := intensity_HV_Profile(postselected_state, pointeur.Time)

	//plotIntensity(pointeur.Time, post_intensity_profile, "post_selected")
	write_to_csv(t, post_intensity_profile, "time", "intensity", "post_initial_intensity")
	//---------------------------------------------------------------------------------------------------------------------------------

	E_1, E_2 := make([]complex128, len(pointeur.Time)), make([]complex128, len(pointeur.Time))

	for i := 0; i < len(pointeur.Time); i++ {
		E_1[i] = postselected_state[i].Horizontal
		E_2[i] = postselected_state[i].Vertical
	}

	g_1 := degree_of_coherence(postselected_state)

	//---------------------------------------------------------------------------------------------------------------------------------
	//plot degree of coherence
	pts := make(plotter.XYs, len(g_1))
	deg_g1 := make([]float64, len(g_1))
	for i, degree := range g_1 {
		//pts[i].X = float64(i)
		pts[i].Y = cmplx.Abs(degree)
		deg_g1[i] = cmplx.Abs(degree)
	}
	write_to_csv(t, deg_g1, "time", "coherence", "degree_of_coherence")

	for i := range t {
		pts[i].X = t[i]
	}

	p := plot.New()

	p.Title.Text = "Degree of Coherence"
	p.X.Label.Text = "Index"
	p.Y.Label.Text = "Magnitude"

	line, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println("Error creating line plot:", err)
		return
	}

	p.Add(line)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "degree_of_coherence.png"); err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}
}

//---------------------------------------------------------------------------------------------------------------------------------

func write_to_csv(x []float64, y []float64, x_name string, y_name string, name_of_file string) {
	// Create a new CSV file
	file, err := os.Create(name_of_file + ".csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{x_name, y_name}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header:", err)
		return
	}

	// Write the data rows
	for i := 0; i < len(x); i++ {
		record := []string{floatToString(x[i]), floatToString(y[i])}
		err = writer.Write(record)
		if err != nil {
			log.Fatal("Error writing CSV record:", err)
		}
	}

	fmt.Println("Data has been written to", name_of_file, ".csv")
}

func floatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
