package main
  
// Importing fmt and time
import (
    "fmt"
    "time"
)
  
// Main function
func main() {
  
    // Calling Sleep method
    for {
    time.Sleep(1 * time.Second)
  
    // Printed after sleep is over
    fmt.Println("Hello this is very cool.")
  }
}
