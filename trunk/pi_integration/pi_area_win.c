#include <windows.h>
#include <stdio.h>
#include <math.h>
#define NUM_RECT 10000000
#define NUMTHREADS 4
double gPi = 0.0;
CRITICAL_SECTION gCS;
DWORD WINAPI Area(LPVOID pArg)
{
  int myNum = *((int *)pArg);
  double h = 2.0 / NUM_RECT;
  double partialSum = 0.0, x;  // local to each thread
  int begin =  myNum    * (NUM_RECT / NUMTHREADS);
  int end   = (myNum+1) * (NUM_RECT / NUMTHREADS);
  if (nyNum == (NUMTHREADS-1)) end = NUM_RECT;
  for ( int i = begin; i < end; ++i ) {  // compute rectangles in range
      x = -1 + (i + 0.5f) * h;
     partialSum += sqrt(1.0f - x*x) * h;
  }
  EnterCriticalSection(&gCS);
    gPi += partialSum;  // add partial to global final answer
  LeaveCriticalSection(&gCS);
  return 0;
}
int main(int argc, char **argv)
{
  HANDLE threadHandles[NUMTHREADS];
  int tNum[NUMTHREADS];
  InitializeCriticalSection(&gCS);
  for ( int i = 0; i < NUMTHREADS; ++i ){
     tNum[i] = i;
     threadHandles[i] = CreateThread( NULL,   // Security attributes
                                                                                   0,   // Stack size
                                       Area,   // Thread function
                                      (LPVOID)&tNum[i],// Data for Area()
                                       0,       // Thread start mode
                                       NULL);   // Returned thread ID
  }
  WaitForMultipleObjects(NUMTHREADS, threadHandles, TRUE, INFINITE);
  gPi * = 2.0;
  DeleteCriticalSection(&gCS)
  printf("Computed value of Pi:  %12.9f\n", gPi );
}
