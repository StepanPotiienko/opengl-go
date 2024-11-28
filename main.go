package main

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type WindowProperties struct {
	width  int16
	height int16
	title  string
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize OpenGL: %v", err))
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)
}

func changeBackgroundColor(colors [4]float32, window glfw.Window) {
	// Set background color (black)
	gl.ClearColor(colors[0], colors[1], colors[2], colors[3])
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Swap buffers to display the content
	window.SwapBuffers()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// Set OpenGL context version
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	windowProps := WindowProperties{width: 800, height: 600, title: "Minecraft"}
	window, err := glfw.CreateWindow(int(windowProps.width), int(windowProps.height), windowProps.title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	initOpenGL()

	blackColor := [4]float32{0.0, 0.0, 0.0, 1.0}

	// Main rendering loop
	for !window.ShouldClose() {
		changeBackgroundColor(blackColor, *window)
		// Poll events (keyboard, mouse, etc.)
		glfw.PollEvents()
	}
}
