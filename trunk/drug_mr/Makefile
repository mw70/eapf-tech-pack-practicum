# Linux makefile for drug_mr 10/15/2011
# be sure to source setup.sh

CC=gcc
CXX=g++
MPICC=mpicc
CUDA_INSTALL = /usr/local/cuda
TBB_INSTALL = /opt/intel/tbb
BOOST_INSTALL = /opt/boost
ARBB_INSTALL = /opt/intel/arbb/1.0.0.022

CFLAGS=-g
LDFLAGS=-lm
OPENMP_FLAGS=-fopenmp
CUDA_CFLAGS = -I$(CUDA_INSTALL)/include
CUDA_LIBS = -L$(CUDA_INSTALL)/lib -lcuda -lcudart
TBB_LIBS = -I$(TBB_INSTALL)/latest/include -ltbb -L$(TBB_INSTALL)/latest/lib/intel64/cc3.4.3_libc2.3.4_kernel2.6.9/
PTHREADS_LIBS=-lpthread
BOOST_LIBS = -I$(BOOST_INSTALL)/include -L$(BOOST_INSTALL)/lib -lboost_thread
ARBB_LIBS=$(TBB_LIBS) -larbb -L$(ARBB_INSTALL)/lib/intel64/ -I$(ARBB_INSTALL)/include/

default:  dd_serial # dd_omp dd_boost dd_go dd_hadoop
#	hadoop jar DDHadoop.jar edu.stolaf.cs.DDHadoop $(DFS)/in $(DFS)/out

.SUFFIXES: .c .o .cu .cpp .go

dd_serial:  dd_serial.cpp
	$(CXX) -o dd_serial dd_serial.cpp

dd_omp:  dd_omp.cpp
	$(CXX) $(LDFLAGS) $(OPENMP_FLAGS) $(TBB_LIBS) -o dd_omp dd_omp.cpp

dd_boost:  dd_boost.cpp
	$(CXX) $(LDFLAGS) $(TBB_LIBS) $(BOOST_LIBS) -o dd_boost dd_boost.cpp 

dd_go: dd_go.go
	go build $<

dd_hadoop:  DDHadoop.jar

DDHadoop.jar:  classes DDHadoop.java
	javac -classpath $(HADOOP_LIB) -d classes DDHadoop.java
	jar cvf DDHadoop.jar -C classes edu

classes:  
	mkdir classes

clean:
	rm -f *.o dd_serial dd_omp dd_boost dd_go
