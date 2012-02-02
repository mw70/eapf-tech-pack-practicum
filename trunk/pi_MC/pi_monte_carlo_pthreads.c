#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <pthread.h>
#define NUM_ITER 10000000
#define NUM_THREADS 8


unsigned long totalhits;
pthread_mutex_t hits_lock;


void *calc_hits(int *mynum) {
        unsigned long localhits=0,myiter;        
            int i;
        double dx, dy;
        short state[3];
        long myseed;
        struct timeval now;


        gettimeofday(&now, NULL); // get the current time
 /*  initialize the random number generator with the time  */
        myseed=((long) now.tv_sec)*(*mynum + 1);
        state[0]=(short) myseed;
        state[1]=(short) myseed>>16;
        state[2]=(short) myseed<<3;


        myiter = (unsigned long) NUM_ITER/NUM_THREADS;
        if (*mynum < NUM_ITER%NUM_THREADS) 
          myiter++;




        for (i=0; i<myiter; i++) {
          dx = erand48(state);
          dy = erand48(state);
/*  Check if the point (dx,dy) is strictly inside the unit circle */
          if ( (dx*dx + dy*dy) < 1)
            localhits++;                
        }


        pthread_mutex_lock(&hits_lock);
        totalhits += localhits;
        pthread_mutex_unlock(&hits_lock);
        pthread_exit(0);
}


int main(int argc, char** argv) {
        int j;
double Pi;
pthread_t threads[NUM_THREADS];
        int tnum[NUM_THREADS];
        
        pthread_mutex_init(&hits_lock,NULL);
        for (j=0;j<NUM_THREADS;j++){
          tnum[j] = j;
          pthread_create(&(threads[j]), NULL, calc_hits,  &tnum[j]);
        }
        for(j=0;j<NUM_THREADS;j++){
          pthread_join(threads[j], NULL);
        }
        Pi = 4.0f * totalhits/NUM_ITER;
        printf("Pi is approximately %.8g.\n", Pi);
        return 0;
}