bin(x,width)=width*floor(x/width)

bps=1000
testname="TCP\_RR"

set terminal qt noenhanced
set title sprintf("throughput - %s - %db - 100 iterations", testname, bps)
set ylabel "Throughput (bits/s)"
set format y '%.0s %c'
set mxtics 2
set grid xtics mxtics lc rgb "#C0C0C0"

set style data histogram
set style histogram errorbars linewidth 1
set errorbars linecolor black
set bars front
set yrange [0:8000]
set style fill pattern 2
plot "throughput_stats.dat" using 2:($2-($3/2)):($2+($3/2)):xtic(1)