package main

import (
	"io"
	"log"
	nhlapi "nhlapi/nhlApi"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	start := time.Now()

	roasterFile, err := os.OpenFile("roasters.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Error in opening file %v", err)
	}
	defer roasterFile.Close()

	wrt := io.MultiWriter(roasterFile)

	log.SetOutput(wrt)

	log.Println("start time: ", start.String())

	teams, err := nhlapi.GetAllTeams()

	results := make(chan []nhlapi.Roster)

	wg.Add(len(teams))

	for _, team := range teams {
		go doTaskRoaster(team, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	display(results)

	log.Printf("took %v", time.Now().Sub(start).String())

}

func doTaskRoaster(team nhlapi.Team, results chan []nhlapi.Roster) chan []nhlapi.Roster {

	roster, err := nhlapi.GetRosters(team.ID)
	if err != nil {
		log.Fatalf("error getting roster: %v", err)
	}

	results <- roster

	wg.Done()

	return results
}

func display(results chan []nhlapi.Roster) {
	for r := range results {
		for _, ros := range r {
			log.Println("----------------------")
			log.Printf("ID: %d\n", ros.Person.ID)
			log.Printf("Name: %s\n", ros.Person.FullName)
			log.Printf("Position: %s\n", ros.Position.Abbreviation)
			log.Printf("Jersey: %s\n", ros.JerseyNumber)
			log.Println("----------------------")
		}
	}
}
