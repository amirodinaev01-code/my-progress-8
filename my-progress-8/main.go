package main

import (
	"fmt"
	"sync"
)

type Payment interface {
	Wallet()
}

type Order struct {
	Name  string
	Price float64
	Pay   float64
}

func (p *Order) Wallet() {
	fmt.Printf("Payment successful for %s. Enjoy!\n", p.Name)
}

var wg sync.WaitGroup

type Validator func(o *Order) (bool, error)

func main() {
	order := []*Order{
		{Name: "Barbeq", Price: 120.0, Pay: 50.0},
		{Name: "Pizza", Price: 80.0, Pay: 30.0},
		{Name: "Burger", Price: 50.0, Pay: 100.0},
	}

	examination := []Validator{
		func(o *Order) (bool, error) {
			if o.Price > o.Pay {

				return false, fmt.Errorf("Insufficient funds on balance")
			}
			return true, nil
		},
	}
	var wg sync.WaitGroup

	reportChan := make(chan error, len(order))

	for _, i := range order {
		wg.Add(1)

		go func(currentOrder *Order) {
			defer wg.Done()

			ok, err := LogikDish(currentOrder, examination)
			if !ok {
				reportChan <- fmt.Errorf("order %s rejected: %v", currentOrder.Name, err)
			} else {
				reportChan <- fmt.Errorf("PAID:  %s", currentOrder.Name)
				PayService(currentOrder)
			}

		}(i)
	}

	wg.Wait()
	close(reportChan)

	for e := range reportChan {
		fmt.Printf("Result: %v\n", e)
	}
}

func PayService(p Payment) {
	if p == nil {
		fmt.Println("Error: Payment impossible, no data")
		return
	}

	p.Wallet()
}

func LogikDish(o *Order, check []Validator) (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from: %v\n", r)
		}
	}()

	for _, checkFunc := range check {
		valid, err := checkFunc(o)
		if !valid {
			return false, err
		}
	}

	ok = true
	return ok, err
}