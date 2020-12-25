package banlancer

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/micro/go-micro/v2/registry"
)

func TestSmoothWeight(t *testing.T) {
	var services []*registry.Service
	var weights [4]string = [4]string{"1", "1", "5", "3"}

	services = make([]*registry.Service, 1)

	for index := range services {
		service := &registry.Service{
			Name: "TestService" + strconv.Itoa(index),
		}
		services[index] = service

		service.Nodes = make([]*registry.Node, 4)
		nodes := service.Nodes

		for index := range nodes {
			node := &registry.Node{
				Id:       service.Name + "--" + strconv.Itoa(index),
				Metadata: map[string]string{"weights": weights[index]},
			}
			nodes[index] = node
		}
	}

	//測試時必須把Chan給註釋掉
	Banlancer := SmoothWeight()

	next := Banlancer(services)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			node, err := next()
			if err != nil {
				t.Error(err.Error())
			}
			fmt.Println(node)
		}()
	}
	wg.Wait()
}
