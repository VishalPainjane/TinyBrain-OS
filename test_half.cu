#include <cuda_fp16.h>
#include <stdio.h>
int main() {
    unsigned short val = 0x7FFF; // NaN
    half h = *(half*)&val;
    float f = __half2float(h);
    printf("NaN: %f\n", f);
    
    val = 0x7C00; // Inf
    h = *(half*)&val;
    f = __half2float(h);
    printf("Inf: %f\n", f);
    
    val = 0x0000;
    h = *(half*)&val;
    f = __half2float(h);
    printf("Zero: %f\n", f);
    
    return 0;
}