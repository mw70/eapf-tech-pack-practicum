CC=gcc
MPICC=mpicc

CFLAGS=-g
OPENMP_FLAGS=-fopenmp

PTHREADS_LIBS=-lpthread

default:  pi_monte_carlo_serial 

all:  pi_monte_carlo_serial pi_monte_carlo_omp pi_monte_carlo_pthreads pi_monte_carlo_mpi

pi_monte_carlo_serial: pi_monte_carlo_serial.o
	$(CC) $(CFLAGS) -o pi_monte_carlo_serial pi_monte_carlo_serial.o

pi_monte_carlo_omp : pi_monte_carlo_omp.o
	$(CC) $(CFLAGS) $(OPENMP_FLAGS) -o pi_monte_carlo_omp pi_monte_carlo_omp.o


pi_monte_carlo_pthreads : pi_monte_carlo_pthreads.o
	$(CC) $(CFLAGS) -o pi_monte_carlo_pthreads $(PTHREADS_LIBS) pi_monte_carlo_pthreads.o

pi_monte_carlo_mpi: pi_monte_carlo_mpi.o
	$(MPICC) $(CFLAGS) -o pi_monte_carlo_mpi pi_monte_carlo_mpi.o

pi_monte_carlo_mpi.o: pi_monte_carlo_mpi.c
	$(MPICC) $(CFLAGS) -c pi_monte_carlo_mpi.c

clean:
	rm -rf *.o pi_monte_carlo_serial pi_monte_carlo_omp pi_monte_carlo_pthreads pi_monte_carlo_mpi





