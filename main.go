package main

import (
	"errors"
	"fmt"
	"systems-diagrams/diagrams"
	"systems-diagrams/systems"
	"time"
)

func main() {
	systemUpdates := make(chan systems.System)

	systemsService := systems.NewSystemsService(systemUpdates)
	diagramService := diagrams.NewDiagramService(systemUpdates)

	go diagramService.ListenForUpdates()

	system := systems.System{Id: 1, Name: "System 1", Description: "Description 1"}
	systemsService.AddSystem(system)

	element1 := diagrams.DiagramElement{Id: 1, Name: "Element 1", Description: "Element Description 1", LinkedSystemID: &system.Id}
	element2 := diagrams.DiagramElement{Id: 2, Name: "Element 2", Description: "Element Description 2", LinkedSystemID: &system.Id}
	diagramService.AddElement(element1)
	diagramService.AddElement(element2)

	element3 := diagrams.DiagramElement{Id: 3, Name: "Element 3", Description: "Element Description 3", LinkedSystemID: nil}
	element4 := diagrams.DiagramElement{Id: 4, Name: "Element 4", Description: "Element Description 4", LinkedSystemID: nil}
	diagramService.AddElement(element3)
	diagramService.AddElement(element4)

	err := diagramService.EditElement(1, "Updated Element 1", "Updated Element Description 1")
	if err != nil {
		if !errors.Is(err, diagrams.ErrCantEditLinkedElement) {
			fmt.Println("Edit of the diagram element should not be allowed")
		}
	}

	err = diagramService.EditElement(3, "Updated Element 3", "Updated Element Description 3")
	if err != nil {
		fmt.Println("Edit of the diagram element should be allowed")
	}

	systemsService.UpdateSystem(1, "Updated System 1", "Updated Description 1")

	time.Sleep(50 * time.Millisecond)

	for _, element := range diagramService.DiagramElements {
		if element.LinkedSystemID != nil && *element.LinkedSystemID == 1 {
			if element.Name != "Updated System 1" || element.Description != "Updated Description 1" {
				fmt.Println("Linked diagram element should have been updated")
			} else {
				fmt.Printf("Linked diagram element updated correctly: %+v\n", element)
			}
		} else {
			if element.Name == "Updated System 1" || element.Description == "Updated Description 1" {
				fmt.Println("Unlinked diagram element should not have been updated")
			} else {
				fmt.Printf("Unlinked diagram element remains unchanged: %+v\n", element)
			}
		}
	}
}
