#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#define NUM_ITER 10000000


int main(int argc, char** argv) {
         double dx, dy;
        int i;
        long long hits=0;
        struct timeval now;


        gettimeofday(&now, NULL); // get the current time
/*  initialize the random number generator with the time  */
        srand48( (long) now.tv_sec);


        for (i=0; i<NUM_ITER; i++) {
dx = drand48();
dy = drand48();
                /*  Check if the point (dx,dy) is strictly inside the unit circle */
if ( (dx*dx + dy*dy) < 1)
hits++;                
        }


        printf("Pi is approximately %g.\n", 4.*hits/NUM_ITER);
}
