#include <stdio.h>
#include <string.h>
#include <sys/mman.h>
#include <errno.h>
#include <time.h>

#define RUNS 10
#define LENGTH (2UL*1024*1024*1024) // 2GB
#define PATTERN 0b10101010

int dewit() {
  char *mem = mmap(
    /* addr= */ NULL,
    /* length= */ LENGTH,
    /* prot= */ PROT_READ | PROT_WRITE,
    /* flags= */ MAP_SHARED | MAP_ANONYMOUS | MAP_HUGETLB,
    /* fd= */ -1,
    /* offset= */ 0);
  if(mem == MAP_FAILED){
    printf("Mapping Failed: %s\n", strerror(errno));
    return 1;
  }

  // zero-pad entire memory block first, since mmap actually only allocates on first use
  for (long i = 0; i < LENGTH; ++i) {
    *(mem+i) = 0b00000000;
  }

  clock_t start, diff;
  start = clock();
  for (long i = 0; i < LENGTH; ++i) {
    *(mem+i) = PATTERN;
  }
  diff = clock() - start;
  printf("%d\t", diff);
  fflush(stdout);

  start = clock();
  for (long i = 0; i < LENGTH; ++i) {
    if ((char)*(mem+i) != (char)PATTERN) {
      printf("corrupt data at %d\n: %X (expected %X)", mem+i, *(mem+i), PATTERN);
      return 0;
    }
  }
  diff = clock() - start;
  printf("%d\n", diff);

	munmap(mem, LENGTH);
   
  return 0;
}

int main() {
  int r;
  printf("%d READ/WRITE speed tests in CPU clock ticks for a %d bytes memory region (%d clock ticks per s)\n", RUNS, LENGTH, CLOCKS_PER_SEC);
  printf("WRITE\t\tREAD\n");
  for (size_t i = 0; i < RUNS; ++i) {
    r = dewit();
    if (r) {
      return r;
    }
  }
  return 0;
}