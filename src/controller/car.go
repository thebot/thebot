package main

import (
	"log"
	"sync"
)

const (
	Servo = 0x53
	Motor = 0x4D
	Reset = 0x52
)

var bus *I2CBus

func init() {
	var err error

	bus, err = Bus(1)
	if err != nil {
		panic(err)
	}
}

type Car struct {
	addr byte

	curAngle, curSpeed int

	mu *sync.RWMutex
}

func NewCar(addr byte) *Car {
	return &Car{addr: addr, mu: &sync.RWMutex{}}
}

func (c *Car) Turn(angle int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmd := []byte{Servo, byte(angle)}
	if err := bus.WriteBytes(c.addr, cmd); err != nil {
		return err
	}

	log.Printf("Set the angle to %v", angle)

	c.curAngle = angle

	return nil
}

func (c *Car) Speed(speed int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmd := []byte{Motor, byte(speed)}
	if err := bus.WriteBytes(c.addr, cmd); err != nil {
		return err
	}

	log.Printf("Set the speed to %v", speed)

	c.curSpeed = speed

	return nil
}

func (c *Car) Orientation() (speed, angle int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.curSpeed, c.curAngle
}

func (c *Car) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := bus.WriteByte(c.addr, Reset); err != nil {
		return err
	}

	log.Print("Reset the device")

	c.curAngle, c.curSpeed = 0, 0

	return nil
}