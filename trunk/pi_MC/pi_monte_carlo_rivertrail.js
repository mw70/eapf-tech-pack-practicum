
var num_iter = 1000000;
var hit_check = new ParallelArray(num_iter, function f(idx) { var x=Math.random(); var y=Math.random(); return (x*x+y*y <=1 ) ? 1 : 0;});
var pi = 4 * hit_check.reduce(function f(a,b) { return a+b;})/num_iter;
