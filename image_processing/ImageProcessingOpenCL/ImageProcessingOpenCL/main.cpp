//
//  main.cpp
//  ImageProcessingOpenCL
//
//  Created by Kevin Goldsmith on 9/17/11.
//  Copyright 2011 Adobe Systems, Inc. All rights reserved.
//

#include <iostream>
#include <OpenCL/OpenCL.h>

int main (int argc, const char * argv[])
{
    cl_device_id deviceID;
    int err = clGetDeviceIDs(NULL, CL_DEVICE_TYPE_GPU, 1, &deviceID, NULL);
    if ( err != CL_SUCCESS )
    {
        std::cout << "error getting first GPU device!";
        return EXIT_FAILURE;
    }

    cl_context context = clCreateContext(NULL, 1, &deviceID, NULL, NULL, &err);
    if ( ( err != CL_SUCCESS ) || !context )
    {
        std::cout << "error creating context!";
        return EXIT_FAILURE;
    }
    
    cl_command_queue commandQueue = clCreateCommandQueue(context, deviceID, 0, &err);
    if ( ( err != CL_SUCCESS ) || !commandQueue )
    {
        std::cout << "error creating Command Queue!";
        return EXIT_FAILURE;
    }

    const char *KernelSource = "\n" \
    "__kernel void square(                                                  \n" \
    "   __global float* input,                                              \n" \
    "   __global float* output,                                             \n" \
    "   const unsigned int count)                                           \n" \
    "{                                                                      \n" \
    "   int i = get_global_id(0);                                           \n" \
    "   if(i < count)                                                       \n" \
    "       output[i] = input[i] * input[i];                                \n" \
    "}                                                                      \n" \
    "\n";
    
    cl_program program = clCreateProgramWithSource(context, 1, (const char **) &KernelSource, NULL, &err);
    if ( ( err != CL_SUCCESS ) || !program )
    {
        std::cout << "error creating Program!";
        return EXIT_FAILURE;
    }

    // insert code here...
    std::cout << "Hello, World!\n";
    return EXIT_SUCCESS;
}

