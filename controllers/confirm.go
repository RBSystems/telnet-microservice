package controllers

// import (
// 	"net/http"
// 	"time"
//
// 	"github.com/byuoitav/ftp-microservice/helpers"
// 	"github.com/labstack/echo"
// 	"github.com/ziutek/telnet"
// )
//
// func CommandConfirm(c echo.Context) error {
// 	req := helpers.Request{}
// 	c.Bind(&req)
//
// 	if len(req.Port) < 1 {
// 		req.Port = "41795"
// 	}
//
// 	err := sendCommandWithConfirm(req.Command, req.Address, req.Port)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, "Error: "+err.Error())
// 	}
//
// 	return c.JSON(http.StatusOK, "Success")
// }
//
// // Answer "y" if a command asks for confirmation
// func sendCommandWithConfirm(command string, ipAddress string, port string) error {
// 	var conn *telnet.Conn
//
// 	conn, err := telnet.Dial("tcp", ipAddress+":"+port)
// 	if err != nil {
// 		return err
// 	}
//
// 	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF) This is apparently very important
//
// 	_, err = conn.Write([]byte(command + "\n"))
//
// 	if err != nil {
// 		return err
// 	}
//
// 	time.Sleep(1000 * time.Millisecond) // Wait for the prompt to appear
//
// 	conn.Write([]byte("y")) // Send the yes confirmation
//
// 	return nil
// }
