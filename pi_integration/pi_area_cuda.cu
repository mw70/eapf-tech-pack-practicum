/*  calculating pi via area under the curve
 *  This code uses an algorithm fairly easily ported to all parallel methods.
 *  Since it calculates pi, it is easy to verify that results are correct.
 *  It can also be used to explore accuracy of results and techniques for managing error.
 */


#include <stdio.h>
#include <stdlib.h>
#include <cuda.h>


#define NUMRECT 10000000


/*  students learn in grammar school that the area of a circle is pi*radius*radius.
 *  They learn in high school that the formula of a circle is x^2 + y^2 = radius^2.
 *
 *  These facts allows students calculating pi by estimating area of mid-point rectangles
 *
 *  Area of unit circle is pi, y = sqrt(1-x^2) is formula for semicircle from -1 to 1
 */


// constants useful to CUDA
 const int        threadsPerBlock        = 256; 
 const int        blocksPerGrid                =  32;
 const int        totalThreads                = threadsPerBlock * blocksPerGrid;
 const float        overallWidth                 = 2.0f; 
 const float        block_width                 = overallWidth / blocksPerGrid;


  __global__ void calcArea(int *d_rectPerThread, float *d_width, float *partPiByBlock) {
        __shared__ float partPiByThread[threadsPerBlock];
        int reduce_i = blockDim.x / 2;        // index for reducing thread results to single block value
        float width = *d_width;
        int rectPerThread = *d_rectPerThread;
        float x                = -1.0f 
+  (overallWidth * blockIdx.x) / blocksPerGrid
+  (block_width * threadIdx.x) / threadsPerBlock
-  width / 2;
        float partPi         = 0.0f;
        
            for (int i = 0; i < rectPerThread; i++) {
                x += width;
                partPi += width * sqrtf(1.0f - x * x);
            }
            partPiByThread[threadIdx.x] = partPi;


        // reduce all threads in the block to a single block value
        while (reduce_i != 0) {
                __syncthreads();
                if (threadIdx.x < reduce_i)
                        partPiByThread[threadIdx.x] += partPiByThread[threadIdx.x + reduce_i]; 
                reduce_i /= 2;
        }
        
        // store block result in correct spot for reducing on CPU side
        if (threadIdx.x == 0) 
                partPiByBlock[blockIdx.x] = partPiByThread[0];
 }




int main(int argc, char **argv) {


        int        numRect;                                        // number of rectangles
        int *d_rectPerThread, rectPerThread;                        // number of rectangles per thread
        int        i;                                                // loop index
        float        *d_width, width;                                // width of each rectangle
        float   *d_partPiByBlock, h_partPiByBlock[blocksPerGrid]; // partial pi values returned by CUDA
        float        pi, halfPI = 0.0;                                // sum of area of rectangles gives pi/2


        numRect = argc == 2 ? atoi(argv[1]) : NUMRECT;        // get number of rectangles
        rectPerThread = numRect / totalThreads;
        numRect = rectPerThread * totalThreads;
        width = overallWidth / numRect;                        // calculate width of each rectangle


        cudaMalloc((void**)&d_rectPerThread, sizeof(int));
        cudaMalloc((void**)&d_width, sizeof(int));
        cudaMalloc((void**)&d_partPiByBlock, sizeof(float) * blocksPerGrid); 
                
        cudaMemcpy(d_rectPerThread, &rectPerThread, sizeof(int), cudaMemcpyHostToDevice);
        cudaMemcpy(d_width, &width, sizeof(float), cudaMemcpyHostToDevice);
        
                calcArea <<<blocksPerGrid, threadsPerBlock>>> (d_rectPerThread, d_width, d_partPiByBlock);
        
        cudaMemcpy(h_partPiByBlock, d_partPiByBlock, sizeof(float) * blocksPerGrid, 
cudaMemcpyDeviceToHost);


        for (i = 0; i < blocksPerGrid; ++i) 
                halfPI += h_partPiByBlock[i];        /* calculate pi/4, with room for better error mgmt */
        pi = 2.0 * halfPI;


        printf ("\n==\n==\t%20s = %15.10f\n",                "pi",                        pi);
        printf ("==\t%20s = %15d\n",                        "total rectangles",         numRect);
        printf ("==\t%20s = %15d\n==\n\n",                "CUDA threads",        totalThreads);
        return 0;
}