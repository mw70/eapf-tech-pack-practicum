#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#define NUM_ITER 10000000


int main(int argc, char** argv) {
         double dx, dy;
        int i;
        long long hits=0;
#pragma omp parallel
{
        short state[3];
        long myseed;
        struct timeval now;


        gettimeofday(&now, NULL); // get the current time
/*  initialize the random number generator with the time  */
        myseed=((long) now.tv_sec)*(omp_get_thread_num()+1);
        state[0]=(short) myseed;
        state[1]=(short) myseed>>16;
        state[2]=(short) myseed<<3;




#pragma omp for reduction(+:hits) private(dx,dy)
        for (i=0; i<NUM_ITER; i++) {
dx = erand48(state);
dy = erand48(state);
  /*  Check if the point (dx,dy) is strictly inside the unit circle */
if ( (dx*dx + dy*dy) < 1)
hits++;                
        }


} // end of the parallel block.  Back to one thread only.


printf("Pi is approximately %g.\n", 4.*hits/NUM_ITER);
}
