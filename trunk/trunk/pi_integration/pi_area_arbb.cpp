# include <iostream>
# include <cstdlib>
# include <arbb.hpp>
const long num_steps=10485760;
void calc_height(arbb::f64 &y, arbb::f64 width, arbb::f64 i)
{
            arbb::f64 x = (i+0.5f)*width - 1.0f;
            y = arbb::sqrt(1.0 - x*x);
}
void calc_pi(arbb::dense<arbb::f64> &area, arbb::f64 width, arbb::dense<arbb::f64> iVals)
{
            arbb::map(calc_height)(area, width, iVals);
            area = area * width;
}
int main(int argc, char *argv[])
{
            arbb::dense<arbb::f64> iterations = arbb::indices(arbb::f64(0.0), num_steps, arbb::f64(1.0));
            arbb::dense<arbb::f64> areas(num_steps);
            arbb::f64 width = 2.0f / num_steps;
        arbb::call(calc_pi)(areas, width, iterations);
        arbb::f64 pi = arbb::sum(areas) * 2.0f;
        std::cout << "Pi =" << arbb::value(pi) << std::endl;
            return 0;
}