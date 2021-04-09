# load "throughput_TRP_RR.gnuplot"
set style rectangle back fc rgb "red" fs solid 1.0 border lt -1
bin(x,width)=width*floor(x/width)
bps=1000
testname="TCP\_RR"
set terminal qt noenhanced
set title sprintf("throughput - %s - 100 iterations",testname)
set xlabel "Throughput (bits/s)"
set ylabel "Samples"
set format x '%.0s %c'
set mxtics 2
set grid xtics mxtics lc rgb "#C0C0C0"
set key left top Left reverse
plot \
  sprintf("bare-metal/macbook-to-imac-%s-%db-x100.txt",testname,bps) using (bin($2*bps,100*bps)):(1.0) smooth freq with lines title sprintf("bare metal - %db",bps), \
  sprintf("k8s-calico-bgp/macbook-to-imac-calico-bgp-%s-%db-x100.txt",testname,bps) using (bin($2*bps,100*bps)):(1.0) smooth freq with lines title sprintf("K8S Calico - %db",bps), \
  sprintf("bare-metal-2/macbook-to-imac-bare-metal-%s-%db-x100.txt",testname,bps) using (bin($2*bps,100*bps)):(1.0) smooth freq with lines title sprintf("bare metal 2 - %db",bps), \
  sprintf("docker-swarm/macbook-to-imac-%s-%db-x100.txt",testname,bps) using (bin($2*bps,100*bps)):(1.0) smooth freq with lines title sprintf("docker swarm - %db",bps), \
  sprintf("docker-hostnet/macbook-to-imac-TCP_RR-%db-x100.txt",bps) using (bin($2*bps,100*bps)):(1.0) smooth freq with lines title sprintf("docker host net - %db",bps)