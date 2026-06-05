package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	file, err := os.Create("rpm_data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Physics and state tracking constants
	const (
		idleRPM     = 800.0   // Standard warm engine idle
		shiftRedline = 4200.0  // RPM point where driver shifts gears
		maxGears    = 5
	)

	state := 0 // 0: Idle, 1: Accelerating, 2: Clutch In (Shifting), 3: Cruising/End
	currentRPM := idleRPM
	currentGear := 1
	stateSteps := 0

	// 1,500 entries @ 20ms intervals = 30 seconds of seamless dashboard data
	for i := 0; i < 1500; i++ {
		stateSteps++
		
		// Micro-vibrations from engine harmonics
		noise := (rand.Float64() - 0.5) * 20.0

		switch state {
		case 0: // Phase 1: Hold steady at Idle (150 steps = 3 seconds)
			currentRPM = idleRPM
			if stateSteps > 150 {
				state = 1
				stateSteps = 0
				currentGear = 1
			}

		case 1: // Phase 2: Accelerating through a gear
			// Higher gears accelerate slower due to mechanical gear ratios
			// Gear 1 climbs fast (+45 RPM), Gear 5 climbs slower (+18 RPM)
			accelerationRate := 50.0 - float64(currentGear)*6.0
			currentRPM += accelerationRate

			// If engine hits the shift point redline, step on the clutch
			if currentRPM >= shiftRedline {
				currentRPM = shiftRedline
				state = 2
				stateSteps = 0
			}

		case 2: // Phase 3: Shifting gears (Clutch-in drop)
			// The engine drops 180 RPM every 20ms while the clutch is pressed
			currentRPM -= 180.0

			// Target landing RPM based on gear ratios (higher gears drop less)
			// e.g., Shifting 1->2 drops down to ~1800 RPM. 4->5 drops to ~2600 RPM.
			targetDropRPM := 1200.0 + (float64(currentGear) * 350.0)

			if currentRPM <= targetDropRPM {
				currentRPM = targetDropRPM
				stateSteps = 0
				
				if currentGear < maxGears {
					currentGear++
					state = 1 // Go back to accelerating in the next gear
				} else {
					state = 3 // Finished all 5 gears, transition to cruise
				}
			}

		case 3: // Phase 4: Cruising in top gear, then slowing back down to idle
			if stateSteps < 200 {
				// Cruise at a steady speed (~3000 RPM)
				currentRPM = 2800.0 + (rand.Float64()-0.5)*40.0
			} else {
				// Coast/Brake back down to idle speed
				currentRPM -= 35.0
				if currentRPM <= idleRPM {
					currentRPM = idleRPM
					state = 0 // Reset entire simulation loop
					stateSteps = 0
				}
			}
		}

		// Enforce physical boundaries
		finalRPM := currentRPM + noise
		if finalRPM < idleRPM-50.0 && state != 2 {
			finalRPM = idleRPM + noise
		}

		// Write a clean dataset to file
		fmt.Fprintf(file, "%.2f\n", finalRPM)
	}

	fmt.Println("Successfully generated a realistic 5-gear dataset in rpm_data.txt!")
}

