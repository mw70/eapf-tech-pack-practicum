# source this on bash before make in order to define environment variables 
# for Tech Pack 2 example codes on Manycore Testing Lab 10/15/2011

# Go
export PATH="$PATH:$HOME/go/bin"

# Intel tools (compiler, TBB, and ArBB)
export PATH="$PATH:/opt/intel/bin"
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/opt/intel/arbb/1.0.0.022/lib/intel64/"
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH::/opt/intel/tbb/latest/lib/intel64/cc3.4.3_libc2.3.4_kernel2.6.9"
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH::/opt/boost/lib"
