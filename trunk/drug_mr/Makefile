# MTL makefile for drug_mr 10/15/2011
# be sure to source setup.sh

CC=gcc
CXX=g++
MPICC=mpicc
GOC=6g
GOLINK=6l

CFLAGS=-g
LDFLAGS=-lm
OPENMP_FLAGS=-fopenmp
CUDA_CFLAGS = -I/usr/local/cuda/include
CUDA_LIBS = -L/usr/local/cuda/lib -lcuda -lcudart
TBB_LIBS = -I/opt/intel/tbb/latest/include -ltbb -L/opt/intel/tbb/latest/lib/intel64/cc3.4.3_libc2.3.4_kernel2.6.9/
PTHREADS_LIBS=-lpthread
BOOST_LIBS = -I/opt/boost/include -L/opt/boost/lib -lboost_thread
ARBB_LIBS=$(TBB_LIBS) -larbb -L/opt/intel/arbb/1.0.0.022/lib/intel64/ -I/opt/intel/arbb/1.0.0.022/include/

all:  dd_serial dd_omp dd_boost # dd_hadoop
#	hadoop jar DDHadoop.jar edu.stolaf.cs.DDHadoop $(DFS)/in $(DFS)/out

.SUFFIXES: .c .o .cu .cpp .go .6

dd_serial:  dd_serial.cpp
	$(CXX) -o dd_serial dd_serial.cpp

dd_omp:  dd_omp.cpp
	$(CXX) $(LDFLAGS) $(OPENMP_FLAGS) $(TBB_LIBS) -o dd_omp dd_omp.cpp

dd_boost:  dd_boost.cpp
	$(CXX) $(LDFLAGS) $(TBB_LIBS) $(BOOST_LIBS) -o dd_boost dd_boost.cpp 

dd_hadoop:  DDHadoop.jar

DDHadoop.jar:  classes DDHadoop.java
	javac -classpath $(HADOOP_LIB) -d classes DDHadoop.java
	jar cvf DDHadoop.jar -C classes edu

classes:  
	mkdir classes

