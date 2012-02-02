// Estimate pi as twice the area under a semicircle

package main

import (
       "fmt"
       "flag"
       "math"
       "runtime"
)

func main() {
     // Command line flags, run as
     //  ./pi_area_go -rect=10000000 -grain=10000 -cpu=4
     nRect := flag.Int("rect", 1e7, "total number of rectangles")
     nGrain := flag.Int("grain", 1e4, "number of rectangles per task")
     nCPU := flag.Int("cpu", 1, "number of threads to use")
     flag.Parse()
     runtime.GOMAXPROCS(*nCPU) //set the number of OS threads to use.
     
     nParts := 0   // keeps track of the number of tasks
     parts := make(chan float64) // Channel delivery floats to collect parts of Pi
     for i := 0; i < *nRect; i += *nGrain {
       nParts += 1 
       end := i + *nGrain // Calculate the highest index value for this task
       if end > *nRect {  // This is the last task allocation and it doesn't divide evenly
          end = *nRect
	}
       
       go func(begin, end int)  {
          sum := 0.0
	  h := 2.0 / float64(*nRect)
	  for i := begin; i< end; i++ {
	    x := -1 + (float64(i)+0.5)*h
	    sum += math.Sqrt(1-x*x) * h
	  }
	  parts <- sum // Return the partial sum on the channel
	}(i,end)
      }

      sum := 0.0  //Note this is in the main thread, *not* the local declarations in each task
      for j := 1; j < nParts; j++ {
        sum += <- parts
      }
      pi := sum * 2.0
      fmt.Printf("Pi = %.8g using %d tasks on %d threads.\n", pi, nParts, *nCPU)
}

