# MTL makefile for pi_integration 10/15/2011
# be sure to source setup.sh

CC=gcc
CXX=g++
MPICC=mpicc

CFLAGS=-g
LDFLAGS=-lm
OPENMP_FLAGS=-fopenmp
CUDA_CFLAGS = -I/usr/local/cuda/include
CUDA_LIBS = -L/usr/local/cuda/lib -lcuda -lcudart
TBB_LIBS = -ltbb -L/opt/intel/tbb/latest/lib/intel64/cc3.4.3_libc2.3.4_kernel2.6.9/
PTHREADS_LIBS=-lpthread
ARBB_LIBS=$(TBB_LIBS) -larbb -L/opt/intel/arbb/1.0.0.022/lib/intel64/ -I/opt/intel/arbb/1.0.0.022/include/

default:  pi_area_serial 

.SUFFIXES: .c .o .cu .cpp .go .6

all:  pi_area_serial pi_area_omp pi_area_go pi_area_arbb pi_area_pthreads
#      pi_area_mpi pi_area_cuda 

pi_area_serial:	pi_area_serial.o
	$(CC) $(LDFLAGS) -o pi_area_serial pi_area_serial.o

pi_area_pthreads: pi_area_pthreads.o
	$(CC) $(LDFLAGS) $(PTHREADS_LIBS) -o pi_area_pthreads pi_area_pthreads.o

pi_area_omp: pi_area_omp.o
	$(CC) $(LDFLAGS) $(OPENMP_FLAGS) -o pi_area_omp pi_area_omp.o

pi_area_go: pi_area_go.go
	go build $<

pi_area_mpi: pi_area_mpi.o
	$(MPICC) $(LDFLAGS) -o pi_area_mpi pi_area_mpi.o
pi_area_mpi.o: pi_area_mpi.c
	$(MPICC) -c $(CFLAGS) pi_area_mpi.c

pi_area_cuda: pi_area_cuda.cu
	nvcc -o pi_area_cuda pi_area_cuda.cu

pi_area_arbb: pi_area_arbb.cpp
	icpc -o pi_area_arbb pi_area_arbb.cpp $(ARBB_LIBS)

clean:
	rm -f *.o *.t pi_area_serial pi_area_omp pi_area_go pi_area_mpi pi_area_cuda

.cu.o:
	$(NVCC) -c $< $(NVCCFLAGS) $(CUDA_CFLAGS) $(CUDA_LIBS)

