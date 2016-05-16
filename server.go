package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/byuoitav/ftp-microservice/helpers"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
	"github.com/ziutek/telnet"
)

func sendCommand(c echo.Context) error {
	req := helpers.Request{}
	c.Bind(&req)

	var conn *telnet.Conn

	if req.Port == "" {
		req.Port = ":23"
	}

	conn, err := telnet.Dial("tcp", req.IPAddress+":"+req.Port)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error with contacting host: "+err.Error())
	}

	defer conn.Close()
	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)

	// Cheap cop-out way to deal with getting the version of the touchpanels. Split out xmodem into own endpoint?
	// TODO: Figure out a better way to handle this
	if strings.EqualFold(req.Command, "xget ~.LocalInfo.vtpage") {
		resp, err := getProjectInfo(req, conn)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error with contacting host: "+err.Error())
		}

		return c.JSON(http.StatusOK, strings.TrimSpace(string(resp)))
	}

	conn.SetReadDeadline(time.Now().Add(45 * time.Second))

	if req.Prompt == "" {
		p, err := getPrompt(req, conn)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error with contacting host %s", err.Error())
			return
		}

		req.Prompt = p
	}

	_, err = conn.Write([]byte(req.Command + "\n\n")) // Send a second newline so we get the prompt

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}
	err = conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}

	response, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	response = response[:len(response)-len(req.Prompt)] // Ghetto trim the prompt off the response
	response = response[len(req.Command):]

	switch req.Command {
	case "iptable":
		ipTable, err := getIPTable(string(response))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		bits, err := json.Marshal(ipTable)
		fmt.Fprintf(w, "%s", bits)
		return
	default:
		fmt.Fprintf(w, "%s", strings.TrimSpace(string(response)))
		return
	}
}

func getProjectInfo(req request, conn *telnet.Conn) (string, error) {
	fmt.Printf("%s Getting project info...\n", req.IPAddress)

	defer conn.Close()

	if req.Prompt == "" {
		prompt, _ := getPrompt(req, conn)
		req.Prompt = prompt
	}
	conn.Write([]byte("udir \\romdisk\\user\\display\\\n\n"))
	conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	resp1, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s %s\n", req.IPAddress, resp1)

	if !strings.Contains(string(resp1), ".vtpage") {
		return "File ~.LocalInfo.vtpage does not exist.\n", nil
	}

	conn.Write([]byte("cd \\romdisk\\user\\display\\\n"))
	conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	resp, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s\n", resp)
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	conn.Write([]byte(req.Command + "\n\n"))

	conn.SkipUntil("[BEGIN_INFO]", "ERROR")
	fmt.Printf("%s skipped\n", req.IPAddress)
	resp, err = conn.ReadUntil("[END_INFO]", "Panel", "not")
	if err != nil {
		return "", err
	}

	fmt.Printf("%s Response: %s\n", req.IPAddress, string(resp))

	conn.Close() // Actively close the xmodem connection

	return string(resp), nil
}

func sendCommandConfirm(c echo.Context) error {
	bits, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body: %s\n", err.Error())
		return
	}

	var req request

	err = json.Unmarshal(bits, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error with the request body: %s", err.Error())
		return
	}

	if len(req.Port) < 1 {
		req.Port = "41795"
	}

	err = sendCommandWithConfirm(req.Command, req.IPAddress, req.Port)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "Success!")
}

// Answer "y" if a command asks for confirmation
func sendCommandWithConfirm(command string, ipAddress string, port string) error {
	var conn *telnet.Conn

	conn, err := telnet.Dial("tcp", ipAddress+":"+port)
	if err != nil {
		return err
	}

	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF) This is apparently very important

	_, err = conn.Write([]byte(command + "\n"))

	if err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond) // Wait for the prompt to appear

	conn.Write([]byte("y")) // Send the yes confirmation

	return nil
}

func main() {
	port := ":8001"
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Get("/health", health.Check)

	e.Post("/sendCommand", controllers.sendCommand)
	e.Post("/sendCommandConfirm", controllers.sendCommandConfirm)
	e.Post("/sendCommand/getPrompt", controllers.getPromptHandler)

	fmt.Printf("The Telnet Microservice is listening on %s\n", port)
	e.Run(fasthttp.New(port))
}
