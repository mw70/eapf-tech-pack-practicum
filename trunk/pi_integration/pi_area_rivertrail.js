
var num_iter = 1000000;
var width = 2.0/num_iter;
var pos = new ParallelArray(num_iter, function f(idx) {return -1.0 + idx*width;});
var height = pos.map(function f(val) {return Math.sqrt(1.0 - val*val);});
var pi = 2 * width * height.reduce(function f(a,b) {return a+b;});

