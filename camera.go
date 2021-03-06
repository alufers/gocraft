package main

import "github.com/go-gl/mathgl/mgl32"

type CameraMovement int

const (
	MoveForward CameraMovement = iota
	MoveBackward
	MoveLeft
	MoveRight
)

type Camera struct {
	CamPos mgl32.Vec3
	up     mgl32.Vec3
	right  mgl32.Vec3
	front  mgl32.Vec3
	wfront mgl32.Vec3

	rotatex, rotatey float32

	Sens float32

	flying bool
}

func NewCamera(pos mgl32.Vec3) *Camera {
	c := &Camera{
		CamPos:  pos,
		front:   mgl32.Vec3{0, 0, -1},
		rotatey: 0,
		rotatex: -90,
		Sens:    0.14,
		flying:  false,
	}
	c.updateAngles()
	return c
}

func (c *Camera) Matrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.CamPos, c.CamPos.Add(c.front), c.up)
}

func (c *Camera) SetPos(pos mgl32.Vec3) {
	c.CamPos = pos
}

func (c *Camera) Pos() mgl32.Vec3 {
	return c.CamPos
}

func (c *Camera) Front() mgl32.Vec3 {
	return c.front
}

func (c *Camera) FlipFlying() {
	c.flying = !c.flying
}

func (c *Camera) Flying() bool {
	return c.flying
}

func (c *Camera) OnAngleChange(dx, dy float32) {
	if mgl32.Abs(dx) > 200 || mgl32.Abs(dy) > 200 {
		return
	}
	c.rotatex += dx * c.Sens
	c.rotatey += dy * c.Sens
	if c.rotatey > 89 {
		c.rotatey = 89
	}
	if c.rotatey < -89 {
		c.rotatey = -89
	}
	c.updateAngles()
}

func (c *Camera) OnMoveChange(dir CameraMovement, delta float32) {
	if c.flying {
		delta = 5 * delta
	}
	switch dir {
	case MoveForward:
		if c.flying {
			c.CamPos = c.CamPos.Add(c.front.Mul(delta))
		} else {
			c.CamPos = c.CamPos.Add(c.wfront.Mul(delta))
		}
	case MoveBackward:
		if c.flying {
			c.CamPos = c.CamPos.Sub(c.front.Mul(delta))
		} else {
			c.CamPos = c.CamPos.Sub(c.wfront.Mul(delta))
		}
	case MoveLeft:
		c.CamPos = c.CamPos.Sub(c.right.Mul(delta))
	case MoveRight:
		c.CamPos = c.CamPos.Add(c.right.Mul(delta))
	}
}
func (c *Camera) updateAngles() {
	front := mgl32.Vec3{
		cos(radian(c.rotatey)) * cos(radian(c.rotatex)),
		sin(radian(c.rotatey)),
		cos(radian(c.rotatey)) * sin(radian(c.rotatex)),
	}
	c.front = front.Normalize()
	c.right = c.front.Cross(mgl32.Vec3{0, 1, 0}).Normalize()
	c.up = c.right.Cross(c.front).Normalize()
	c.wfront = mgl32.Vec3{0, 1, 0}.Cross(c.right).Normalize()
}
