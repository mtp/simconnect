package simconnect_test

import (
	"log"
	"time"

	sim "github.com/micmonay/simconnect"
)

func connect() *sim.EasySimConnect {
	sc, err := sim.NewEasySimConnect()
	if err != nil {
		panic(err)
	}
	sc.SetLoggerLevel(sim.LogInfo)
	c, err := sc.Connect("MyApp")
	if err != nil {
		panic(err)
	}
	<-c // wait connection confirmation
	return sc
}

// ExampleGetSimVar this example show how to get SimVar with Easysim
func Example_getSimVar() {
	sc := connect()
	cSimVar := sc.ConnectToSimVarObject(
		sim.SimVarPlaneAltitude(),
		sim.SimVarPlaneLatitude(sim.UnitDegrees), // you can force the units
		sim.SimVarPlaneLongitude(),
		sim.SimVarIndicatedAltitude(),
		sim.SimVarAutopilotAltitudeLockVar(),
		sim.SimVarAutopilotMaster(),
	)
	for i := 0; i < 1; i++ {
		result := <-cSimVar
		for _, simVar := range result {
			f, err := simVar.GetFloat64()
			if err != nil {
				panic(err)
			}
			log.Printf("%#v\n", f)
		}

	}
	<-sc.Close() // wait close confirmation
	// Output:

}

func Example_getSimVarWithIndex() {
	sc := connect()
	cSimVar := sc.ConnectToSimVarObject(
		sim.SimVarGeneralEngRpm(1),
		sim.SimVarTransponderCode(1),
	)
	for i := 0; i < 1; i++ {
		result := <-cSimVar
		for _, simVar := range result {

			if simVar.Name == sim.SimVarTransponderCode().Name {
				i, err := simVar.GetInt()
				if err != nil {
					panic(err)
				}
				log.Printf("%s : %x\n", simVar.Name, i)
			} else {
				f, err := simVar.GetFloat64()
				if err != nil {
					panic(err)
				}
				log.Printf("%s : %f\n", simVar.Name, f)
			}
		}

	}
	<-sc.Close() // wait close confirmation
	// Output:
}

//
func Example_setSimVar() {
	sc := connect()
	newalt := sim.SimVarPlaneAltitude()
	newalt.SetFloat64(6000.0)
	sc.SetSimObject(newalt)
	time.Sleep(1000 * time.Millisecond)
	<-sc.Close() // wait close confirmation
	// NOEXEC Output:
}

func Example_getLatLonAlt() {
	sc := connect()
	cSimVar := sc.ConnectToSimVarObject(
		sim.SimVarStructLatlonalt(),
	)
	for i := 0; i < 1; i++ {
		result := <-cSimVar
		for _, simVar := range result {
			latlonalt, err := simVar.GetDataLatLonAlt()
			if err != nil {
				panic(err)
			}
			log.Printf("%s : %#v\nIn Feet %#v\n", simVar.Name, latlonalt, latlonalt.GetFeets())
		}

	}
	<-sc.Close() // wait close confirmation
	// Output:
}

func Example_getXYZ() {
	sc := connect()
	cSimVar := sc.ConnectToSimVarObject(
		sim.SimVarEyepointPosition(),
	)
	for i := 0; i < 1; i++ {
		result := <-cSimVar
		for _, simVar := range result {
			xyz, err := simVar.GetDataXYZ()
			if err != nil {
				panic(err)
			}
			log.Printf("%s : %#v\n", simVar.Name, xyz)
		}

	}
	<-sc.Close() // wait close confirmation
	// Output:
}

func Example_getString() {
	sc := connect()
	cSimVar := sc.ConnectToSimVarObject(
		sim.SimVarTitle(),
		sim.SimVarCategory(),
	)
	for i := 0; i < 1; i++ {
		result := <-cSimVar
		for _, simVar := range result {
			str := simVar.GetString()
			log.Printf("%s : %#v\n", simVar.Name, str)
		}

	}
	<-sc.Close() // wait close confirmation
	// Output:
}

// Example_showText Actually color no effect in the sim
func Example_showText() {
	sc := connect()
	ch, err := sc.ShowText("Test", 1, sim.SIMCONNECT_TEXT_TYPE_PRINT_GREEN)
	if err != nil {
		panic(err)
	}
	log.Println(<-ch)
	<-sc.Close() // wait close confirmation
	// Output:
}

//Example_simEvent You can wait chan if you will surre the event has finish with succes. If your app finish before all event probably not effect.
func Example_simEvent() {
	sc := connect()
	aileronsSet := sc.NewSimEvent(sim.KeyAxisAileronsSet)
	throttleSet := sc.NewSimEvent(sim.KeyThrottleSet)
	altVarInc := sc.NewSimEvent(sim.KeyApAltVarInc)
	altVarDec := sc.NewSimEvent(sim.KeyApAltVarDec)
	log.Println(<-aileronsSet.RunWithValue(-16383))
	log.Println(<-throttleSet.RunWithValue(16383))
	for i := 0; i < 10; i++ {
		<-altVarInc.Run()
	}
	for i := 0; i < 10; i++ {
		<-altVarDec.Run()
	}
	<-sc.Close() // wait close confirmation
	// Output:
}
