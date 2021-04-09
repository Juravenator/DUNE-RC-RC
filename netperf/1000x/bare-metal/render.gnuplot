# load "render.gnuplot"
# set style line 1 lc rgb '#8b1a0e' pt 1 ps 1 lt 1 lw 2 # --- red
# set style line 2 lc rgb '#5e9c36' pt 6 ps 1 lt 1 lw 2 # --- green
bin(x,width)=width*floor(x/width)
set title "bare metal throughput - 1.000 iterations"
set xlabel "Throughput (bits/s)"
set ylabel "Samples"
set format x '%.0s %c'
set mxtics 2
set grid xtics mxtics lc rgb "#C0C0C0"
# set xrange [6000 to 14000]
set key left top Left reverse
# set xtics 0, 1000
# set style fill transparent solid .4
# set style fill pattern
plot \
  'macbook-to-imac-TCP_RR-1b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_RR - 1b', \
  'macbook-to-imac-TCP_STREAM-1b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_STREAM - 1b', \
  'macbook-to-imac-UDP_RR-1b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_RR - 1b', \
  'macbook-to-imac-UDP_STREAM-1b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_STREAM - 1b', \
  'macbook-to-imac-TCP_RR-500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_RR - 500b', \
  'macbook-to-imac-TCP_STREAM-500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_STREAM - 500b', \
  'macbook-to-imac-UDP_RR-500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_RR - 500b', \
  'macbook-to-imac-UDP_STREAM-500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_STREAM - 500b', \
  'macbook-to-imac-TCP_RR-1000b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_RR - 1000b', \
  'macbook-to-imac-TCP_STREAM-1000b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_STREAM - 1000b', \
  'macbook-to-imac-UDP_RR-1000b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_RR - 1000b', \
  'macbook-to-imac-UDP_STREAM-1000b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_STREAM - 1000b', \
  'macbook-to-imac-TCP_RR-1500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_RR - 1500b', \
  'macbook-to-imac-TCP_STREAM-1500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'TCP\_STREAM - 1500b', \
  'macbook-to-imac-UDP_RR-1500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_RR - 1500b', \
  'macbook-to-imac-UDP_STREAM-1500b-x1000.txt' using (bin($2,100)):(1.0) smooth freq with boxes title 'UDP\_STREAM - 1500b'
