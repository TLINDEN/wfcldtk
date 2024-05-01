out=images/output.png
in=images/inputtilemap.png
w=8
h=8
cell=100



all: build render show

build:
	go build .

render:
	./wfcldtk $(in) $(w) $(h) $(cell) $(out)

show:
	display $(out)
