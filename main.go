package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/unit/constant"
)

type Polarization struct {
	Time      []float64
	Amplitude []float64
	Delay     float64
	Intensity []float64
	Polarized []float64
	AngleRad  float64
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
	amplitude := make([]complex128, len(t))
	w := 2 * math.Pi * float64(constant.LightSpeedInVacuum) / wavelength
	k := (2 * math.Pi) / wavelength
	imaginaryUnit := complex(0, 1) // Imaginary unit (complex128)
	for j := 0; j < len(t); j++ {
		exponent := complex(-math.Pow(t[j]/(2*pulseWidth), 2), 0)                              // Convert the exponent to complex128
		amplitude[j] = cmplx.Exp(exponent) * cmplx.Exp(complex128(imaginaryUnit*(k*z-w*t[j]))) // Gaussian function
	}
	return amplitude
}

func applyTemporalDelay(amplitude []float64, t []float64, delay float64) []float64 {
	delayedAmplitude := make([]float64, len(amplitude))
	for i := 0; i < len(amplitude); i++ {
		delayedAmplitude[i] = amplitude[i] * math.Exp(-2*math.Pi*delay*t[i]) // Apply delay
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
	freq := float64(constant.LightSpeedInVacuum) / wavelength
	w := 2 * math.Pi * freq
	pulseWidth := 10e-9    // Width of the pulse in the temporal frame in seconds (sigma_t)
	delay := 1e-9          // Temporal delay in seconds for one of the polarization components
	polarizerAngle := 45.0 // Orientation of the polarizer in degrees

	// Generate time array
	timePoints := 1000
	startTime := 0.0
	stopTime := pulseWidth
	t := linspace(startTime, stopTime, timePoints)

	// Create Polarization struct
	polarization := Polarization{
		Time:     t,
		Delay:    delay,
		AngleRad: polarizerAngle * math.Pi / 180.0,
	}

	// Generate Gaussian temporal profile
	polarization.Amplitude = generateGaussianProfile(polarization.Time, pulseWidth)

	// Apply temporal delay to one polarization component
	delayedAmplitude := applyTemporalDelay(polarization.Amplitude, polarization.Time, polarization.Delay)

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
}
