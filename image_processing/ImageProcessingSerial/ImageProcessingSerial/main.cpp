//
//  main.cpp
//  ImageProcessingSerial
//
//  Created by Kevin Goldsmith on 9/16/11.
//  Copyright 2011 Adobe Systems, Inc. All rights reserved.
//

// this uses the libjpeg library and some sample code from Junaed Sattar


#include <iostream>
#include <math.h>
#include <jpeglib.h>

using std::cout;
using std::endl;

// typedefs
typedef struct imageData {
    unsigned char *rawImage;
    unsigned int width;
    unsigned int height;
} imageData_t;

// function declarations
imageData_t * readJpegFile( const char * );
int writeJpegFile( const char *filename, imageData_t * imageData );
imageData_t * processImage( imageData_t *input );


imageData_t *
readJpegFile( const char *filename )
{
	struct jpeg_decompress_struct cinfo;
	struct jpeg_error_mgr jerr;
	JSAMPROW row_pointer[1];
    unsigned int width, height, channels = 0;
    unsigned char *raw_image = NULL;
	
	FILE *infile = fopen( filename, "rb" );
	unsigned long location = 0;
	int i = 0;
	
	if ( !infile )
	{
		return NULL;
	}

	cinfo.err = jpeg_std_error( &jerr );
	jpeg_create_decompress( &cinfo );
	jpeg_stdio_src( &cinfo, infile );
	jpeg_read_header( &cinfo, TRUE );
    
    width = cinfo.image_width;
    height = cinfo.image_height;
    channels = cinfo.num_components;
    if ( ( cinfo.out_color_space != JCS_RGB ) || ( channels != 3 ) )
    {
        cout << "this sample only processes 3 channel images" << endl;
        
        jpeg_destroy_decompress(&cinfo);
        fclose(infile);
        return NULL;
    }
    
	jpeg_start_decompress( &cinfo );
	raw_image = (unsigned char*)malloc( cinfo.output_width*cinfo.output_height*cinfo.num_components );
	row_pointer[0] = (unsigned char *)malloc( cinfo.output_width*cinfo.num_components );

	while( cinfo.output_scanline < cinfo.image_height )
	{
		jpeg_read_scanlines( &cinfo, row_pointer, 1 );
		for( i=0; i<cinfo.image_width*cinfo.num_components;i++) 
			raw_image[location++] = row_pointer[0][i];
	}

	jpeg_finish_decompress( &cinfo );
	jpeg_destroy_decompress( &cinfo );
	free( row_pointer[0] );
	fclose( infile );

    imageData_t * returnValue = (imageData_t *) malloc(sizeof(imageData_t));
    if ( returnValue != NULL )
    {
        returnValue->rawImage = raw_image;
        returnValue->width = width;
        returnValue->height = height;
    }
    
	return returnValue;
}

int
writeJpegFile( const char *filename, imageData_t * imageData )
{
	struct jpeg_compress_struct cinfo;
	struct jpeg_error_mgr jerr;
	
	/* this is a pointer to one row of image data */
	JSAMPROW row_pointer[1];
	FILE *outfile = fopen( filename, "wb" );
	
	if ( !outfile )
	{
        cout << "Error opening output jpeg file " << filename << endl;
		return -1;
	}
	cinfo.err = jpeg_std_error( &jerr );
	jpeg_create_compress(&cinfo);
	jpeg_stdio_dest(&cinfo, outfile);
    
	cinfo.image_width = imageData->width;	
	cinfo.image_height = imageData->height;
	cinfo.input_components = 3;
	cinfo.in_color_space = JCS_RGB;

	jpeg_set_defaults( &cinfo );
	jpeg_start_compress( &cinfo, TRUE );

	while( cinfo.next_scanline < cinfo.image_height )
	{
		row_pointer[0] = &imageData->rawImage[ cinfo.next_scanline * cinfo.image_width *  cinfo.input_components];
		jpeg_write_scanlines( &cinfo, row_pointer, 1 );
	}

	jpeg_finish_compress( &cinfo );
	jpeg_destroy_compress( &cinfo );
	fclose( outfile );

	return 1;
}

// copy the buffer
imageData_t *
processImage( imageData_t *input )
{
    imageData_t *output = (imageData_t *) malloc( sizeof(imageData_t) );
    output->width = input->width;
    output->height = input->height;
    output->rawImage = (unsigned char *) malloc( input->width * input->height * 3 );
    unsigned int position = 0;
    
    for ( unsigned int i = 0; i < input->height; i++ )
    {
        for ( unsigned int j = 0; j < input->width; j++ )
        {
            float inRGB[3];
            float outRGB[3];
            for ( unsigned int k = 0; k < 3; k++ )
            {
                inRGB[k] = (input->rawImage[ position + k ]/255.0f);
            }

            outRGB[0] = (inRGB[0] * 0.400f) + (inRGB[1] * 0.769f) + (inRGB[2] * 0.189f);
            outRGB[1] = (inRGB[0] * 0.349f) + (inRGB[1] * 0.686f) + (inRGB[2] * 0.168f);
            outRGB[2] = (inRGB[0] * 0.272f) + (inRGB[1] * 0.534f) + (inRGB[2] * 0.131f);

            for ( unsigned int k = 0; k < 3; k++ )
            {
                outRGB[k] = (outRGB[k] > 1.0) ? 1.0 : outRGB[k];
                output->rawImage[ position++ ] = (unsigned char) (outRGB[k]*255.0f);
            }
        }
    }
    
    return output;
}

int main (int argc, const char * argv[])
{
    imageData_t * inputBuffer = readJpegFile("/Users/kgoldsmi/Development/EAPF Tech Pack 2/ImageProcessingSerial/test.jpg");
    
    if ( inputBuffer == NULL )
    {
        return -1;
    }
    
    clock_t startTime = clock();
    imageData_t * outputBuffer = processImage(inputBuffer);
    clock_t endTime = clock();

    if ( writeJpegFile( "/Users/kgoldsmi/Development/EAPF Tech Pack 2/ImageProcessingSerial/test2.jpg", outputBuffer) < 0 )
    {
        cout << "error writing output file" << endl;
    }
    
    cout << "ticks: " << endTime - startTime << endl;
    
    free(inputBuffer);
    free(outputBuffer);
    
    return 0;
}

