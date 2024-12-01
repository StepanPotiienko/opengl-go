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

const (
	vertexShaderSource = `
		#version 410
		layout(location = 0) in vec3 position;
		void main() {
			gl_Position = vec4(position, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 fragColor;
		void main() {
			fragColor = vec4(0.5, 0.0, 0.0, 1.0); // Red color
		}
	` + "\x00"
)

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize OpenGL: %v", err))
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)
}

func createShaderProgram() uint32 {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vertexSource, freeVertex := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, vertexSource, nil)
	freeVertex()
	gl.CompileShader(vertexShader)

	var success int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetShaderInfoLog(vertexShader, logLength, nil, &log[0])
		panic(fmt.Sprintf("Failed to compile vertex shader: %s", string(log)))
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentSource, freeFragment := gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, fragmentSource, nil)
	freeFragment()
	gl.CompileShader(fragmentShader)

	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetShaderInfoLog(fragmentShader, logLength, nil, &log[0])
		panic(fmt.Sprintf("Failed to compile fragment shader: %s", string(log)))
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	gl.GetProgramiv(program, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])
		panic(fmt.Sprintf("Failed to link program: %s", string(log)))
	}

	return program
}

func createVAO(vertices []float32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return vao
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	windowProps := WindowProperties{width: 800, height: 600, title: "OpenGL Triangle"}
	window, err := glfw.CreateWindow(int(windowProps.width), int(windowProps.height), windowProps.title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	initOpenGL()

	vertices := []float32{
		-0.5, -0.5, 0.0, // Bottom-left
		0.5, -0.5, 0.0, // Bottom-right
		0.0, 0.5, 0.0, // Top
	}

	shaderProgram := createShaderProgram()
	vao := createVAO(vertices)
	defer gl.DeleteVertexArrays(1, &vao)

	for !window.ShouldClose() {
		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
