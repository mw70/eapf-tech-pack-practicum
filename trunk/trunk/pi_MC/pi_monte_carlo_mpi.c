#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <mpi.h>
#define NUM_ITER 10000000


int main(int argc, char** argv) {
          double dx, dy;
        int i,myrank,commsize;
        unsigned long hits=0,myiter,totalhits;
        long myseed;
        struct timeval now;




           MPI_Init(&argc,&argv);


/* Figure out how many iterations are needed for each process */
        MPI_Comm_rank(MPI_COMM_WORLD,&myrank);
        MPI_Comm_size(MPI_COMM_WORLD,&commsize);
        myiter = (unsigned long) NUM_ITER/commsize;
        if (myrank < NUM_ITER%commsize) 
          myiter++;




 /*  initialize the random number generator with the time  */
         gettimeofday(&now, NULL); // get the current time
        myseed=((long) now.tv_sec)*(myrank + 1);
        srand48(myseed);


        for (i=0; i<myiter; i++) {
          dx = drand48();
          dy = drand48();
/*  Check if the point (dx,dy) is strictly inside the unit circle */
          if ( (dx*dx + dy*dy) < 1)
            hits++;                
        }


        MPI_Reduce(&hits, &totalhits, 1, MPI_UNSIGNED_LONG, MPI_SUM, 0, MPI_COMM_WORLD);


        if (myrank == 0) {
          printf("Pi is approximately %.8g.\n", 4.*totalhits/NUM_ITER);
        }
        MPI_Finalize();
        return 0;
}