%.o:%.c %.h
	gcc -c -g -m64 -o $@ $<

%.o:%.cc %.hh
	g++ -c -g -m64 -o $@ $<

libbliss.a: bliss_C.o defs.o graph.o heap.o orbit.o partition.o timer.o uintseqhash.o utils.o
	ar crf $@ $^

.PHONY: clean
clean:
	-@rm libbliss.a *.o 2>/dev/null

